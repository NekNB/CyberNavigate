import { GetAllScenarios } from "@/lib/api-server/Simulator/Simulator";
import { redirect } from "next/navigation";

export default async function SimulatorIndexPage() {
  const scenarios = await GetAllScenarios();

  if (!scenarios || scenarios.length === 0) {
    return <div>Статьи не найдены</div>;
  }

  // Перенаправляем на ID первой статьи до того, как клиенту отдается HTML
  redirect(`/simulator/${scenarios[0].id}`);
}
