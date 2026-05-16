package proxy

import (
	"fmt"

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
		router.All(service.Cfg.Path+"/*", func(c fiber.Ctx) error {
			// log.Info(c.Params("*"))
			// log.Info(c.Request())
			target := fmt.Sprintf("%s://%s:%d",
				service.Cfg.Protocol,
				service.Cfg.Host,
				service.Cfg.Port,
			)

			path := c.Params("*")

			// CN/SAN backend cert
			proxyClient.TLSConfig.ServerName = service.Name
			if err := proxy.Do(c, target+service.Cfg.Path+"/"+path, proxyClient); err != nil {
				log.Printf("Proxy error: %v", err)
				return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
					"error": "Service unavailable",
				})
			}
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
