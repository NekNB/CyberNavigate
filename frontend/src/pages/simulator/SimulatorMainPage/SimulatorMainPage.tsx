import { useEffect, useState, type FC } from "react";
import { GetAllScenarios } from "../../../app/api/Simulator/Simulator";
import type { IScenario } from "../../../types/simulator";
import ScenarioPage from "./ScenarioPage/ScearioPage";
import styles from "./SimulatorMainPage.module.css";

const SimulatorMainPage: FC = () => {
  const [scenarios, setScenarios] = useState<IScenario[]>([]);
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const [currentScenario, setCurrentScenario] = useState<IScenario>();

  useEffect(() => {
    const fetch = async () => {
      try {
        const scenarios = await GetAllScenarios();
        setScenarios(scenarios);
        if (!scenarios?.length) {
          setError("Сценарии не найдены");
          return;
        }
        setCurrentScenario(scenarios[0]);
      } catch (error: any) {
        setError(error);
      } finally {
        setIsLoading(false);
      }
    };

    fetch();
  }, []);

  return isLoading ? (
    <p>Загрузка сценариев</p>
  ) : error ? (
    <p>Произошла ошибка при получении данных</p>
  ) : (
    <div className={styles.menu}>
      <nav className={styles.scenariosList}>
        {scenarios?.length &&
          scenarios.map((scenario) => {
            return (
              <div
                className={styles.scenarioTitle}
                key={scenario.id}
                onClick={() => {
                  setCurrentScenario(scenario);
                }}
              >
                {scenario.title}
              </div>
            );
          })}
      </nav>
      <main className={styles.mainContent}>
        {currentScenario && <ScenarioPage scenario={currentScenario} />}
      </main>
    </div>
  );
};

export default SimulatorMainPage;
