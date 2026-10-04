"use client";

import { IArticle } from "@/types/articles";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { FC } from "react";
import styles from "./SideBar.module.css";
interface Props {
  articlesList: IArticle[];
}

const SideBar: FC<Props> = ({ articlesList }) => {
  const path = usePathname();
  return (
    <div className={styles.articleList}>
      {articlesList.map((item) => {
        const slug = item.slug;
        const name = item.title;
        const isActive = path?.includes(slug);

        return (
          <Link
            key={item.id}
            className={`${styles.articleItem} ${isActive ? styles.active : ""}`}
            href={`/articles/${item.slug}`}
          >
            {name}
          </Link>
        );
      })}
    </div>
  );
};

export default SideBar;
