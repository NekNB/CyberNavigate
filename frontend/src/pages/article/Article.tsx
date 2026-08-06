import DOMPurify from "dompurify";
import type { FC } from "react";
import { useEffect, useState } from "react";
import {
  GetArticleText,
  GetPublishedArticles,
} from "../../app/api/Article/Article";
import Footer from "../../components/Footer/Footer";
import Header from "../../components/Header/Header";
import { EArticleStatus, type IArticle } from "../../types/articles";
import styles from "./Article.module.css";
const Article: FC = () => {
  const [articles, setArticles] = useState<IArticle[]>([]);
  const [selectedId, setSelectedId] = useState<string>();
  const [articleContent, setArticleContent] = useState<string>("");

  const [searchQuery, setSearchQuery] = useState<string>("");
  const [isLoadingList, setIsLoadingList] = useState<boolean>(true);
  const [isLoadingText, setIsLoadingText] = useState<boolean>(false);
  const [error, setError] = useState<string>();

  // Управление мобильной шторкой
  const [isMenuOpen, setIsMenuOpen] = useState<boolean>(false);

  // 1. Загрузка списка всех статей
  useEffect(() => {
    const fetchArticlesList = async () => {
      try {
        setIsLoadingList(true);
        const articles = await GetPublishedArticles();

        // const data: ArticleHeader[] = await response.json();
        setArticles(articles);

        if (articles.length > 0) {
          const firstId = articles[0].id;
          if (firstId !== undefined) setSelectedId(firstId);
        }
      } catch (err) {
        console.error("Не удалось загрузить статьи:", err);
        setError("Не удалось загрузить список статей");
      } finally {
        setIsLoadingList(false);
      }
    };

    fetchArticlesList();
  }, []);

  // 2. Загрузка текста выбранной статьи
  useEffect(() => {
    if (!selectedId) return;

    const fetchArticleText = async () => {
      try {
        setIsLoadingText(true);
        const articleText = await GetArticleText(selectedId);
        const clearHTML = DOMPurify.sanitize(articleText);
        setArticleContent(clearHTML);
      } catch (err) {
        console.error("Ошибка загрузки текста:", err);
        setArticleContent("Ошибка при загрузке содержимого статьи.");
      } finally {
        setIsLoadingText(false);
      }
    };

    fetchArticleText();
  }, [selectedId]);

  // Двухэтапная фильтрация для поиска
  const searchLower = searchQuery.trim().toLowerCase();

  // 1. Совпадения по НАЗВАНИЮ
  const matchesByTitle = articles.filter((item) => {
    const name = (item.title || item.title || "").toLowerCase();
    return searchLower !== "" && name.includes(searchLower);
  });

  // 2. Совпадения по ТЕКСТУ (исключая те, что уже найдены по названию)
  const matchesByText = articles.filter((item) => {
    const name = (item.title || item.title || "").toLowerCase();
    const bodyText = // item.content ||
      // item.text ||
      (articleContent || "").toLowerCase();

    const isTitleMatch = name.includes(searchLower);
    const isTextMatch = bodyText.includes(searchLower);

    return searchLower !== "" && !isTitleMatch && isTextMatch;
  });

  // Полный неизменный список для шторки шапки
  const fullHeaderArticles = articles.map((item): IArticle => ({
    id: item.id!,
    title: item.title || "Без названия",
    status: EArticleStatus.PUBLISHED,
  }));

  const activeArticle = articles.find((item) => item.id === selectedId);

  return (
    <>
      <Header
        articles={fullHeaderArticles}
        activeArticleId={selectedId ?? undefined}
        onSelectArticle={(id) => setSelectedId(id)}
        isMenuOpen={isMenuOpen}
        setIsMenuOpen={setIsMenuOpen}
      />

      <main className={styles.main}>
        {/* Левая колонка сайдбара (ПК) */}
        <div className={styles.articleList}>
          {isLoadingList && (
            <div className={styles.statusMessage}>Загрузка статей...</div>
          )}
          {error && (
            <div className={styles.statusMessage} style={{ color: "#ff8888" }}>
              {error}
            </div>
          )}

          {!isLoadingList && !error && articles.length === 0 && (
            <div className={styles.statusMessage}>Список пуст</div>
          )}

          {!isLoadingList &&
            !error &&
            articles.map((item) => {
              const id = item.id;
              const name = item.title || item.title || "Без названия";
              const isActive = selectedId === id;

              return (
                <div
                  key={id}
                  className={`${styles.articleItem} ${isActive ? styles.active : ""}`}
                  onClick={() => id !== undefined && setSelectedId(id)}
                >
                  {name}
                </div>
              );
            })}
        </div>

        {/* Правая колонка: Поиск и статья */}
        <div className={styles.articlePage}>
          <div className={styles.searchBox}>
            <button
              className={styles.mobileMenuTrigger}
              onClick={() => setIsMenuOpen(true)}
            >
              ☰ Выбрать статью из списка
            </button>

            <div className={styles.searchInputWrapper}>
              <input
                type="text"
                className={styles.searchInput}
                placeholder="Поиск по статьям..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />

              {/* Выпадающий список совпадений с разбивкой по типам */}
              {searchQuery.trim() !== "" && (
                <div className={styles.searchResultsDropdown}>
                  {matchesByTitle.length === 0 && matchesByText.length === 0 ? (
                    <div className={styles.searchNoResults}>
                      Ничего не найдено по запросу «{searchQuery}»
                    </div>
                  ) : (
                    <>
                      {/* Блок 1: Совпадения по названию */}
                      {matchesByTitle.length > 0 && (
                        <div className={styles.searchCategoryGroup}>
                          <div className={styles.searchCategoryHeader}>
                            📌 ПО НАЗВАНИЮ
                          </div>
                          {matchesByTitle.map((item) => {
                            const id = item.id;
                            const name =
                              item.title || item.title || "Без названия";
                            return (
                              <div
                                key={`title-${id}`}
                                className={styles.searchResultItem}
                                onClick={() => {
                                  if (id !== undefined) setSelectedId(id);
                                  setSearchQuery("");
                                }}
                              >
                                🔍 {name}
                              </div>
                            );
                          })}
                        </div>
                      )}

                      {/* Блок 2: Совпадения по тексту */}
                      {matchesByText.length > 0 && (
                        <div className={styles.searchCategoryGroup}>
                          <div className={styles.searchCategoryHeader}>
                            📄 В ТЕКСТЕ СТАТЕЙ
                          </div>
                          {matchesByText.map((item) => {
                            const id = item.id;
                            const name =
                              item.title || item.title || "Без названия";
                            return (
                              <div
                                key={`text-${id}`}
                                className={styles.searchResultItem}
                                onClick={() => {
                                  if (id !== undefined) setSelectedId(id);
                                  setSearchQuery("");
                                }}
                              >
                                📝 {name}
                              </div>
                            );
                          })}
                        </div>
                      )}
                    </>
                  )}
                </div>
              )}
            </div>
          </div>

          {/* Карточка статьи */}
          <article className={styles.articleCard}>
            {activeArticle ? (
              <>
                <h1 className={styles.articleTitle}>
                  {activeArticle.title || activeArticle.title}
                </h1>
                <div className={styles.articleText}>
                  {isLoadingText ? (
                    <p>Загрузка текста статьи...</p>
                  ) : (
                    <div
                      style={{ whiteSpace: "pre-line" }}
                      dangerouslySetInnerHTML={{ __html: articleContent }}
                    ></div>
                  )}
                </div>
              </>
            ) : (
              <div className={styles.articleText}>
                <p>Выберите статью из списка</p>
              </div>
            )}
          </article>
        </div>
      </main>
      <Footer />
    </>
  );
};

export default Article;
