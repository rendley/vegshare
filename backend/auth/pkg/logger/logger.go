package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

func New() *logrus.Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)                  // Вывод в консоль
	log.SetFormatter(&logrus.JSONFormatter{}) // Устанавливаем JSON-формат
	log.SetLevel(logrus.InfoLevel)            // Уровень логирования
	return log
}
