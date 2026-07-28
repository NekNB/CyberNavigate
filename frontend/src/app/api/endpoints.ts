export const REFRESH_ENDPOINT = "/auth/refresh";

export const isRefreshEndpoint = (url: string | undefined): boolean => {
  if (!url) return false;

  return url.includes(REFRESH_ENDPOINT);
};

export const AUTH_ENDPOINTS = [REFRESH_ENDPOINT, "/auth/login"];

export const isAuthEndpoint = (url: string | undefined): boolean => {
  if (!url) return false;

  return AUTH_ENDPOINTS.some((endpoint) => {
    return url.includes(endpoint);
  });
};
