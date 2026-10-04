import type { FC } from "react";
import { useNavigate } from "react-router";
import type { IScenario } from "../../../../types/simulator";
import styles from "./ScenarioPage.module.css";
interface ScenarioProps {
  scenario: IScenario;
}

const ScenarioPage: FC<ScenarioProps> = ({ scenario }) => {
  const navigate = useNavigate();

  return (
    <div className={styles.scenarioPage}>
      <h3 className={styles.title}>{scenario.title}</h3>
      <p className={`${styles.text} ${styles.difficulty}`}>
        <b>Сложность: </b>
        {scenario.difficulty}
      </p>
      <p className={`${styles.text} ${styles.description}`}>
        <b>Описание: </b>
        {scenario.description}
      </p>
      <button
        className={`${styles.buttonChoiceScenario}`}
        onClick={() => navigate(`/scenario/${scenario.id}`)}
      >
        Выбрать сценарий
      </button>
    </div>
  );
};

export default ScenarioPage;
