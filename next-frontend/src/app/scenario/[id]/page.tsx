import Scenario from "@/pages/scenario/Scenario";

type Props = {
  params: Promise<{ id: string }>;
};

export default async function ScenarioPage({ params }: Props) {
  // Если у вас Next.js 15+, используйте: const { id } = await params;
  const { id } = await params;

  // Вы можете передать id в ваш существующий компонент
  return <Scenario id={id} />;
}
