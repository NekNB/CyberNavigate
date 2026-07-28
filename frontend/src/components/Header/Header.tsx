import { useEffect, useState, type FC } from "react";
import { GetUser } from "../../app/api/User/User";
import Login from "../Login/Login";

import { Logout } from "../../app/api/Auth/Auth";
import Register from "../Register/Register";
import styles from "./Header.module.css";
import logo from "/assets/logo.png";
import openDoor from "/assets/open-door.svg";
import profile from "/assets/profile.svg";
const Header: FC = () => {
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
          "Вы уверены, что хотите выйти из профиля? Тогда ваш прогресс не будет сохранен",
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

  return (
    <header className={styles.header}>
      {isLoginOpen ? (
        <Login
          onClose={() => {
            setIsLoginOpen(false);
          }}
          openRegister={() => {
            setIsLoginOpen(false);
            setIsRegisterOpen(true);
          }}
        />
      ) : null}
      {isRegisterOpen ? (
        <Register
          onClose={() => {
            setIsRegisterOpen(false);
          }}
          openLogin={() => {
            setIsRegisterOpen(false);
            setIsLoginOpen(true);
          }}
        />
      ) : null}
      <a href="/" className={styles.logo}>
        <img src={logo} />
      </a>
      <nav className={styles.nav}>
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
      <img
        className={styles.profile}
        onClick={openProfile}
        src={isAuth ? profile : openDoor}
      />
    </header>
  );
};

export default Header;
