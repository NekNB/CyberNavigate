package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env           string       `yaml:"env" env-default:"local"`
	Server        serverConfig `yaml:"server" env-required:"true"`
	Certs         certsConfig  `yaml:"certs" env-required:"true"`
	PublicKeyPath string       `yaml:"public_key_path"`
}

type certsConfig struct {
	CaCertPath string `yaml:"ca_cert" env-required:"true"`

	ClientCertPath string `yaml:"client_cert" env-required:"true"`
	ClientKeyPath  string `yaml:"client_key" env-required:"true"`

	PublicCertPath string `yaml:"public_cert" env-required:"true"`
	PublicKeyPath  string `yaml:"public_key" env-required:"true"`
}
type serverConfig struct {
	Host          string                  `yaml:"host"`
	HTTPPort      int                     `yaml:"http_port"`
	HTTPSPort     int                     `yaml:"https_port"`
	Services      []map[string]serviceCfg `yaml:"services"`
	GatewayOrigin string                  `yaml:"gateway_origin"`
	AllowHeaders  []string                `yaml:"allow_headers"`
	AllowOrigins  []string                `yaml:"allow_origins"`
}

type serviceCfg struct {
	Path         string        `yaml:"path"`
	Protocol     string        `yaml:"protocol" env-default:"http"`
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	ConnDuration time.Duration `yaml:"conn_duration"`
}

// "Must" means the function will panic rather than return an error
// Only used during application startup
func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}
	return MustLoadByPath(path)
}

func MustLoadByPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

type Service struct {
	Name string
	Cfg  serviceCfg
}

func Normalize(cfg *Config) []Service {
	var result []Service

	for _, item := range cfg.Server.Services {
		for name, svc := range item {
			result = append(result, Service{
				Name: name,
				Cfg:  svc,
			})
		}
	}

	return result
}
