package proxy

import (
	"fmt"
	"slices"
	"strings"

	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/NekNB/CyberNavigate/backend/gateway-server/internal/config"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/proxy"
	"github.com/sirupsen/logrus"

	"github.com/valyala/fasthttp"
)

func Register(cfg *config.Config, log *logrus.Logger, router fiber.Router) {

	proxyClient, err := proxyConfig(cfg, log)
	if err != nil {
		panic(err)
	}

	services := config.Normalize(cfg)
	for _, service := range services {

		serviceTLSConfig := proxyClient.TLSConfig.Clone()
		serviceTLSConfig.ServerName = service.Name

		serviceClient := &fasthttp.Client{
			TLSConfig:                     serviceTLSConfig,
			NoDefaultUserAgentHeader:      true,
			DisableHeaderNamesNormalizing: true,
			MaxIdleConnDuration:           service.Cfg.ConnDuration,
		}

		router.All(service.Cfg.Path+"/*", func(c fiber.Ctx) error {

			if c.Method() == fiber.MethodOptions {
				setGatewayCorsHeaders(c, cfg)
			}

			target := fmt.Sprintf("%s://%s:%d",
				service.Cfg.Protocol,
				service.Cfg.Host,
				service.Cfg.Port,
			)

			path := c.Params("*")

			if err := proxy.Do(c, target+service.Cfg.Path+"/"+path, serviceClient); err != nil {
				log.Printf("Proxy error: %v", err)
				return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
					"error": "Service unavailable",
				})
			}
			// 3. Очищаем CORS заголовки, которые мог прислать бэкенд
			// Используем указатель, чтобы избежать ошибки noCopy
			h := &c.Response().Header
			h.Del(fiber.HeaderAccessControlAllowOrigin)
			h.Del(fiber.HeaderAccessControlAllowMethods)
			h.Del(fiber.HeaderAccessControlAllowHeaders)
			h.Del(fiber.HeaderAccessControlAllowCredentials)
			h.Del(fiber.HeaderAccessControlExposeHeaders)
			h.Del(fiber.HeaderAccessControlMaxAge)

			setGatewayCorsHeaders(c, cfg)
			return nil
		})

	}
}

func proxyConfig(cfg *config.Config, log *logrus.Logger) (*fasthttp.Client, error) {

	gatewayCert, err := tls.LoadX509KeyPair(
		cfg.Certs.ClientCertPath,
		cfg.Certs.ClientKeyPath,
	)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	caCert, err := os.ReadFile(cfg.Certs.CaCertPath)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	caPool := x509.NewCertPool()

	if !caPool.AppendCertsFromPEM(caCert) {
		log.Error(err)
		return nil, fmt.Errorf("failed append ca")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{
			gatewayCert,
		},

		RootCAs: caPool,

		MinVersion: tls.VersionTLS13,
	}

	return &fasthttp.Client{
		TLSConfig: tlsConfig,

		NoDefaultUserAgentHeader: true,
		DisablePathNormalizing:   true,
		MaxIdleConnDuration:      60,
	}, nil
}

func setGatewayCorsHeaders(c fiber.Ctx, cfg *config.Config) {
	// 1. Получаем Origin из входящего запроса
	origin := string(c.Request().Header.Peek(fiber.HeaderOrigin))

	// Если заголовка Origin нет (например, это запрос от сервера к серверу или curl),
	// CORS-заголовки не нужны вообще.
	if origin == "" {
		return
	}

	// Если Origin не разрешен, просто выходим.
	// Браузер увидит отсутствие Access-Control-Allow-Origin и заблокирует ответ.
	if !slices.Contains(cfg.Server.AllowOrigins, origin) {
		return
	}

	// 3. Устанавливаем заголовки
	h := &c.Response().Header

	// Возвращаем именно тот Origin, с которого пришел запрос (нельзя использовать "*" при куках!)
	h.Set(fiber.HeaderAccessControlAllowOrigin, origin)
	h.Set(fiber.HeaderAccessControlAllowCredentials, "true")
	h.Set(fiber.HeaderAccessControlAllowMethods, strings.Join([]string{
		"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS",
	}, ", "))
	// Сообщаем кэшу/CDN, что ответ зависит от заголовка Origin (важно для CORS)
	h.Add(fiber.HeaderVary, fiber.HeaderOrigin)

	// 4. Указываем, какие заголовки разрешаем читать фронтенду
	if len(cfg.Server.AllowHeaders) > 0 {
		h.Set(fiber.HeaderAccessControlExposeHeaders, strings.Join(cfg.Server.AllowHeaders, ", "))
	}
}
