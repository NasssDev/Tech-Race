package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var esp32Conn *websocket.Conn
var esp32Connected = false
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
	esp32Conn, _, err = websocket.DefaultDialer.Dial("ws://192.168.18.50/ws", nil) // Replace with your ESP32 IP
	if err != nil {
		fmt.Println("Erreur lors de la connexion à l'ESP32:", err)
		return
	}
	esp32Connected = true
	fmt.Println("Connecté à l'ESP32")
}

func sendMessageToESP32(message Message) {
	if esp32Connected {
		mutex.Lock()
		defer mutex.Unlock()
		err := esp32Conn.WriteJSON(message)
		if err != nil {
			fmt.Println("Erreur lors de l'envoi du message à l'ESP32:", err)
			esp32Connected = false
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

func handleBack(w http.ResponseWriter, r *http.Request) {
	msg := Message{Cmd: 1, Data: []int{-500, -500, -500, -500}}
	sendMessageToESP32(msg)
	json.NewEncoder(w).Encode(msg)
}

func handleLeft(w http.ResponseWriter, r *http.Request) {
	msg := Message{Cmd: 1, Data: []int{-500, 500, -500, 500}}
	sendMessageToESP32(msg)
	json.NewEncoder(w).Encode(msg)
}

func handleRight(w http.ResponseWriter, r *http.Request) {
	msg := Message{Cmd: 1, Data: []int{500, -500, 500, -500}}
	sendMessageToESP32(msg)
	json.NewEncoder(w).Encode(msg)
}

func handleCamOn(w http.ResponseWriter, r *http.Request) {
	msg := Message{Cmd: 9, Data: 1}
	sendMessageToESP32(msg)
	json.NewEncoder(w).Encode(msg)
}

func handleCamOff(w http.ResponseWriter, r *http.Request) {
	msg := Message{Cmd: 9, Data: 0}
	sendMessageToESP32(msg)
	json.NewEncoder(w).Encode(msg)
}

func handleHead(w http.ResponseWriter, r *http.Request) {
	var data struct {
		X int `json:"x"`
		Y int `json:"y"`
	}
	json.NewDecoder(r.Body).Decode(&data)
	msg := Message{Cmd: 3, Data: []int{data.X, data.Y}}
	sendMessageToESP32(msg)
	json.NewEncoder(w).Encode(msg)
}

func handleHeadFace(w http.ResponseWriter, r *http.Request) {
	var data struct {
		X int `json:"x"`
	}
	json.NewDecoder(r.Body).Decode(&data)
	msg := Message{Cmd: 2, Data: data.X}
	sendMessageToESP32(msg)
	json.NewEncoder(w).Encode(msg)
}

func main() {
	go connectToESP32()

	router := mux.NewRouter()

	router.HandleFunc("/ws", handleConnections)
	router.HandleFunc("/front", handleFront).Methods("POST")
	router.HandleFunc("/back", handleBack).Methods("POST")
	router.HandleFunc("/left", handleLeft).Methods("POST")
	router.HandleFunc("/right", handleRight).Methods("POST")
	router.HandleFunc("/camon", handleCamOn).Methods("POST")
	router.HandleFunc("/camoff", handleCamOff).Methods("POST")
	router.HandleFunc("/head", handleHead).Methods("POST")
	router.HandleFunc("/headface", handleHeadFace).Methods("POST")

	fmt.Println("Serveur démarré sur le port 8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Erreur lors du démarrage du serveur:", err)
	}
}
