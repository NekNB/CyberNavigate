import { useEffect, useState, type FC } from "react";
import { GetUser } from "../../app/api/User/User";
import Login from "../Login/Login";
import { Logout } from "../../app/api/Auth/Auth";
import Register from "../Register/Register";
import styles from "./Header.module.css";
import logo from "/assets/logo.png";
import openDoor from "/assets/open-door.svg";
import profile from "/assets/profile.svg";

interface ArticleItem {
  id: string | number;
  title: string;
}

interface HeaderProps {
  articles?: ArticleItem[];
  activeArticleId?: string | number;
  onSelectArticle?: (id: string | number) => void;
  isMenuOpen?: boolean;
  setIsMenuOpen?: (open: boolean) => void;
}

const Header: FC<HeaderProps> = ({
  articles = [],
  activeArticleId,
  onSelectArticle,
  isMenuOpen = false,
  setIsMenuOpen,
}) => {
  const [isAuth, setIsAuth] = useState(false);
  const [isLoginOpen, setIsLoginOpen] = useState(false);
  const [isRegisterOpen, setIsRegisterOpen] = useState(false);

  useEffect(() => {
    const checkAuth = async () => {
      try {
        await GetUser();
        setIsAuth(true);
      } catch (error) {
        setIsAuth(false);
      }
    };
    checkAuth();
  }, [isLoginOpen]);

  const openProfile = async () => {
    if (isAuth) {
      if (
        window.confirm(
          "Вы уверены, что хотите выйти из профиля? Тогда ваш прогресс не будет сохранен"
        )
      ) {
        await Logout();
        setIsAuth(false);
        setIsLoginOpen(true);
      }
    } else {
      setIsLoginOpen(true);
    }
  };

  const closeMenu = () => {
    if (setIsMenuOpen) setIsMenuOpen(false);
  };

  return (
    <header className={styles.header}>
      {isLoginOpen && (
        <Login
          onClose={() => setIsLoginOpen(false)}
          openRegister={() => {
            setIsLoginOpen(false);
            setIsRegisterOpen(true);
          }}
        />
      )}
      {isRegisterOpen && (
        <Register
          onClose={() => setIsRegisterOpen(false)}
          openLogin={() => {
            setIsRegisterOpen(false);
            setIsLoginOpen(true);
          }}
        />
      )}

      <div className={styles.leftSection}>
        <a href="/" className={styles.logo}>
          <img src={logo} alt="КиберНавигатор" />
        </a>

        <nav className={styles.navMainLinks}>
          <a className={styles.headerLink} href="/">
            Главная
          </a>
          <a className={styles.headerLink} href="/articles">
            Статьи
          </a>
          <a className={styles.headerLink} href="/simulator">
            Симулятор
          </a>
        </nav>
      </div>

      <div className={styles.rightSection}>
        <img
          className={styles.profile}
          onClick={openProfile}
          src={isAuth ? profile : openDoor}
          alt="Профиль"
        />
      </div>

      <aside className={`${styles.articleDrawer} ${isMenuOpen ? styles.drawerActive : ""}`}>
        <div className={styles.drawerHeader}>
          <span className={styles.drawerTitle}>Каталог статей</span>
          <button className={styles.closeBtn} onClick={closeMenu} aria-label="Закрыть">
            ✕
          </button>
        </div>

        <ul className={styles.mobileArticlesList}>
          {articles.map((art) => (
            <li
              key={art.id}
              className={`${styles.mobileArticleItem} ${
                art.id === activeArticleId ? styles.mobileArticleActive : ""
              }`}
              onClick={() => {
                if (onSelectArticle) onSelectArticle(art.id);
                closeMenu();
              }}
            >
              {art.title}
            </li>
          ))}
        </ul>
      </aside>

      {isMenuOpen && <div className={styles.overlay} onClick={closeMenu} />}
    </header>
  );
};

export default Header;
