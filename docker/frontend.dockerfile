# Этап 1: Сборка приложения
FROM node:20-alpine AS builder
WORKDIR /app
# Копируем только package файлы для кеширования зависимостей
COPY /frontend/package*.json ./
RUN npm ci
# Копируем исходный код и собираем
COPY /frontend .
COPY /configs/frontend/dev.yaml .
ENV CONFIG_PATH=./dev.yaml
RUN npm run build

# Этап 2: Отдача через Nginx
FROM nginx:alpine
# Копируем собранные файлы из первого этапа в папку Nginx
COPY --from=builder /app/dist /usr/share/nginx/html
# Заменяем стандартный конфиг Nginx на наш (нужен для роутинга React/Vue)
COPY /frontend/nginx/nginx.conf /etc/nginx/conf.d/default.conf
CMD ["nginx", "-g", "daemon off;"]