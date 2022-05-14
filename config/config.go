package config

import (
	"fmt"
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
)

func InitConfig() error {

	env := os.Getenv("DOUYIN_ENV")
	if env == "" {
		env = "config"
	}

	viper.SetConfigName(env)
	viper.AddConfigPath("./config")
	viper.SetConfigType("yml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error while reading config file: %v", err))
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Configuration file changed")
	})

	if err := viper.Unmarshal(&G.Config); err != nil {
		panic(fmt.Errorf("fatal error while decode config file: %v", err))
	}

	err := initGorm()
	if err != nil {
		return err
	}
	err = initRedis()
	if err != nil {
		return err
	}
	return nil
}
