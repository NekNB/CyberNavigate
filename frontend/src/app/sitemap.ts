import { MetadataRoute } from "next";

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
  const dynamicSlugs = ["kak-zashchitit-dannye", "chto-takoe-fishing"];

  const articlePages = dynamicSlugs.map((slug) => ({
    url: `${baseUrl}/articles/${slug}`,
    lastModified: new Date(),
    changeFrequency: "weekly" as const,
    priority: 0.6,
  }));

  return [...staticPages, ...articlePages];
}
