package env

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Env struct {
	App struct {
		Name string `yaml:"name"`
	}

	Log struct {
		Level       string `yaml:"level"`
		Environment string `yaml:"environment"`
	} `yaml:"log"`

	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
		User     string `yaml:"user"`
		Pass     string `yaml:"pass"`
		SSLMode  string `yaml:"sslmode"`
		Timezone string `yaml:"timezone"`
	} `yaml:"database"`
}

func Load() {
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		zap.L().Fatal("failed to read config", zap.Error(err))
	}
}
