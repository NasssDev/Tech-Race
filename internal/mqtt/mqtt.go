package mqtt

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
	"hetic/tech-race/internal/models"
	"strconv"
	"strings"
	"time"
)

type MQTTClient struct {
	client      MQTT.Client
	db          models.DatabaseInterface
	isAutopilot bool
	isCollision bool
	valueTrack  int
}

func NewMQTTClient(db models.DatabaseInterface) *MQTTClient {
	opts := MQTT.NewClientOptions().AddBroker("tcp://192.168.16.82:1883")
	client := MQTT.NewClient(opts)
	return &MQTTClient{client: client, db: db}
}

func (m *MQTTClient) ConnectAndSubscribe(isAutopilot bool) error {
	m.isAutopilot = isAutopilot
	//_, err := m.db.GetCurrentSessionID()
	// TODO Why isn't the error being returned
	//if err != nil {
	//	//fmt.Println(err)
	//	return nil
	//}
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
		//fmt.Println(err)
		return
	}
	switch topic {
	case "esp32/track":
		m.valueTrack, err = strconv.Atoi(string(msg.Payload()))
		if err != nil {
			//fmt.Println(err)
			return
		}
		if m.valueTrack < 7 {
			timestamp := time.Now()
			data := models.LineTracking{LineTrackingValue: m.valueTrack, IDSession: sessionID, Timestamp: timestamp}
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
		m.isCollision = false
		if distance < 10 {
			m.isCollision = true
		}
		if m.isCollision == true {
			fmt.Println("Collision detected", m.isCollision, distance)
			timestamp := time.Now()
			data := models.Collision{Distance: distance, IsCollision: m.isCollision, Timestamp: timestamp, IDSession: sessionID}
			err = m.db.InsertSonarData(data)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
	//sessionID, err = m.db.GetCurrentSessionID()
	//if err != nil {
	//	//fmt.Println(err)
	//	return
	//}
	if m.isAutopilot {
		//println("autopilot", m.isCollision, m.valueTrack)
		c, _, err := websocket.DefaultDialer.Dial("ws://192.168.16.10/ws", nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer c.Close()
		var payload map[string]interface{}
		if m.isCollision {
			payload = map[string]interface{}{
				"cmd":  1,
				"data": [4]int{0, 0, 0, 0},
			}
			err = c.WriteJSON(payload)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {

			switch m.valueTrack {
			case 7:
				payload = map[string]interface{}{
					"cmd":  1,
					"data": [4]int{500, 500, 500, 500},
				}
			case 3:
				payload = map[string]interface{}{
					"cmd":  1,
					"data": [4]int{0, 0, 500, 500},
				}
			case 6:
				payload = map[string]interface{}{
					"cmd":  1,
					"data": [4]int{500, 500, 0, 0},
				}
			case 0:
				payload = map[string]interface{}{
					"cmd":  1,
					"data": [4]int{0, 0, 0, 0},
				}
			}
			err = c.WriteJSON(payload)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

}
func (m *MQTTClient) Disconnect() {
	m.client.Disconnect(250)
}
