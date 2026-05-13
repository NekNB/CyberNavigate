import "./Header.css";
import type { FC } from "react";

const Header: FC = () => {
  return (
    <header className="header">
      <div className="header-background-wrapper">
        <div className="header-background">
          <img
            className="header-logo"
            alt="Логотип"
            src="/icons/logo.svg"
          />
        </div>
      </div>
      <div className="header-nav-background" />
      <nav className="header-nav">
        <div className="header-nav-item">Главная</div>
        <div className="header-nav-item">Статьи</div>
        <div className="header-nav-item">Симулятор</div>
      </nav>
    </header>
  );
};

export default Header;