import type { IScenario } from "../../../types/simulator";
import { apiServer } from "../Api";

export const GetAllScenarios = async (): Promise<IScenario[]> => {
  try {
    const response = await apiServer.get<IScenario[]>("simulator/scenarios");
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const GetScenarioById = async (
  scenarioId: string,
): Promise<IScenario> => {
  try {
    const response = await apiServer.get<IScenario>(
      `simulator/scenarios/${scenarioId}`,
    );
    return response.data;
  } catch (error) {
    throw error;
  }
};
