export const EArticleStatus = {
  PUBLISHED: "published",
  DRAFT: "draft",
  ARCHIVE: "archive",
};
export type EArticleStatus =
  (typeof EArticleStatus)[keyof typeof EArticleStatus];

export interface IArticle {
  id: string;
  title: string;
  status: EArticleStatus;
}
