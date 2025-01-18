package bot

import (
	"log"
	"os"
	"quiz-mod/internal/config"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func Start(envFile, confName, path string) {
	logrus.Info("starting")
	defer logrus.Info("stopped")

	logLvl := os.Getenv("LOG_LEVEL")

	switch logLvl {
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	case "fatal":
		logrus.SetLevel(logrus.PanicLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	_, err := config.LoadConfig(confName, path)
	if err != nil {
		logrus.Fatalf("unable to load config: %v", err)
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}
