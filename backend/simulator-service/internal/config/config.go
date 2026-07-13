package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env     string         `yaml:"env" env-default:"local"`
	HTTP    HTTPConfig     `yaml:"http"`
	Storage StoragesConfig `yaml:"storages"`
	Certs   certsConfig    `yaml:"certs" env-required:"true"`
}

type certsConfig struct {
	CaCertPath string `yaml:"ca_cert" env-required:"true"`

	ServerCertPath string `yaml:"server_cert" env-required:"true"`
	ServerKeyPath  string `yaml:"server_key" env-required:"true"`
}

type HTTPConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type StoragesConfig struct {
	Postgres DatabaseConfig `yaml:"postgres"`
}

type DatabaseConfig struct {
	Database string `yaml:"database"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string
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

	fetchDatabasePasswords(&cfg)

	return &cfg
}

// Получение пути к конфигу из аргументов запуска или Переменной CONFIG_PATH
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

// Получение паролей к базам данных из Env
func fetchDatabasePasswords(cfg *Config) {
	cfg.Storage.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	fmt.Println(os.Getenv("POSTGRES_PASSWORD"))
}
