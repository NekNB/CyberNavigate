import axios, { AxiosError, type InternalAxiosRequestConfig } from "axios";

import { isAuthEndpoint, isRefreshEndpoint } from "./endpoints";

// Очередь обновления
let isRefreshing = false;
let failedQueue: Array<{
  resolve: (value?: unknown) => void;
  reject: (error: any) => void;
}> = [];

const processQueue = (error: any) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve();
    }
  });
  failedQueue = [];
};

const apiClient = axios.create({
  baseURL: __APP_CONFIG__.backend.url,
  timeout: 10000,
  headers: {
    "Content-Type": "application/json",
  },
  withCredentials: true,
});

apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & {
      _retry?: boolean;
    };

    if (isAuthEndpoint(originalRequest.url) && error.response?.status === 401) {
      return Promise.reject(error);
    }

    if (
      isRefreshEndpoint(originalRequest.url) &&
      error.response?.status === 422
    ) {
      return Promise.reject(error);
    }

    // защита от рекурсии. Если уже была попытка и она не увенчалась успехом
    if (error.response?.status === 401 && !originalRequest._retry) {
      // Если сам запрос на обновление токенов провалился => logout
      if (isRefreshEndpoint(originalRequest.url)) {
        return Promise.reject(error);
      }

      // Если уже происходит процесс обновления -> добавляем в очередь
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        }).then(() => {
          return apiClient(originalRequest);
        });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        await apiClient.put(`${__APP_CONFIG__.backend.url}/auth/refresh`);

        processQueue(null);

        return apiClient(originalRequest);
      } catch (refreshError) {
        processQueue(refreshError);

        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }
    return Promise.reject(error);
  },
);

export default apiClient;
