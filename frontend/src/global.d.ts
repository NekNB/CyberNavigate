// src/global.d.ts

import type { AppConfig } from "./types/config";

// Описываем структуру нашего YAML файла

// Объявляем глобальную переменную
declare global {
  const __APP_CONFIG__: AppConfig;
}

export {};
