import MobileSideBar from "@/components/ui/MobileSideBar/MobileSideBar";
import {
  GetAllScenarios,
  GetScenarioById,
} from "@/lib/api-server/Simulator/Simulator";
import { Metadata } from "next";
import Link from "next/link";
import { FC } from "react";
import styles from "./Scenario.module.css";
export const dynamic = "force-dynamic";
interface PageProps {
  params: Promise<{ simulatorId: string }>;
}

export async function generateMetadata({
  params,
}: PageProps): Promise<Metadata> {
  const scenario = await GetScenarioById((await params).simulatorId);
  return {
    title: scenario.title,
    description: scenario.description,
  };
}
const SimulatorPage: FC<PageProps> = async ({ params }) => {
  const { simulatorId } = await params;
  const scenario = await GetScenarioById(simulatorId);
  const scenarioList = await GetAllScenarios();
  return (
    <div className={styles.scenarioPage}>
      <MobileSideBar
        baseUrl="/simulator"
        drawerTitle="Список сценариев"
        elements={scenarioList}
        openButtonTitle="Выбрать сценарий из списка"
      />
      <h3 className={styles.title}>{scenario.title}</h3>
      <p className={`${styles.text} ${styles.difficulty}`}>
        <b>Сложность: </b>
        {scenario.difficulty}
      </p>
      <p className={`${styles.text} ${styles.description}`}>
        <b>Описание: </b>
        {scenario.description}
      </p>
      <Link
        className={`${styles.buttonChoiceScenario}`}
        href={`/scenario/${scenario.id}`}
      >
        Выбрать сценарий
      </Link>
    </div>
  );
};
export default SimulatorPage;
