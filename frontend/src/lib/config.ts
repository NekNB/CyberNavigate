// src/lib/config.ts
import { readFileSync } from "fs";
import path from "path";
import { load as yamlLoad } from "js-yaml";
import "server-only";

export interface AppConfig {
  internalUrl: string;
  externalUrl: string;
}

// 1. Читаем путь из переменной
const envPath = process.env.CONFIG_PATH || "config.yaml";

// 2. Игнорируем трейсинг Next.js для динамического пути
const configPath = path.isAbsolute(envPath)
  ? envPath
  : path.resolve(/*turbopackIgnore: true*/ process.cwd(), envPath);

// 3. Читаем и парсим файл единовременно при сборке
let staticConfig: AppConfig;

try {
  const fileContents = readFileSync(configPath, "utf8");
  staticConfig = yamlLoad(fileContents) as AppConfig;
} catch (e) {
  console.error(`❌ Ошибка чтения конфига при сборке по пути: ${configPath}`, e);
  throw e;
}

export function loadConfig(): AppConfig {
  return staticConfig;
}