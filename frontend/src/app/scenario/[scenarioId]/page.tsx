import ScenarioClient from "./ScenarioClient";

interface PageProps {
  params: Promise<{ scenarioId: string }>;
}
export default async function ScenarioPage({ params }: PageProps) {
  const { scenarioId } = await params;

  return <ScenarioClient scenarioId={scenarioId} />;
}
