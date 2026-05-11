package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env    string       `yaml:"env" env-default:"local"`
	Server ServerConfig `yaml:"server" env-required:"true"`
}

type ServerConfig struct {
	Port     int                     `yaml:"port"`
	Services []map[string]ServiceCfg `yaml:"services"`
}

type ServiceCfg struct {
	Path     string `yaml:"path"`
	Protocol string `yaml:"protocol" env-default:"http"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
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
	Cfg  ServiceCfg
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
