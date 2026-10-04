import type {
  IFile,
  IResults,
  IScenario,
  IStep,
} from "../../../types/simulator";
import apiClient from "../Api";

export const GetAllScenarios = async (): Promise<IScenario[]> => {
  try {
    const response = await apiClient.get<IScenario[]>("simulator/scenarios");
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const CreateSession = async (scenarioId: string): Promise<void> => {
  await apiClient.post("/simulator/sessions", {
    scenarioId: scenarioId,
  });
};

export const GetStep = async (): Promise<{ code: number; step?: IStep }> => {
  const response = await apiClient.get<IStep | void>("/simulator/step");

  switch (response.status) {
    case 200:
      return {
        code: response.status,
        step: response.data as IStep,
      };
    case 204:
      return { code: response.status };
    default:
      throw new Error(`Status code unknown: ${response.status}`);
  }
};

export const GetFile = async (fileId: string): Promise<IFile> => {
  return (await apiClient.get<IFile>(`/simulator/files/${fileId}`)).data;
};
export const SendAnswer = async (answerId: string): Promise<void> => {
  return await apiClient.post("/simulator/action/answer", {
    answerId: answerId,
  });
};

export const GetResults = async (): Promise<IResults> => {
  const results = await apiClient.delete<IResults>("/simulator/sessions");
  return results.data;
};
