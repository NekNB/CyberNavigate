// Footer.tsx
import type { FC } from "react";
import styles from "./Footer.module.css";

const Footer: FC = () => {
  return (
    <div className={styles.footer}>
      <h3 className={styles.title}>Наши контакты</h3>
      <a href="mailto:nikita.58baydin4@mail.ru" className={styles.contacts}>
        nikita.58baydin4@mail.ru
      </a>
      <a href="tel:+79960505076" className={styles.contacts}>
        +7-996-050-50-76
      </a>
    </div>
  );
};

export default Footer;
