package services

import (
	"encoding/json"
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

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Erreur lors de l'upgrade:", err)
		return
	}
	defer ws.Close()

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Erreur lors de la lecture du message:", err)
			break
		}
		fmt.Printf("Message reçu: %+v\n", msg)
	}
}

func handleFront(w http.ResponseWriter, r *http.Request) {
	msg := Message{Cmd: 1, Data: []int{500, 500, 500, 500}}
	sendMessageToESP32(msg)
	json.NewEncoder(w).Encode(msg)
}

func stopCar(w http.ResponseWriter, r *http.Request) {
	msg := Message{Cmd: 1, Data: []int{0, 0, 0, 0}}
	sendMessageToESP32(msg)
	json.NewEncoder(w).Encode(msg)
}

func handleBack(w http.ResponseWriter, r *http.Request) {
	msg := Message{Cmd: 1, Data: []int{-500, -500, -500, -500}}
	sendMessageToESP32(msg)
	json.NewEncoder(w).Encode(msg)
}

func handleLeft(w http.ResponseWriter, r *http.Request) {
	msg := Message{Cmd: 1, Data: []int{0, 0, 2000, 2000}}
	sendMessageToESP32(msg)
	json.NewEncoder(w).Encode(msg)
}

func handleRight(w http.ResponseWriter, r *http.Request) {
	// FL BL FR BR
	msg := Message{Cmd: 1, Data: []int{2000, 2000, 0, 0}}
	sendMessageToESP32(msg)
	json.NewEncoder(w).Encode(msg)
}

func front() {
	msg := Message{Cmd: 1, Data: []int{500, 500, 500, 500}}
	sendMessageToESP32(msg)
}

func stop() {
	msg := Message{Cmd: 1, Data: []int{0, 0, 0, 0}}
	sendMessageToESP32(msg)
}

func back() {
	msg := Message{Cmd: 1, Data: []int{-500, -500, -500, -500}}
	sendMessageToESP32(msg)
}

func left() {
	msg := Message{Cmd: 1, Data: []int{0, 0, 2000, 2000}}
	sendMessageToESP32(msg)
}

func right() {
	msg := Message{Cmd: 1, Data: []int{2000, 2000, 0, 0}}
	sendMessageToESP32(msg)
}

func runAutopilot() {
	if autopilotRunning {
		//w.WriteHeader(http.StatusBadRequest)
		//w.Write([]byte("Autopilot already running"))
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
				stop()
				return
			case <-ticker.C:
				// Simulate autopilot behavior
				front()
				time.Sleep(1 * time.Second)
				left()
				time.Sleep(1 * time.Second)
				front()
				time.Sleep(1 * time.Second)
				right()
			}
		}
	}()
	println("autopilot started")
}

func stopAutopilot(w http.ResponseWriter, r *http.Request) {
	if autopilotRunning {
		stop()
		close(autopilotStop)
		autopilotRunning = false
		w.Write([]byte("Autopilot stopped"))
	} else {
		w.Write([]byte("Autopilot not running"))
	}
}
