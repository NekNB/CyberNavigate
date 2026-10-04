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
        const id = item.id;
        const name = item.title;
        const isActive = path?.includes(id);

        return (
          <Link
            key={id}
            className={`${styles.articleItem} ${isActive ? styles.active : ""}`}
            href={`/articles/${item.id}`}
          >
            {name}
          </Link>
        );
      })}
    </div>
  );
};

export default SideBar;
