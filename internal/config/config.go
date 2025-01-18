package config

import (
	"encoding/json"
	"os"
	"quiz-mod/internal/model"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	FirstLevel  []model.SimpleQuestion `json:"first_level"`
	SecondLevel []*model.HardQuestion  `json:"second_level"`
	ThirdLevel  []model.SimpleQuestion `json:"third_level"`
}

func LoadConfig(filename, path string) (*Config, error) {
	viper.SetConfigFile(filename)
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	cfg := Config{}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	logrus.Infof("Config: first level: %d question(s)", len(cfg.FirstLevel))
	logrus.Infof("Config: second level: %d question(s)", len(cfg.SecondLevel))
	logrus.Infof("Config: third level: %d question(s)", len(cfg.ThirdLevel))

	return &cfg, nil
}
