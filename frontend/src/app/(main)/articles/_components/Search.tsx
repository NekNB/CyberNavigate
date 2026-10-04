"use client";

import { IArticle } from "@/types/articles";
import { useRouter } from "next/navigation";
import { FC, useEffect, useRef, useState } from "react";
import styles from "./Search.module.css";

interface Props {
  articlesList: IArticle[];
}

const Search: FC<Props> = ({ articlesList: articleList }) => {
  const [searchQuery, setSearchQuery] = useState("");
  const router = useRouter();
  // Двухэтапная фильтрация для поиска
  const searchLower = searchQuery.trim().toLowerCase();
  const searchRef = useRef<HTMLDivElement>(null);
  // 1. Совпадения по НАЗВАНИЮ
  const matchesByTitle = articleList.filter((item) => {
    const name = (item.title || item.title || "").toLowerCase();
    return searchLower !== "" && name.includes(searchLower);
  });

  const searchResultsDropdown = searchQuery.trim() !== "";
  // 🎯 Закрываем выпадающий список при клике вне компонента поиска
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        searchRef.current &&
        !searchRef.current.contains(event.target as Node)
      ) {
        setSearchQuery(""); // Сбрасываем/закрываем поиск
      }
    };

    // Слушаем клик на всем документе
    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);
  return (
    <div ref={searchRef} className={styles.searchBox}>
      <div className={styles.searchInputWrapper}>
        <input
          type="text"
          className={styles.searchInput}
          placeholder="Поиск по статьям..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />

        {/* Выпадающий список совпадений с разбивкой по типам */}
        {searchResultsDropdown && (
          <div className={styles.searchResultsDropdown}>
            {matchesByTitle.length === 0 ? (
              <div className={styles.searchNoResults}>
                Ничего не найдено по запросу «{searchQuery}»
              </div>
            ) : (
              <>
                {/* Блок 1: Совпадения по названию */}
                {matchesByTitle.length > 0 && (
                  <div className={styles.searchCategoryGroup}>
                    {matchesByTitle.map((item) => {
                      const id = item.id;
                      const name = item.title;
                      return (
                        <div
                          key={`title-${id}`}
                          className={styles.searchResultItem}
                          onClick={() => {
                            if (id !== undefined)
                              router.push(`/articles/${item.id}`);
                            setSearchQuery("");
                          }}
                        >
                          🔍 {name}
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
  );
};

export default Search;
