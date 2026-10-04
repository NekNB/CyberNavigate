import { GetAllScenarios } from "@/lib/api-server/Simulator/Simulator";
import SideBar from "./_components/SideBar";
import styles from "./Simulator.module.css";
export default async function SimulatorLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  // Получаем список сценариев
  const scenarios = await GetAllScenarios();

  return (
    <section className={styles.simulator}>
      <SideBar scenariosList={scenarios} />
      {children}
    </section>
  );
}
