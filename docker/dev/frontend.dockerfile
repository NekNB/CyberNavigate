# 1. Зависимости
FROM node:24-alpine AS deps
RUN apk add --no-cache libc6-compat
WORKDIR /app

# Копируем манифесты зависимостей (подходит для npm, yarn, pnpm)
COPY /frontend/package.json /frontend/package-lock.json* ./
RUN npm ci

# 2. Сборка приложения
FROM node:24-alpine AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY /frontend .

# Отключаем телеметрию Next.js во время сборки
ENV NEXT_TELEMETRY_DISABLED=1
ENV CONFIG_PATH="./dev.yaml"
COPY /configs/frontend/dev.yaml .

RUN  npm run build

# 3. Финальный образ для запуска
FROM node:24-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

# Создаем не-root пользователя для безопасности
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public

# Автоматически копируем минимальные зависимости standalone-сборки
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs
ENV CONFIG_PATH="./dev.yaml"
COPY /configs/frontend/dev.yaml .
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

CMD ["node", "server.js"]