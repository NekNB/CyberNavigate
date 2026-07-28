import type { FC } from "react";
import Footer from "../../components/Footer/Footer";
import Header from "../../components/Header/Header";
import styles from "./Article.module.css";
const Article: FC = () => {
  return (
    <>
      <Header />
      <main className={styles.main}>
        <div className={styles.articleList}>
          <p>Тут будет список статей</p>
        </div>
        <div className={styles.articlePage}>
          <p>Тут будет статья</p>
        </div>
      </main>
      <Footer />
    </>
  );
};

export default Article;
