import {
  GetArticleById,
  GetArticleText,
} from "@/lib/api-server/Article/Article";
import DOMPurify from "isomorphic-dompurify";
import styles from "./Article.module.css";
type PageProps = {
  params: Promise<{ slug: string }>;
};

export default async function Article({ params }: PageProps) {
  const { slug } = await params;

  //  Загрузка текста выбранной статьи
  const articleText = await GetArticleText(articleId);
  const articleData = await GetArticleById(articleId);
  const clearHTML = DOMPurify.sanitize(articleText);

  return (
    <div className={styles.articleCard}>
      <h1 className={styles.articleTitle}>{articleData.title}</h1>
      <div className={styles.articleText}>
        <div
          style={{ whiteSpace: "pre-line" }}
          dangerouslySetInnerHTML={{ __html: clearHTML }}
        ></div>
      </div>
    </div>
  );
}
