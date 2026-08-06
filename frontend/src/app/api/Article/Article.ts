import { EArticleStatus, type IArticle } from "../../../types/articles";
import apiClient from "../Api";

export const GetPublishedArticles = async (): Promise<IArticle[]> => {
  const response = await apiClient.get<IArticle[]>("/articles");
  return response.data.filter((article) => {
    if (article.status === EArticleStatus.PUBLISHED) {
      return article;
    }
  });
};
export const GetArticleText = async (articleId: string): Promise<string> => {
  const response = await apiClient.get<string>(`/articles/${articleId}/text`);
  return response.data;
};
