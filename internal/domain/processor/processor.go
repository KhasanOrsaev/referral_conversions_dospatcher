package processor

import (
	clickhouse_client_git "git.fin-dev.ru/dmp/clickhouse_client.git"
)

type Processor struct {
	clickhouse_client_git.ClickHouseClient
}

func NewProcessor() *Processor {
	return &Processor{}
}
