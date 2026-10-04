import { GetPublishedArticles } from "@/lib/api-server/Article/Article";
import { redirect } from "next/navigation";

export default async function ArticlesIndexPage() {
  const articles = await GetPublishedArticles();

  if (!articles || articles.length === 0) {
    return <div>Статьи не найдены</div>;
  }

  // Перенаправляем на ID первой статьи до того, как клиенту отдается HTML
  redirect(`/articles/${articles[0]?.slug}`);
}
