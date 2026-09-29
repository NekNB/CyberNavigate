import { AxiosError } from "axios";
import apiClient from "../Api";

export const Login = async (
  username: string,
  password: string,
): Promise<number> => {
  try {
    const response = await apiClient.post("/auth/login", {
      username: username,
      password: password,
    });
    return response.status;
  } catch (error) {
    if (error instanceof AxiosError && error.response?.status === 401) {
      return 401;
    }
    if (error instanceof AxiosError) {
      console.log(error.response?.status);
    }
    throw error;
  }
};

export const Logout = async (): Promise<number> => {
  try {
    const response = await apiClient.delete("/auth/logout");
    return response.status;
  } catch (error) {
    throw error;
  }
};
