import Link from "next/link";
import type { FC } from "react";
import styles from "./Feature.module.css";
const Feature: FC = () => {
  return (
    <section className={styles.feature}>
      <div className={styles.container}>
        <h3 className={styles.featuresText}>
          Полезные статьи
          <span className={styles.DesktopOnly}>
            , которые помогут Вам избежать большинства угроз в современном мире
          </span>
        </h3>
        <div className={styles.cardsGrid}>
          <Link href="/articles" className={styles.card}>
            <h3 className={styles.title}>Фишинг</h3>
            <p className={styles.text}>
              Обман через письма и сообщения для кражи ваших данных
            </p>
          </Link>

          <Link href="/articles" className={styles.card}>
            <h3 className={styles.title}>Социальная инженерия</h3>
            <p className={styles.text}>
              Психологическое давление по телефону для выманивания денег жертвы
            </p>
          </Link>

          <Link href="/articles" className={styles.card}>
            <h3 className={styles.title}>Мошенничество</h3>
            <p className={styles.text}>
              Обман ради выгоды: фейковые магазины, инвестиции и пирамиды
            </p>
          </Link>

          <Link href="/articles" className={styles.card}>
            <h3 className={styles.title}>Технические атаки</h3>
            <p className={styles.text}>
              Внедрение вредоносных программ для кражи данных и блокировки
              устройств
            </p>
          </Link>
        </div>

        <div className={styles.cardsButton}>
          <Link href="/articles" className={styles.buttonMore}>
            Больше наших статей
          </Link>
        </div>
      </div>
    </section>
  );
};

export default Feature;
