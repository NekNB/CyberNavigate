package middlewares

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func LoggerMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		// 1. Сохраняем тело входящего запроса
		reqBody := c.Body()

		// Собираем данные ЗАПРОСА
		reqLog := map[string]any{
			"phase":      "REQUEST",
			"method":     c.Method(),
			"path":       c.Path(),
			"query":      c.Queries(),
			"host":       c.Hostname(),             // Host, на который пришел запрос
			"remote_ip":  c.IP(),                   // IP клиента
			"remote_hdr": c.Get("X-Forwarded-For"), // Реальный IP, если за прокси
			"headers":    formatHeaders(c.GetReqHeaders()),
			"body":       string(reqBody),
		}

		// Сериализуем в одну строку JSON (без отступов)
		if jsonBytes, err := json.Marshal(reqLog); err == nil {
			log.Println(string(jsonBytes))
		}

		// ВАЖНО: Возвращаем тело обратно в контекст,
		// чтобы следующие обработчики (handlers) могли его прочитать
		c.Request().SetBody(reqBody)

		// 2. Передаем запрос дальше по цепочке
		err := c.Next()

		// 3. Собираем данные ОТВЕТА (после того, как отработал handler)
		resLog := map[string]interface{}{
			"phase":   "RESPONSE",
			"status":  c.Response().StatusCode(),
			"headers": formatHeaders(c.GetRespHeaders()),
			"body":    string(c.Response().Body()),
			"error":   nil, // По умолчанию ошибок нет
		}

		// Если handler вернул ошибку, добавляем её в лог
		if err != nil {
			resLog["error"] = err.Error()
		}

		// Сериализуем и выводим ответ
		if jsonBytes, err := json.Marshal(resLog); err == nil {
			log.Println(string(jsonBytes))
		}

		return err // Обязательно возвращаем ошибку, если она была
	}
}

// Вспомогательная функция: Fasthttp возвращает заголовки в виде map[string][]string
// Эта функция склеивает массивы значений в одну строку через запятую для компактности
func formatHeaders(h map[string][]string) map[string]string {
	result := make(map[string]string)
	for k, v := range h {
		result[k] = strings.Join(v, ", ")
	}
	return result
}

// Использование:
// app.Use(LoggerMiddleware())
