import type { FC } from "react";
import styles from "./Hero.module.css";
const Hero: FC = () => {
  return (
    <section className={styles.hero}>
      <div className={styles.container}>
        <h1 className={styles.title}>
          Они не дремлют! <br />
          Больше знаешь - крепче спишь
        </h1>
        <p className={styles.text}>
          «КиберНавигатор» — проект по обучению молодежи кибербезопасности: он
          помогает распознавать онлайн-угрозы, защищать личные данные и
          противостоять мошенникам через простое обучение, практические
          симуляции и актуальную базу знаний
        </p>
      </div>
    </section>
  );
};

export default Hero;
