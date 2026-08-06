package parser

import (
	"io/fs"

	"github.com/NekNB/CyberNavigate/init/internal/assets"
	"github.com/NekNB/CyberNavigate/init/internal/http"
	"github.com/NekNB/CyberNavigate/init/internal/init/scenarios"
	"github.com/sirupsen/logrus"
)

func InitScenarios(apiClient *http.APIClient, log *logrus.Logger) error {

	// Читаем содержимое директории "simulator" из встроенной ФС
	entries, err := fs.ReadDir(assets.SimulatorFS, "simulator")
	if err != nil {
		log.Warnf("Ошибка чтения директории: %v", err)
	}

	for _, entry := range entries {
		// Пропускаем поддиректории (если они вдруг есть), берем только файлы
		if entry.IsDir() {
			continue
		}

		// Формируем полный путь к файлу внутри embed.FS
		filePath := "simulator/" + entry.Name()

		// Читаем содержимое файла в слайс байтов
		data, err := fs.ReadFile(assets.SimulatorFS, filePath)
		if err != nil {
			log.Warnf("Ошибка чтения файла %s: %v", filePath, err)
			continue // Если один файл битый, продолжаем со следующим
		}

		// Вызываем вашу функцию
		log.Infof("Генерируем сценарий из: %s\n", entry.Name())
		if err := scenarios.GenerateScenario(data, apiClient, log); err != nil {
			log.Warn(err)
			return err
		}
	}
	return nil
}
