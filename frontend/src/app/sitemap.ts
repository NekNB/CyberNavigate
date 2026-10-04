import { GetPublishedArticles } from "@/lib/api-server/Article/Article";
import { GetAllScenarios } from "@/lib/api-server/Simulator/Simulator";
import { MetadataRoute } from "next";
export const dynamic = "force-dynamic";
export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const baseUrl = "https://кибер-навигатор.рф";

  // 1. Статические страницы
  const staticPages = [
    {
      url: baseUrl,
      lastModified: new Date(),
      changeFrequency: "weekly" as const,
      priority: 1, // Главная страница — самый высокий приоритет
    },
    {
      url: `${baseUrl}/articles`,
      lastModified: new Date(),
      changeFrequency: "weekly" as const,
      priority: 0.8,
    },
  ];

  // 2. Динамические страницы (например, статьи)
  // Получаем массив slug'ов из вашей базы или API
  // const articles = await fetchArticles();
  const articles = await GetPublishedArticles();
  const dynamicArticleSlugs = articles.map((article) => {
    return article.slug;
  });
  const articlePages: MetadataRoute.Sitemap = dynamicArticleSlugs.map(
    (slug) => ({
      url: `${baseUrl}/articles/${slug}`,
      lastModified: new Date(),
      changeFrequency: "weekly" as const,
      priority: 0.6,
    }),
  );

  const scenarios = await GetAllScenarios();
  const dynamicScenarioUUIDs = scenarios.map((scenario) => {
    return scenario.id;
  });
  const scenarioPage: MetadataRoute.Sitemap = dynamicScenarioUUIDs.map(
    (id) => ({
      url: `${baseUrl}/simulator/${id}`,
      lastModified: new Date(),
      changeFrequency: "monthly" as const,
      priority: 0.5,
    }),
  );

  return [...staticPages, ...articlePages, ...scenarioPage];
}
