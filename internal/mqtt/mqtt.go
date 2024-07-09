package mqtt

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"hetic/tech-race/internal/models"
	"strconv"
	"strings"
	"time"
)

type MQTTClient struct {
	client MQTT.Client
	db     models.DatabaseInterface
}

func NewMQTTClient(db models.DatabaseInterface) *MQTTClient {
	opts := MQTT.NewClientOptions().AddBroker("tcp://192.168.52.82:1883")
	client := MQTT.NewClient(opts)
	return &MQTTClient{client: client, db: db}
}

func (m *MQTTClient) ConnectAndSubscribe() error {
	_, err := m.db.GetCurrentSessionID()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	if token := m.client.Subscribe("esp32/track", 0, m.MessageHandler); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	if token := m.client.Subscribe("esp32/sonar", 0, m.MessageHandler); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}
func (m *MQTTClient) MessageHandler(client MQTT.Client, msg MQTT.Message) {
	topic := msg.Topic()
	sessionID, err := m.db.GetCurrentSessionID()
	if err != nil {
		fmt.Println(err)
		return
	}
	switch topic {
	case "esp32/track":
		value, err := strconv.Atoi(string(msg.Payload()))
		if err != nil {
			fmt.Println(err)
			return
		}
		if value < 7 {
			data := models.LineTracking{LineTrackingValue: value, IDSession: sessionID}
			err = m.db.InsertTrackData(data)
			if err != nil {
				fmt.Println(err)
			}
		}
	case "esp32/sonar":
		distanceStr := strings.TrimSpace(string(msg.Payload()))
		distance, err := strconv.ParseFloat(distanceStr, 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		isCollision := false
		if distance < 10 {
			isCollision = true
		}
		if isCollision == true {
			timestamp := time.Now()
			data := models.Collision{Distance: distance, IsCollision: isCollision, Timestamp: timestamp, IDSession: sessionID}
			err = m.db.InsertSonarData(data)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
}
func (m *MQTTClient) Disconnect() {
	m.client.Disconnect(250)
}
