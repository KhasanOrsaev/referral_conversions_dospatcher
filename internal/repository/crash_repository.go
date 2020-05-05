package repository

import (
	dispatcher_interface "git.fin-dev.ru/dmp/dispatcher-interface.git"
	rabbitmq_client_git "git.fin-dev.ru/dmp/rabbitmq_client.git"
	"github.com/spf13/viper"
)

type RabbitCrashRepository struct {
	V *viper.Viper
}

type CrashRepositoryInterface interface {
	NewSourceRepository(v *viper.Viper) dispatcher_interface.Crash
}

func (RabbitCrashRepository) NewCrashRepository(v *viper.Viper) (dispatcher_interface.Crash, error) {
	rabbit := rabbitmq_client_git.NewClient()
	err := rabbit.SetConfig(v)
	return rabbit, err
}
