import babel from "@rolldown/plugin-babel";
import react, { reactCompilerPreset } from "@vitejs/plugin-react";
import fs from "fs";
import { load as yamlLoad } from "js-yaml"; // <-- ИСПРАВЛЕННЫЙ ИМПОРТ
import path from "path";
import { defineConfig } from "vite";
import type { AppConfig } from "./src/types/config";

// https://vite.dev/config/
export default defineConfig(() => {
  // 1. Читаем путь из переменной окружения (если нет — берем config.yaml по умолчанию)
  const envPath = process.env.CONFIG_PATH || "config.yaml";

  // 2. Преобразуем в абсолютный путь
  const configPath = path.resolve(__dirname, envPath);

  let appConfig: Partial<AppConfig> = {};

  try {
    // 3. Читаем и парсим YAML
    const fileContents = fs.readFileSync(configPath, "utf8");
    // <-- ИСПОЛЬЗУЕМ yamlLoad ВМЕСТО yaml.load
    appConfig = (yamlLoad(fileContents) as AppConfig) || {};
  } catch (e) {
    console.warn(`⚠️ Не удалось прочитать конфиг по пути: ${configPath}`, e);
  }

  return {
    define: {
      __APP_CONFIG__: JSON.stringify(appConfig),
    },
    server: {
      host: appConfig.server?.address || "localhost",
      port: appConfig.server?.port || 3777,
      //  ДОБАВЛЯЕМ ПРОКСИРОВАНИЕ 
      proxy: {
        '/api': {
          target: 'http://127.0.0.1:9080',
          changeOrigin: true,
          secure: false,
          
        },
      },
    },
    plugins: [react(), babel({ presets: [reactCompilerPreset()] })],
  };
});

