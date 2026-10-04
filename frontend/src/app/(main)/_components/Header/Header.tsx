"use client";

import { useEffect, useState, type FC } from "react";

import { Logout } from "@/lib/api-client/Auth/Auth";
import { GetUser } from "@/lib/api-client/User/User";

import Login from "@/components/features/Login/Login";
import Register from "@/components/features/Register/Register";
import styles from "./Header.module.css";

const Header: FC = () => {
  const [isAuth, setIsAuth] = useState(false);
  const [isLoginOpen, setIsLoginOpen] = useState(false);
  const [isRegisterOpen, setIsRegisterOpen] = useState(false);

  useEffect(() => {
    const checkAuth = async () => {
      try {
        await GetUser();
        setIsAuth(true);
      } catch {
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

  const [isVisible, setIsVisible] = useState(true);
  const [lastScrollY, setLastScrollY] = useState(0);
  useEffect(() => {
    const handleScroll = () => {
      const currentScrollY = window.scrollY;
      if (currentScrollY > lastScrollY && currentScrollY > 80) {
        setIsVisible(false);
      } else {
        setIsVisible(true);
      }

      setLastScrollY(currentScrollY);
    };
    window.addEventListener("scroll", handleScroll);
  }, [lastScrollY]);

  return (
    <header className={`${styles.header} ${!isVisible ? styles.hidden : ""}`}>
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
          <img src="/assets/logo.webp" alt="КиберНавигатор" />
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
          src={isAuth ? "/assets/profile.svg" : "/assets/open-door.svg"}
          alt="Профиль"
        />
      </div>
    </header>
  );
};

export default Header;
