import { AxiosError } from "axios";
import type { IUser } from "../../../types/users";
import apiClient from "../Api";

export const GetUser = async (): Promise<IUser> => {
  try {
    const response = await apiClient.get<IUser>("/users/me");
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const RegisterUser = async (
  username: string,
  password: string,
): Promise<number> => {
  try {
    const response = await apiClient.post("/users", {
      username: username,
      password: password,
    });
    return response.status;
  } catch (error) {
    if (error instanceof AxiosError && error.response?.status === 400) {
      return 400;
    }
    throw error;
  }
};
