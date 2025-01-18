package bot

import "github.com/sirupsen/logrus"

func Start(confName, path string) {
	logrus.Info("starting")
	defer logrus.Info("stopped")
}
