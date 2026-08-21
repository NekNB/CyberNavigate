import type { FC } from "react";
import styles from "./Stats.module.css";
import graph from "/assets/graph-line.webp";

const Stats: FC = () => {
  return (
    <section className={styles.stats}>
      <div className={`${styles.container} ${styles.card}`}>
        <div>
          <img src={graph} className={styles.graphImg} />
        </div>

        <div>
          <h2 className={styles.title}>Атаки стали серьезнее</h2>
          <p className={styles.text}>
            Ущерб от кибератак в 2026 превысит 11,36 трлн. $ по данным
            PIR-Center — ведущей российской неправительственной организации,
            специализирующаяся на изучении вопросов глобальной безопасности.
            <br />C распространением технологий ИИ у злоумышленников появились
            новые методы и векторы атак
          </p>
        </div>
      </div>
    </section>
  );
};

export default Stats;
