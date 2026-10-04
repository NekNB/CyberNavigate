// components/ClientConfigProvider.tsx
"use client";

import { setClientApiBaseUrl } from "@/lib/api-client/Api";
import React, { createContext } from "react";

interface ConfigContextType {
  clientApiUrl: string;
}

const ConfigContext = createContext<ConfigContextType | null>(null);

export function ClientConfigProvider({
  clientApiUrl,
  children,
}: {
  clientApiUrl: string;
  children: React.ReactNode;
}) {
  setClientApiBaseUrl(clientApiUrl);
  return (
    <ConfigContext.Provider value={{ clientApiUrl }}>
      {children}
    </ConfigContext.Provider>
  );
}
