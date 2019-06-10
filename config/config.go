package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Debug bool
	Database
	Protocol struct {
		Grpc Host
		Http Host
	}
	Logging struct {
		LogLevel      int
		LogTimeFormat string
	}
}

type Database struct {
	MySqlDB struct {
		Host     string
		Port     string
		Username string
		Password string
		DBName   string
		SSLMode  bool
	}
}

type Host struct {
	Address string
	Port    string
}

var config Config

func InitViper(currentPath string, env string) {
	fmt.Println("env", env)
	switch env {
	case "develop":
		viper.SetConfigName("dev-config")
		break
	case "staging":
		viper.SetConfigName("staging-config")
		break
	case "prod":
		viper.SetConfigName("config")
		break
	default:
		viper.SetConfigName("config")
		break
	}
	viper.SetConfigType("yaml")
	viper.AddConfigPath(currentPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	viper.WatchConfig() // Watch for changes to the configuration file and recompile
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Println(err)
	}
}

func GetViper() *Config {
	return &config
}
