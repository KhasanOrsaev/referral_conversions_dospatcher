package repository

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	dispatcher_interface "git.fin-dev.ru/dmp/dispatcher-interface.git"
	logger2 "git.fin-dev.ru/dmp/logger.git"
	"git.fin-dev.ru/dmp/referral_coversions_dispatcher.git/internal/domain/processor"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type ProcessRepository struct {
	V *viper.Viper
	processor.Processor
}

const (
	querySelect = "select distinct user_id, application_id from cs.events where source='userid' and dt_create> '2020-04-30 11:00:00' and application_id_native in (%s) format JSONCompact"
)

func NewProcessRepository(v *viper.Viper) (dispatcher_interface.Processor, error) {
	p := ProcessRepository{}
	err := p.SetConfig(v)

	return &p, err
}

type Message map[string]interface{}

type Row struct {
	EventID             *uuid.UUID             `json:"event_id,omitempty"`
	UserID              *int                   `json:"user_id,omitempty"`
	ApplicationID       *uuid.UUID             `json:"application_id,omitempty"`
	ApplicationIDNative string                 `json:"application_id_native,omitempty"`
	Type                string                 `json:"type,omitempty"`
	Source              string                 `json:"source,omitempty"`
	OfferID             string                 `json:"offer_id,omitempty"`
	OfferStatus         string                 `json:"offer_status,omitempty"`
	DtEvent             string                 `json:"dt_event,omitempty"`
	DataJSON            map[string]interface{} `json:"data_json,omitempty"`
}

// ProcessData обработка данных, обогащение и отправка в канал отправки
func (p *ProcessRepository) ProcessData(outChannel <-chan map[interface{}][]byte,
	processedChannel chan<- map[interface{}][]byte, logger interface{}) {

	for row := range outChannel {
		for i, v := range row {
			var messages []Message
			err := json.Unmarshal(v, &messages)
			if err != nil {
				logger.(*logger2.Logger).Debug("unmarshal message", "",
					&map[string]interface{}{
						"submodule":   "processor",
						"method":      "ProcessData",
						"data":        string(v),
						"stack_trace": errors.Wrap(err, -1).ErrorStack(),
					}, err)
				continue
			}
			cnt := len(messages)
			processedData := make([]Row, cnt)
			userIDs := make([]string, cnt)
			userIDsMap := make(map[string]int, cnt)
			for m := range messages {
				if userID, ok := messages[m]["source"]; ok && userID.(string) != "" {
					userIDs[m] = userID.(string)
					userIDsMap[userID.(string)] = m
					processedData[m].ApplicationIDNative = userID.(string)
					delete(messages[m], "source")
				}
				// application id taken as hasoffers_id
				if hasOffer, ok := messages[m]["hasoffers_id"]; ok {
					hash := md5.New()
					hash.Write([]byte(strconv.Itoa(int(hasOffer.(float64)))))
					appID, _ := uuid.Parse(hex.EncodeToString(hash.Sum(nil)))
					processedData[m].ApplicationID = &appID
				}
				if offerID, ok := messages[m]["offer_id"]; ok {
					processedData[m].OfferID = strconv.Itoa(int(offerID.(float64)))
					delete(messages[m], "offer_id")
				}
				if offerStatus, ok := messages[m]["status"]; ok {
					processedData[m].OfferStatus = offerStatus.(string)
					delete(messages[m], "status")
				}
				if dtEvent, ok := messages[m]["lead_date"]; ok {
					processedData[m].DtEvent = dtEvent.(string)
					delete(messages[m], "lead_date")
				}
				id := uuid.New()
				processedData[m].EventID = &id
				processedData[m].Type = "lead"
				processedData[m].Source = "referral conversion"
				processedData[m].DataJSON = messages[m]
			}
			p.process(processedData, userIDs, userIDsMap, processedChannel, logger, i)

		}
	}
}

// process поиск пользователей и отправка данных в канал
func (p *ProcessRepository) process(processedData []Row, userIDs []string, userIDsMap map[string]int,
	processedChannel chan<- map[interface{}][]byte, logger interface{}, i interface{}) {
	findedUsers, err := p.findUsers(userIDs)
	if err != nil {
		logger.(*logger2.Logger).Fatal("find users", "",
			&map[string]interface{}{
				"submodule":   "processor",
				"method":      "process",
				"stack_trace": errors.Wrap(err, -1).ErrorStack(),
			}, err)
	}
	if findedUsers == nil {
		return
	}
	for _, user := range *findedUsers {
		index := userIDsMap[user[1]]
		id, _ := strconv.Atoi(user[0])
		processedData[index].UserID = &id
	}
	for index := range processedData {
		j, err := json.Marshal(processedData[index])
		if err != nil {
			logger.(*logger2.Logger).Fatal("marshal row", "",
				&map[string]interface{}{
					"submodule":   "processor",
					"method":      "process",
					"stack_trace": errors.Wrap(err, -1).ErrorStack(),
				}, err)
		}
		processedChannel <- map[interface{}][]byte{i: j}
	}
}

func (p *ProcessRepository) findUsers(userIDs []string) (*[][]string, error) {
	s := "'" + strings.Join(userIDs, "','") + "'"
	body := []byte(fmt.Sprintf(querySelect, s))
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	httpClient := http.Client{Transport: tr}
	request, err := http.NewRequest("POST", fmt.Sprintf("%s://%s:%s@%s:%s",
		p.Configuration.Protocol,
		p.Configuration.User,
		p.Configuration.Password,
		p.Configuration.Host,
		p.Configuration.Port,
	), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, err
	}
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	responseData := struct {
		Data [][]string `json:"data"`
	}{}
	err = json.Unmarshal(b, &responseData)
	if err != nil {
		return nil, err
	}
	return &responseData.Data, nil
}
