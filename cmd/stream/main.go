package main

import (
	"bufio"
	"fmt"
	"hetic/tech-race/internal/config"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var cfg = config.LoadStreamInfo()

type Client struct {
	conn   net.Conn
	writer *bufio.Writer
}

type VideoRelay struct {
	clients     map[*Client]bool
	clientsLock sync.Mutex
	broadcast   chan []byte
}

func NewVideoRelay() *VideoRelay {
	return &VideoRelay{
		clients:   make(map[*Client]bool),
		broadcast: make(chan []byte),
	}
}

func (vr *VideoRelay) run() {
	for {
		data := <-vr.broadcast
		vr.clientsLock.Lock()
		for client := range vr.clients {
			_, err := client.writer.Write(data)
			if err != nil {
				log.Printf("Erreur d'envoi au client %s: %v", client.conn.RemoteAddr(), err)
				vr.removeClient(client)
				continue
			}
			err = client.writer.Flush()
			if err != nil {
				log.Printf("Erreur de flush pour le client %s: %v", client.conn.RemoteAddr(), err)
				vr.removeClient(client)
			}
		}
		vr.clientsLock.Unlock()
	}
}

func (vr *VideoRelay) addClient(conn net.Conn) {
	client := &Client{conn: conn, writer: bufio.NewWriter(conn)}
	vr.clientsLock.Lock()
	vr.clients[client] = true
	vr.clientsLock.Unlock()

	// Envoyer les en-têtes HTTP pour le flux MJPEG
	client.writer.WriteString("HTTP/1.1 200 OK\r\n")
	client.writer.WriteString("Access-Control-Allow-Origin: *\r\n")
	client.writer.WriteString(fmt.Sprintf("Content-Type: multipart/x-mixed-replace; boundary=%s\r\n", cfg.StreamBoundary))
	client.writer.WriteString("\r\n")
	client.writer.Flush()

	log.Printf("Nouveau client connecté: %s", conn.RemoteAddr())
}

func (vr *VideoRelay) removeClient(client *Client) {
	vr.clientsLock.Lock()
	delete(vr.clients, client)
	vr.clientsLock.Unlock()
	client.conn.Close()
	log.Printf("Client déconnecté: %s", client.conn.RemoteAddr())
}

func (vr *VideoRelay) handleESP32Connection() {
	for {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", cfg.Esp32Address, cfg.Esp32Port))
		if err != nil {
			log.Printf("Erreur de connexion à l'ESP32: %v. Tentative de reconnexion dans 5 secondes...", err)
			time.Sleep(5 * time.Second)
			continue
		}
		defer conn.Close()

		log.Println("Connecté à l'ESP32")

		reader := bufio.NewReader(conn)
		for {
			// Lire jusqu'à la limite
			_, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					log.Println("Connexion à l'ESP32 fermée. Tentative de reconnexion...")
					break
				}
				log.Printf("Erreur de lecture de la limite: %v", err)
				continue
			}

			// Lire les en-têtes
			headers := make([]string, 0)
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					log.Printf("Erreur de lecture des en-têtes: %v", err)
					break
				}
				if line == "\r\n" {
					break
				}
				headers = append(headers, line)
			}

			// Extraire la longueur du contenu
			var contentLength int
			for _, header := range headers {
				if strings.HasPrefix(header, "Content-Length:") {
					_, err := fmt.Sscanf(header, "Content-Length: %d", &contentLength)
					if err != nil {
						return
					}
					break
				}
			}

			if contentLength == 0 {
				log.Println("Longueur de contenu invalide")
				continue
			}

			// Lire le corps de l'image
			body := make([]byte, contentLength)
			_, err = io.ReadFull(reader, body)
			if err != nil {
				log.Printf("Erreur de lecture du corps de l'image: %v", err)
				continue
			}

			// Construire le message complet
			message := fmt.Sprintf("\r\n--%s\r\n", cfg.StreamBoundary)
			message += strings.Join(headers, "")
			message += "\r\n"

			// Envoyer l'en-tête et le corps aux clients
			vr.broadcast <- []byte(message)
			vr.broadcast <- body
		}
	}
}

func main() {
	videoRelay := NewVideoRelay()
	go videoRelay.run()
	go videoRelay.handleESP32Connection()

	listener, err := net.Listen("tcp", cfg.RelayAddress)
	if err != nil {
		log.Fatal("Erreur d'écoute:", err)
	}
	defer listener.Close()

	log.Printf("Serveur de relais vidéo en écoute sur %s", cfg.RelayAddress)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Erreur d'acceptation de la connexion:", err)
			continue
		}
		go videoRelay.addClient(conn)
	}
}
