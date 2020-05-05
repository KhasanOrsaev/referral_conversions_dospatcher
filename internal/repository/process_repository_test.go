package repository

import (
	"fmt"
	dispatcher_interface "git.fin-dev.ru/dmp/dispatcher-interface.git"
	"git.fin-dev.ru/dmp/logger.git"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

var (
	processorClient dispatcher_interface.Processor
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	v := viper.New()
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)
	processorClient, err = NewProcessRepository(v)
	if err != nil {
		logger.Fatal("init processor client", "", nil, err)
	}
	err = processorClient.OpenConnection()
	if err != nil {
		logger.Fatal(" open processor connection", "", nil, err)
	}
}

func TestProcessRepository_ProcessData(t *testing.T) {
	processorChannel := make(chan map[interface{}][]byte, 10)
	inputChannel := make(chan map[interface{}][]byte, 10)
	logg, err := logger.NewLogger(map[string]interface{}{
		"module":    "test",
		"submodule": "process",
		"level":     4,
		"format":    "[%s] %s.%s message: %s context: %s extra: %s",
	})
	if err != nil {
		t.Fatal("create new logger", "", nil, err)
	}

	data, err := ioutil.ReadFile("test_data")
	if err != nil {
		t.Fatal("create new logger", "", nil, err)
	}
	inputChannel <- map[interface{}][]byte{1: data}
	close(inputChannel)
	go func() {
		processorClient.ProcessData(inputChannel, processorChannel, logg)
		close(processorChannel)
		err := processorClient.CloseConnection()
		if err != nil {
			t.Error(err)
		}
	}()
	for v := range processorChannel {
		for _, vv := range v {
			fmt.Println(string(vv))
		}
	}
}
