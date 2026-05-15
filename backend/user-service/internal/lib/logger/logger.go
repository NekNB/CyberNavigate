package logger

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

func Init(runMode string) (*logrus.Logger, error) {

	projectRoot, err := projectRoot()

	if err != nil {
		return nil, err
	}

	log := logrus.New()

	// Настройка logger под разные режимы запуска local/dev/prod
	switch runMode {
	case "local":
		log.SetLevel(logrus.DebugLevel)
		log.SetReportCaller(true)
		log.SetFormatter(
			&logrus.TextFormatter{
				ForceColors:     true,
				FullTimestamp:   true,
				TimestampFormat: "2006-01-02 15:04:05",
				CallerPrettyfier: func(f *runtime.Frame) (string, string) {
					// Обрезаем путь до корня проекта
					file := f.File
					if idx := strings.Index(file, projectRoot); idx != -1 {
						file = file[idx+len(projectRoot):]
					}
					return "", fmt.Sprintf(" => %s:%d =>", file, f.Line)
				},
			},
		)
	case "dev":
		log.SetLevel(logrus.DebugLevel)
		log.SetFormatter(&logrus.JSONFormatter{})
	case "prod":
		log.SetLevel(logrus.InfoLevel)
		log.SetFormatter(&logrus.JSONFormatter{})
	default:
		return nil, fmt.Errorf("Not known logLevel: %s", runMode)
	}

	log.Infoln("the logger is configured")

	return log, nil

}

func projectRoot() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("Cannot get caller info")
	}

	projectRoot := filepath.Join(file, "../../..")
	projectRoot = strings.Join(strings.Split(projectRoot, "\\"), "/")

	return projectRoot, nil
}
