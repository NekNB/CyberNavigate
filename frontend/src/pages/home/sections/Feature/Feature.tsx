import type { FC } from "react";
import styles from "./Feature.module.css";
const Feature: FC = () => {
  return (
    <>
      <div className={styles.container}>
        <p className={styles.featuresText}>
          Большинство угроз в кибер пространстве можно избежать лишь сохраняя
          внимательность и предосторожность
        </p>
        <div className={styles.cardsGrid}>
          <div className={styles.card}>
            <h3 className={styles.title}>Фишинг</h3>
            <p className={styles.text}>
              Обман через письма и сообщения для кражи ваших данных
            </p>
          </div>

          <div className={styles.card}>
            <h3 className={styles.title}>Социальная инженерия</h3>
            <p className={styles.text}>
              Психологическое давление по телефону для выманивания денег жертвы
            </p>
          </div>

          <div className={styles.card}>
            <h3 className={styles.title}>Мошенничество</h3>
            <p className={styles.text}>
              Обман ради выгоды: фейковые магазины, инвестиции и пирамиды
            </p>
          </div>

          <div className={styles.card}>
            <h3 className={styles.title}>Технические атаки</h3>
            <p className={styles.text}>
              Внедрение вредоносных программ для кражи данных и блокировки
              устройств
            </p>
          </div>
        </div>

        <div className={styles.cardsButton}>
          <a href="#" className={styles.buttonMore}>
            Больше наших статей
          </a>
        </div>
      </div>
    </>
  );
};

export default Feature;
