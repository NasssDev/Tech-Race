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
		log.Printf("Broadcasting data to %d clients", len(vr.clients)) // Log client count before broadcasting
		vr.clientsLock.Lock()
		for client := range vr.clients {
			_, err := client.writer.Write(data)
			if err != nil {
				log.Printf("Error sending to client %s: %v", client.conn.RemoteAddr(), err)
				vr.removeClient(client)
				continue
			}
			err = client.writer.Flush()
			if err != nil {
				log.Printf("Error flushing to client %s: %v", client.conn.RemoteAddr(), err)
				vr.removeClient(client)
			}
		}
		vr.clientsLock.Unlock()
		log.Println("Broadcast completed.") // Log broadcast completion
	}
}

func (vr *VideoRelay) addClient(conn net.Conn) {
	client := &Client{conn: conn, writer: bufio.NewWriter(conn)}
	vr.clientsLock.Lock()
	vr.clients[client] = true
	vr.clientsLock.Unlock()

	log.Printf("New client connected: %s", conn.RemoteAddr()) // Log new client connection

	// Send HTTP headers for MJPEG stream
	_, err := client.writer.WriteString("HTTP/1.1 200 OK\r\n")
	if err != nil {
		log.Printf("Error writing HTTP header to client %s: %v", conn.RemoteAddr(), err)
		vr.removeClient(client)
		return
	}
	client.writer.WriteString("Access-Control-Allow-Origin: *\r\n")
	client.writer.WriteString(fmt.Sprintf("Content-Type: multipart/x-mixed-replace; boundary=%s\r\n", cfg.StreamBoundary))
	client.writer.WriteString("\r\n")
	err = client.writer.Flush()
	if err != nil {
		log.Printf("Error flushing initial headers to client %s: %v", conn.RemoteAddr(), err)
		vr.removeClient(client)
		return
	}
}

func (vr *VideoRelay) removeClient(client *Client) {
	vr.clientsLock.Lock()
	delete(vr.clients, client)
	vr.clientsLock.Unlock()
	client.conn.Close()
	log.Printf("Client disconnected: %s", client.conn.RemoteAddr()) // Log client disconnection
}

func (vr *VideoRelay) handleESP32Connection() {
	for {
		log.Printf("Attempting to connect to ESP32 at %s:%s...", cfg.Esp32Address, cfg.Esp32Port)
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", cfg.Esp32Address, cfg.Esp32Port))
		if err != nil {
			log.Printf("Error connecting to ESP32: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}
		defer conn.Close()

		log.Println("Connected to ESP32")

		reader := bufio.NewReader(conn)
		for {
			log.Println("Waiting for data from ESP32...")
			// Read until boundary
			_, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					log.Println("ESP32 connection closed. Reconnecting...")
					break
				}
				log.Printf("Error reading boundary from ESP32: %v", err)
				continue
			}

			// Read headers
			log.Println("Reading headers from ESP32...")
			headers := make([]string, 0)
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					log.Printf("Error reading headers from ESP32: %v", err)
					break
				}
				if line == "\r\n" {
					break
				}
				headers = append(headers, line)
			}

			// Extract content length
			var contentLength int
			for _, header := range headers {
				if strings.HasPrefix(header, "Content-Length:") {
					_, err := fmt.Sscanf(header, "Content-Length: %d", &contentLength)
					if err != nil {
						log.Printf("Error parsing Content-Length header: %v", err)
						return
					}
					break
				}
			}

			if contentLength == 0 {
				log.Println("Invalid content length from ESP32, skipping this frame.")
				continue
			}

			// Read image body
			log.Printf("Reading image body of length %d...", contentLength)
			body := make([]byte, contentLength)
			_, err = io.ReadFull(reader, body)
			if err != nil {
				log.Printf("Error reading image body from ESP32: %v", err)
				continue
			}

			// Construct complete message
			message := fmt.Sprintf("\r\n--%s\r\n", cfg.StreamBoundary)
			message += strings.Join(headers, "")
			message += "\r\n"

			log.Println("Sending frame to clients...")
			// Send header and body to clients
			vr.broadcast <- []byte(message)
			vr.broadcast <- body
			log.Println("Frame sent to clients.") // Log successful frame relay
		}
	}
}

func main() {
	videoRelay := NewVideoRelay()
	go videoRelay.run()
	go videoRelay.handleESP32Connection()

	log.Printf("Starting video relay server on %s", cfg.RelayAddress)
	listener, err := net.Listen("tcp", cfg.RelayAddress)
	if err != nil {
		log.Fatal("Error listening:", err)
	}
	defer listener.Close()

	log.Printf("Video relay server listening on %s", cfg.RelayAddress)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		log.Printf("New client attempting to connect: %s", conn.RemoteAddr())
		go videoRelay.addClient(conn)
	}
}
