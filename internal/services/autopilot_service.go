package services

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

var esp32Conn *websocket.Conn
var esp32Connected = false
var autopilotRunning = false
var autopilotStop chan struct{}
var mutex = &sync.Mutex{}

type Message struct {
	Cmd  int         `json:"cmd"`
	Data interface{} `json:"data"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func connectToESP32() {
	var err error
	esp32Conn, _, err = websocket.DefaultDialer.Dial("ws://192.168.31.10/ws", nil) // Replace with your ESP32 IP
	if err != nil {
		fmt.Println("Erreur lors de la connexion à l'ESP32:", err)
		esp32Connected = false
		return
	}
	esp32Connected = true
	fmt.Println("Connecté à l'ESP32")
}

func sendMessageToESP32(message Message) {
	if esp32Connected {
		mutex.Lock()
		defer mutex.Unlock()
		if esp32Conn != nil {
			err := esp32Conn.WriteJSON(message)
			if err != nil {
				fmt.Println("Erreur lors de l'envoi du message à l'ESP32:", err)
				esp32Connected = false
			}
		} else {
			fmt.Println("Connexion ESP32 est nil")
		}
	} else {
		fmt.Println("ESP32 non connecté")
	}
}

func goStraight() {
	msg := Message{Cmd: 1, Data: []int{300, 300, 300, 300}}
	sendMessageToESP32(msg)
}

func pushCar() {
	msg := Message{Cmd: 1, Data: []int{500, 500, 500, 500}}
	sendMessageToESP32(msg)
}

func stop() {
	msg := Message{Cmd: 1, Data: []int{0, 0, 0, 0}}
	sendMessageToESP32(msg)
}

func goBack() {
	msg := Message{Cmd: 1, Data: []int{-50, -50, -50, -50}}
	sendMessageToESP32(msg)
}

func goLeft() {
	msg := Message{Cmd: 1, Data: []int{0, 0, 2000, 2000}}
	sendMessageToESP32(msg)
}

func goRight() {
	msg := Message{Cmd: 1, Data: []int{2000, 2000, 0, 0}}
	sendMessageToESP32(msg)
}

func runAutopilot() {
	if autopilotRunning {
		return
	}

	autopilotStop = make(chan struct{})
	autopilotRunning = true
	go func() {
		// Arrêter la voiture en cas d'arrêt d'autopilot
		defer func() {
			stop()
			autopilotRunning = false
		}()

		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-autopilotStop:
				stopAutopilot()
				return
			case <-ticker.C:
				// Simulate autopilot behavior
				pushCar()
				goStraight()
				time.Sleep(2 * time.Second)
			}
		}
	}()
	println("autopilot started")
}

func stopAutopilot() {
	if autopilotRunning {
		stop()
		close(autopilotStop)
		autopilotRunning = false
		println("Autopilot stopped")
	} else {
		println("Autopilot not running")
	}
}
