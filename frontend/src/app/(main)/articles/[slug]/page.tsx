import {
  GetArticleByIdOrSlug,
  GetArticleText,
} from "@/lib/api-server/Article/Article";
import DOMPurify from "isomorphic-dompurify";
import styles from "./Article.module.css";
type PageProps = {
  params: Promise<{ slug: string }>;
};
export const dynamic = "force-dynamic";
export default async function Article({ params }: PageProps) {
  const { slug } = await params;

  //  Загрузка текста выбранной статьи
  const articleData = await GetArticleByIdOrSlug(slug);
  const articleText = await GetArticleText(articleData.id);
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
