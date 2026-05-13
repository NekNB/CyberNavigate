import Header from "../header/Header";
import Footer from "../footer/Footer";
import "./Frame.css";
import type { FC } from "react";

const Frame: FC = () => {
  return (
    <div className="frame">
      <div className="view">
        <Header />

        <div className="rectangle" />

        <aside className="sidebar" />

        <main className="article">
          <div className="article-title">Название статьи</div>
          <div className="article-text">Текст статьи</div>
        </main>

        <Footer />
      </div>
    </div>
  );
};

export default Frame;