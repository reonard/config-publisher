package main

import (
	"config-publisher/db"
	. "config-publisher/lib"
	"config-publisher/processor"
	"encoding/json"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"time"
)

func main() {

	InitCfg()
	fmt.Println(AppCfg.AlarmTopics)

	//s := db.InitMongoDB(AppCfg.MongodbURL)
	//defer s.Close()

	m := db.InitMySQLDB(AppCfg.MySqlURL)
	defer m.Close()

	cliOpt := mqtt.NewClientOptions()

	cliOpt.AddBroker("xx").SetClientID("xx").SetUsername("xx")

	mqttCli := mqtt.NewClient(cliOpt)

	if token := mqttCli.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for{

		time.Sleep(5 * time.Second)

		rows, err:=processor.GetLatestPushData(1)

		if err!=nil{
			fmt.Println("Error Get Push Data ", err)
			rows.Close()
			continue
		}

		for rows.Next(){

			message:=processor.DeviceActionMsg{}

			var ProbeConfigText []byte

			err := rows.Scan(&message.DeviceId, &message.IssueStatusId, &message.ActionType, &ProbeConfigText)
			if err != nil{
				log.Println("Error Scan Push Data ", err)
				continue
			}
			log.Printf("Processing MsgId: %d ", message.IssueStatusId)

			err = json.Unmarshal(ProbeConfigText, &message.ProbeConfig)
			if err != nil{
				log.Printf("Error Unmarshal ProbeConfig: %v", err)
				continue
			}

			publishMsg, err := json.Marshal(message)
			if err != nil{
				log.Printf("Error Marshal Push Msg: %v", err)
				continue
			}

			topic:= fmt.Sprintf("client/%d/action", message.DeviceId)
			if token := mqttCli.Publish(topic, 1, true, publishMsg); token.Wait() && token.Error() != nil {
				log.Println("Could not publish message ", token.Error())
				continue
			}

			processor.UpdateIssueStatus(message.IssueStatusId, processor.STATUS_PUSHED)
		}

		rows.Close()
	}

}
