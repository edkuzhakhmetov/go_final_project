package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()

	log.SetFormatter(&logrus.JSONFormatter{})

	log.SetLevel(logrus.TraceLevel)

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Warn("Не удалось записать логи в файл, используется стандартный вывод")
	}
}

func GetLogger() logrus.FieldLogger {
	return log
}
