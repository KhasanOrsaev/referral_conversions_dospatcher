package config

import (
	"encoding/json"
	"github.com/go-errors/errors"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Configuration struct {
	TimeOut  int `json:"time_out"`
	Services struct {
		Source      *json.RawMessage `json:"source"`
		Destination *json.RawMessage `json:"destination"`
		Crash       *json.RawMessage `json:"crash"`
	} `json:"services"`
	Log struct {
		Format      string `json:"format"`
		ServiceName string `json:"service_name"`
		Level       int `json:"level"`
	} `json:"log"`
}

var (
	conf   Configuration
)

// инициализация, получение данных из конфига, создание конфига, создание логера
func Init() error {
	// если директории нет то создаем
	if _, err := os.Stat("./var/log"); os.IsNotExist(err) {
		err = os.MkdirAll("./var/log", os.ModePerm)
		if err != nil {
			return errors.Wrap(err, -1)
		}
	}
	v := viper.New()
	// обязательно заменять . на _ при поиске переменной
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)

	v.BindEnv("log.level")
	v.BindEnv("log.name")
	v.BindEnv("timeout")

	conf.Log.ServiceName = v.GetString("log.name")+"_dispatcher"
	conf.Log.Format = "[%s] %s.%s message: %s context: %s extra: %s"
	conf.Log.Level = v.GetInt("log.level")
	conf.TimeOut = v.GetInt("timeout")
	return nil
}

func GetConfig() *Configuration {
	return &conf
}
