package msghandler

import (
	"encoding/json"
	"fmt"
	opts "github.com/Juniper/ipfix-translator/options"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type DataManager struct {
	netClient *http.Client
}

type DMMessage struct {
	CollectionName  string           `json:"collection_name"`
	Data            AugmentedMessage `json:"data"`
	TailwindManager interface{}      `json:"tailwind_manager"`
}

func (dm *DataManager) setup(configFile string) error {
	dm.netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	return nil
}

func (dm *DataManager) handleMessages(mhChan chan []byte) {
	var (
		msg []byte
	)
	for {
		select {
		case msg = <-mhChan:
			if opts.Verbose {
				log.Println("Received Message on DM Handler ", string(msg))
			}
			dm.pushDataToDataManager(msg)
		}
	}
}

func (dm *DataManager) splitDataSets(msg *Message) []DMMessage {
	res := make([]DMMessage, len(msg.DataSets))
	for i, dataSet := range msg.DataSets {
		/*
		   Do not compute Timestamp from msg.Header.ExportTime, as if the time is
		   not properly set in the sender, then time will not sync between Txer and
		   Rxer
		   tsInMilliSecs := msg.Header.ExportTime * 1000
		*/
		tsInMilliSecs := time.Now().UnixNano() / 1000000
		res[i] = DMMessage{CollectionName: "ipfix_flow_collection",
			Data: AugmentedMessage{Header: msg.Header, DataSets: dataSet, AgentID: msg.AgentID,
				RoomKey: msg.AgentID, Timestamp: tsInMilliSecs}}
		var emptyTM struct{}
		res[i].TailwindManager = &emptyTM
	}
	return res
}

func (dm *DataManager) serializeDataToDataManager(msg []byte) ([]DMMessage, error) {
	var p Message
	err := json.Unmarshal(msg, &p)
	if err != nil {
		log.Println("data Serialize json.Unmarshal() error ", err)
		return nil, err
	}
	dmMsgs := dm.splitDataSets(&p)
	return dmMsgs, nil
}

func (dm *DataManager) pushDataToDataManager(msg []byte) error {
	var (
		reqUrl      string
		contentType string
	)
	contentType = "application/json"
	dmMsgs, err := dm.serializeDataToDataManager(msg)
	if err != nil {
		log.Println("data serialize->DM error ", err)
		return err
	}
	msgCnt := len(dmMsgs)
	reqUrl = fmt.Sprintf("http://%s:%s/version/2.0/post_event", opts.DataMgrIpAddress, opts.DataMgrPort)
	for idx := 0; idx < msgCnt; idx++ {
		dmMsg, err := json.Marshal(dmMsgs[idx])
		if err != nil {
			log.Println("data json.Marshal() error ", err)
		}
		if opts.Verbose {
			log.Println("Sending POST data to DM ", reqUrl, string(dmMsg))
		}
		response, err := dm.netClient.Post(reqUrl, contentType, strings.NewReader(string(dmMsg)))
		if err != nil {
			log.Println("DataManager POST error ", err)
		} else {
			defer response.Body.Close()
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatalln("Parse error response body ", err)
			}
			if opts.Verbose {
				log.Println("Getting response from DM ", string(body))
			}
		}
	}
	return nil
}