Referral Conversion Loader
====
Dispatcher loading data from rabbitmq to clickhouse.

Installation
--
```go
docker-compose -f deployment/docker-compose.test.yml up -d --build
```

_**Env file example:**_
```.env
CLICKHOUSE_PASSWORD=test
CLICKHOUSE_QUERIES=
CLICKHOUSE_USER=dispatcher
CLICKHOUSE_BULK=100
CLICKHOUSE_HOST=dispatcher-clickhouse

RABBITMQ_HOST=
RABBITMQ_PASSWORD=
RABBITMQ_PORT=
RABBITMQ_QUEUE_ARGUMENTS=
RABBITMQ_QUEUE_DURABLE=
RABBITMQ_QUEUE_NAME=
RABBITMQ_TLS=
RABBITMQ_USER=
RABBITMQ_VIRTUAL_HOST=

CRASH_RABBITMQ_QUEUE_ARGUMENTS=
CRASH_RABBITMQ_QUEUE_DURABLE=
CRASH_RABBITMQ_QUEUE_NAME=

LOG_NAME=referral_consumer
LOG_LEVEL=4

TIMEOUT=2

DOCKER_NETRC=
```

Project structure was taken from [here](https://github.com/golang-standards/project-layout).