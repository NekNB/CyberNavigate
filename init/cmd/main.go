package main

import (
	"github.com/NekNB/CyberNavigate/init/internal/http"
	parser "github.com/NekNB/CyberNavigate/init/internal/init"
	"github.com/NekNB/CyberNavigate/init/internal/lib/logger"
)

func main() {

	log, err := logger.Init("local")
	if err != nil {
		panic(err)
	}

	apiClient := http.New(log, "http://127.0.0.1:9080/api/v1")

	// TODO: поменять ввод username/password
	apiClient.Login("Арей", "123")

	log.Info("Инициализация проекта началась...")
	if err := parser.ProcessArticles(apiClient, log); err != nil {
		log.Warning(err)
	}

	if err := parser.InitScenarios(apiClient, log); err != nil {
		log.Warning(err)
	}

}
