package main

import (
	"fmt"
	logger2 "git.fin-dev.ru/dmp/logger.git"
	"git.fin-dev.ru/dmp/referral_coversions_dispatcher.git/internal/repository"
	"git.fin-dev.ru/dmp/referral_coversions_dispatcher.git/pkg/config"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"strings"
	"sync"
	"time"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logger2.Fatal("err", "", nil, err)
	}
}

func main() {
	v := viper.New()
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)
	c := config.NewConfig(v)
	for {
		logger, err := logger2.NewLogger(map[string]interface{}{
			"module":    c.Log.ServiceName,
			"submodule": "dispatcher",
			"level":     c.Log.Level,
			"format":    c.Log.Format,
		})
		if err != nil {
			logger.Fatal("create new logger", "", nil, err)
		}
		timeStart := time.Now()
		wgDispatcher := sync.WaitGroup{}

		// ----------------- инициализация источника ----------------------------
		sourceRepository := repository.RabbitSourceRepository{V: v}
		sourceClient, err := sourceRepository.NewSourceRepository(v)
		if err != nil {
			logger.Fatal("init source client", "", nil, err)
		}
		err = sourceClient.OpenConnection()
		if err != nil {
			logger.Fatal(" open source connection", "", nil, err)
		}
		// ----------------- завершение инициализации источника ----------------------------

		// ----------------- инициализация пункта назначения ----------------------------
		destinationRepository := repository.ClickhouseDestinationRepository{V: v}
		destinationClient, err := destinationRepository.NewDestinationRepository(v)
		if err != nil {
			logger.Fatal("init destination client", "", nil, err)
		}
		err = destinationClient.OpenConnection()
		if err != nil {
			logger.Fatal(" open destination connection", "", nil, err)
		}
		// ----------------- завершение инициализации пункта назначения ----------------------------

		// ----------------- инициализация обработчика ----------------------------
		processorClient, err := repository.NewProcessRepository(v)
		if err != nil {
			logger.Fatal("init processor client", "", nil, err)
		}
		err = processorClient.OpenConnection()
		if err != nil {
			logger.Fatal(" open processor connection", "", nil, err)
		}
		// ----------------- завершение инициализации обработчика ----------------------------

		// ----------------- инициализация аварийки ----------------------------
		v.SetEnvPrefix("crash")
		crashRepository := repository.RabbitCrashRepository{V: v}
		crashClient, err := crashRepository.NewCrashRepository(v)
		if err != nil {
			logger.Fatal("init crash client", "", nil, err)
		}
		err = crashClient.OpenConnection()
		if err != nil {
			logger.Fatal(" open crash connection", "", nil, err)
		}
		// ----------------- завершение инициализации аварийки ----------------------------

		wgDispatcher.Add(5)
		processorChannel := make(chan map[interface{}][]byte, 10)
		inputChannel := make(chan map[interface{}][]byte, 10)
		confirmChannel := make(chan interface{}, 10)
		crashChannel := make(chan map[uuid.UUID][]byte, 1)

		go func() {
			defer wgDispatcher.Done()
			sourceClient.ReadData(inputChannel, logger)
			close(inputChannel)
		}()

		go func() {
			defer wgDispatcher.Done()
			sourceClient.Confirm(confirmChannel, logger)
			err = sourceClient.CloseConnection()
			if err != nil {
				logger.Error(" close source connection", "", nil, err)
			}
		}()
		// process
		go func() {
			defer wgDispatcher.Done()
			processorClient.ProcessData(inputChannel, processorChannel, logger)
			close(processorChannel)
			err = processorClient.CloseConnection()
			if err != nil {
				logger.Error(" close processor connection", "", nil, err)
			}
		}()

		go func() {
			defer wgDispatcher.Done()
			destinationClient.WriteData(processorChannel, confirmChannel, crashChannel, logger)
			// закрываем каналы при завершении
			close(confirmChannel)
			close(crashChannel)
			err = destinationClient.CloseConnection()
			if err != nil {
				logger.Error(" close destination connection", "", nil, err)
			}
		}()

		go func() {
			defer wgDispatcher.Done()
			crashClient.SaveData(crashChannel, logger)
			err = crashClient.CloseConnection()
			if err != nil {
				logger.Error(" close crash connection", "", nil, err)
			}
		}()
		wgDispatcher.Wait()
		fmt.Println("Время заняло - ", time.Since(timeStart).Seconds())
		time.Sleep(time.Duration(c.TimeOut) * time.Minute)
	}

}
