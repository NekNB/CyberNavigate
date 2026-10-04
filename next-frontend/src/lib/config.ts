// src/lib/config.ts
import { readFileSync } from "fs";
import "server-only"; // Защита от случайного вызова на клиенте

import { load as yamlLoad } from "js-yaml";
import path from "path";

export interface AppConfig {
  internalUrl: string;
  externalUrl: string;
}

export function loadConfig(): AppConfig {
  // 1. Читаем путь из переменной окружения (или fallback на config.yaml)
  const envPath = process.env.CONFIG_PATH || "config.yaml";

  // 2. Преобразуем в абсолютный путь относительно корня проекта
  const configPath = path.isAbsolute(envPath)
    ? envPath
    : path.resolve(process.cwd(), envPath);

  try {
    // 3. Читаем и парсим YAML
    const fileContents = readFileSync(configPath, "utf8");
    return yamlLoad(fileContents) as AppConfig;
  } catch (e) {
    console.warn(`⚠️ Не удалось прочитать конфиг по пути: ${configPath}`, e);
    throw e;
  }
}
