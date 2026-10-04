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
          <Link
            href="/articles/chto-takoe-fishing-i-kak-ne-popast-sia-na-udochku-moshennikov"
            className={styles.card}
          >
            <h3 className={styles.title}>Фишинг</h3>
            <p className={styles.text}>
              Обман через письма и сообщения для кражи ваших данных
            </p>
          </Link>

          <Link
            href="/articles/nastroiki-privatnosti-v-sotssetiakh-za-5-minut"
            className={styles.card}
          >
            <h3 className={styles.title}>Соц Сети</h3>
            <p className={styles.text}>
              Данные в соц сетях могут быть использованы против Вас
            </p>
          </Link>

          <Link
            href="/articles/bezopasnost-v-publichnykh-wi-fi-setiakh"
            className={styles.card}
          >
            <h3 className={styles.title}>Бесплатный Wi-Fi</h3>
            <p className={styles.text}>
              Почему кто-то бесплатно делится с Вами интернетом?
            </p>
          </Link>

          <Link
            href="/articles/meropriiatiia-po-informatsionnoi-gigiene-dlia-studentov-kak-provodit-i-chto-osveshchat"
            className={styles.card}
          >
            <h3 className={styles.title}>Методические рекомендации</h3>
            <p className={styles.text}>
              Все о проведении мероприятий для студентов
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
