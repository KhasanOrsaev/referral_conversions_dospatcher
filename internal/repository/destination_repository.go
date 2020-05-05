package repository

import (
	clickhouse_client_git "git.fin-dev.ru/dmp/clickhouse_client.git"
	dispatcher_interface "git.fin-dev.ru/dmp/dispatcher-interface.git"
	"github.com/spf13/viper"
)

type ClickhouseDestinationRepository struct {
	V *viper.Viper
}

type DestinationRepositoryInterface interface {
	NewDestinationRepository(v *viper.Viper) dispatcher_interface.Destination
}

func (ClickhouseDestinationRepository) NewDestinationRepository(v *viper.Viper) (dispatcher_interface.Destination, error) {
	click := clickhouse_client_git.NewClient()
	err := click.SetConfig(v)
	return click, err
}
