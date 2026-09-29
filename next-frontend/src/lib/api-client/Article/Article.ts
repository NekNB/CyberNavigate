// import { EArticleStatus, type IArticle } from "../../../types/articles";
// import apiClient from "../Api";

// export const GetPublishedArticles = async (): Promise<IArticle[]> => {
//   const response = await apiClient.get<IArticle[]>("/articles");
//   return response.data.filter((article) => {
//     if (article.status === EArticleStatus.PUBLISHED) {
//       return article;
//     }
//   });
// };
// export const GetArticleText = async (articleId: string): Promise<string> => {
//   const response = await apiClient.get<string>(`/articles/${articleId}/text`);
//   return response.data;
// };

import { EArticleStatus, type IArticle } from "../../../types/articles";
import apiClient from "../Api";

export const GetPublishedArticles = async (): Promise<IArticle[]> => {
  const response = await apiClient.get<any>("/articles");
  
  // Достаем массив из ответа бэкенда
  const rawData = response.data;
  const articlesArray: IArticle[] = Array.isArray(rawData)
    ? rawData
    : rawData?.articles || rawData?.data || rawData?.items || [];

  return articlesArray.filter(
    (article) => article.status === EArticleStatus.PUBLISHED
  );
};

export const GetArticleText = async (articleId: string): Promise<string> => {
  const response = await apiClient.get<string>(`/articles/${articleId}/text`);
  return response.data;
};