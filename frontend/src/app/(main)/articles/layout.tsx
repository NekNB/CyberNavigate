import "server-only";

import Search from "@/app/(main)/articles/_components/Search";

import MobileSideBar from "@/components/ui/MobileSideBar/MobileSideBar";
import { GetPublishedArticles } from "@/lib/api-server/Article/Article";
import { type IArticle } from "@/types/articles";
import styles from "./Articles.module.css";
import SideBar from "./_components/SideBar";
export const dynamic = "force-dynamic";
const ArticlesLayout = async ({ children }: { children: React.ReactNode }) => {
  const articlesList: IArticle[] = await GetPublishedArticles();
  return (
    <main className={styles.main}>
      {/* Левая колонка сайдбара (ПК) */}
      <SideBar articlesList={articlesList} />

      {/* Правая колонка: Поиск и статья */}
      <div className={styles.articlePage}>
        {/* Левая колонка сайдбара (Mobile) */}
        <MobileSideBar
          drawerTitle="Каталог Статей"
          elements={articlesList.map((article) => {
            article.id = article.slug;
            return article;
          })}
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
