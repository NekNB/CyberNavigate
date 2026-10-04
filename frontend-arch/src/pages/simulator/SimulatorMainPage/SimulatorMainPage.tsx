import { useEffect, useState, type FC } from "react";
import { GetAllScenarios } from "../../../app/api/Simulator/Simulator";
import type { IScenario } from "../../../types/simulator";
import ScenarioPage from "./ScenarioPage/ScearioPage";
import Header from "../../../components/Header/Header";
import Footer from "../../../components/Footer/Footer";
import styles from "./SimulatorMainPage.module.css";

const SimulatorMainPage: FC = () => {
  const [scenarios, setScenarios] = useState<IScenario[]>([]);
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const [currentScenario, setCurrentScenario] = useState<IScenario>();

  // Стейт для управления открытием шторки
  const [isMenuOpen, setIsMenuOpen] = useState<boolean>(false);

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

  // Мапим сценарии под формат, который ждёт Header ({ id, title })
  const headerScenarios = (scenarios || []).map((sc) => ({
    id: sc.id,
    title: sc.title,
  }));

  return (
    <>
      {/* 1. Наш готовый Header со шторкой сценариев */}
      <Header
        articles={headerScenarios}
        activeArticleId={currentScenario?.id}
        onSelectArticle={(id) => {
          const selected = scenarios.find((s) => s.id === id);
          if (selected) setCurrentScenario(selected);
        }}
        isMenuOpen={isMenuOpen}
        setIsMenuOpen={setIsMenuOpen}
        drawerTitle="Список сценариев"
      />

      <main>
        {isLoading ? (
          <p style={{ textAlign: "center", padding: "40px", color: "#fff" }}>
            Загрузка сценариев...
          </p>
        ) : error ? (
          <p style={{ textAlign: "center", padding: "40px", color: "#ff8888" }}>
            Произошла ошибка при получении данных
          </p>
        ) : (
          <div className={styles.menu}>
            {/* 2. Кнопка вызова шторки для мобилок */}
            <button
              className={styles.mobileMenuTrigger}
              onClick={() => setIsMenuOpen(true)}
            >
              ☰ Выбрать сценарий из списка
            </button>

            {/* Левая колонка со сценариями (ПК) */}
            <nav className={styles.scenariosList}>
              {scenarios?.length &&
                scenarios.map((scenario) => {
                  const isActive = currentScenario?.id === scenario.id;
                  return (
                    <div
                      className={`${styles.scenarioTitle} ${isActive ? styles.active : ''}`}
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

            {/* Контент выбранного сценария */}
            <div className={styles.mainContent}>
              {currentScenario && <ScenarioPage scenario={currentScenario} />}
            </div>
          </div>
        )}
      </main>

      <Footer />
    </>
  );
};

export default SimulatorMainPage;