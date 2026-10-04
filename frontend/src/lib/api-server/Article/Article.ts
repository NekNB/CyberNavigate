import "server-only";

import { EArticleStatus, type IArticle } from "../../../types/articles";
import { apiServer } from "../Api";

export const GetPublishedArticles = async (): Promise<IArticle[]> => {
  const response = await apiServer.get<IArticle[]>("/articles");

  // Достаем массив из ответа бэкенда
  const articles = response.data;

  return articles.filter(
    (article) => article.status === EArticleStatus.PUBLISHED,
  );
};

export const GetArticleText = async (articleId: string): Promise<string> => {
  const response = await apiServer.get<string>(`/articles/${articleId}/text`);
  return response.data;
};

export const GetArticleByIdOrSlug = async (
  articleIdOrSlug: string,
): Promise<IArticle> => {
  const response = await apiServer.get<IArticle>(
    `/articles/${articleIdOrSlug}`,
  );
  return response.data;
};
