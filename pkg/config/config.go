package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type Configuration struct {
	TimeOut int `json:"time_out"`
	Log     struct {
		Format      string `json:"format"`
		ServiceName string `json:"service_name"`
		Level       int    `json:"level"`
	} `json:"log"`
}

var (
	conf Configuration
)

// NewConfig создание нового экземпляра конфига
func NewConfig(v *viper.Viper) Configuration {
	// если директории нет то создаем
	if _, err := os.Stat("./var/log"); os.IsNotExist(err) {
		err = os.MkdirAll("./var/log", os.ModePerm)
		if err != nil {
			log.Fatal("Error on create log directory:", err.Error())
		}
	}

	v.BindEnv("log.level")
	v.BindEnv("log.name")
	v.BindEnv("timeout")

	c := Configuration{}
	c.Log.ServiceName = v.GetString("log.name")
	c.Log.Format = "[%s] %s.%s message: %s context: %s extra: %s"
	c.Log.Level = v.GetInt("log.level")
	c.TimeOut = v.GetInt("timeout")
	return c
}

// NewDefaultConfig обновить настройки экземпляра конфига
func NewDefaultConfig(v *viper.Viper) {
	// если директории нет то создаем
	if _, err := os.Stat("./var/log"); os.IsNotExist(err) {
		err = os.MkdirAll("./var/log", os.ModePerm)
		if err != nil {
			log.Fatal("Error on create log directory:", err.Error())
		}
	}

	v.BindEnv("log.level")
	v.BindEnv("log.name")
	v.BindEnv("timeout")

	conf := Configuration{}
	conf.Log.ServiceName = v.GetString("log.name")
	conf.Log.Format = "[%s] %s.%s message: %s context: %s extra: %s"
	conf.Log.Level = v.GetInt("log.level")
	conf.TimeOut = v.GetInt("timeout")
}

func Config() *Configuration {
	return &conf
}
