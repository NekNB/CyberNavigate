import type { FC } from "react";
import styles from "./Hero.module.css";
const Hero: FC = () => {
  return (
    <section className={styles.hero}>
      <div className={styles.container}>
        <h1 className={styles.title}>
          Хакеры не дремлют! <br />
          <span className={styles.DesktopOnly}>
            Больше знаешь - крепче спишь
          </span>
        </h1>
        <p className={styles.text}>
          «КиберНавигатор» — проект по обучению молодежи кибербезопасности
          <span className={styles.DesktopOnly}>
            : он помогает распознавать онлайн-угрозы, защищать личные данные и
            противостоять мошенникам через простое обучение, практические
            симуляции и актуальную базу знаний
          </span>
        </p>
      </div>
    </section>
  );
};

export default Hero;
