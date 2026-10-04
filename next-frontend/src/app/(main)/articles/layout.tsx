import "server-only";

import Search from "@/app/(main)/articles/_components/Search";

import MobileSideBar from "@/components/ui/MobileSideBar/MobileSideBar";
import { GetPublishedArticles } from "@/lib/api-server/Article/Article";
import { EArticleStatus, type IArticle } from "@/types/articles";
import styles from "./Articles.module.css";
import SideBar from "./_components/SideBar";

const ArticlesLayout = async ({ children }: { children: React.ReactNode }) => {
  const articlesList: IArticle[] = await GetPublishedArticles();

  // Полный неизменный список для шторки шапки
  const fullHeaderArticles = articlesList.map(
    (item): IArticle => ({
      id: item.id!,
      title: item.title || "Без названия",
      status: EArticleStatus.PUBLISHED,
    }),
  );

  return (
    <main className={styles.main}>
      {/* Левая колонка сайдбара (ПК) */}
      <SideBar articlesList={articlesList} />

      {/* Правая колонка: Поиск и статья */}
      <div className={styles.articlePage}>
        {/* Левая колонка сайдбара (Mobile) */}
        <MobileSideBar
          drawerTitle="Каталог Статей"
          elements={articlesList}
          baseUrl="/articles"
          openButtonTitle="Выбрать статью из списка"
        />
        <Search articlesList={articlesList} />
        {children}
      </div>
    </main>
  );
};

export default ArticlesLayout;
