package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os/exec"
	"runtime"
	"time"
)

const (
	serverIP = "192.168.52.10"
	port     = 7000
	boundary = "--123456789000000000000987654321"
)

func main() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, port))
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	//cmd := exec.Command("./ffmpeg-linux", "-f", "mjpeg", "-i", "-", "output.mp4")
	os := runtime.GOOS
	var ffmpegPath string
	if os == "windows" {
		fmt.Println("Windows OS")
		ffmpegPath = "bin/ffmpeg/ffmpeg-windows.exe"
	}
	if os == "linux" {
		fmt.Println("Linux OS")
		ffmpegPath = "bin/ffmpeg/ffmpeg-linux"
	}
	if os == "darwin" {
		fmt.Println("Mac OS")
		ffmpegPath = "bin/ffmpeg/ffmpeg-mac"
	}
	cmd := exec.Command(ffmpegPath, "-f", "mjpeg", "-i", "-", "output.mp4")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("Error creating stdin pipe:", err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("Error creating stderr pipe:", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting ffmpeg-linux:", err)
		return
	}

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	reader := bufio.NewReader(conn)
	buffer := bytes.Buffer{}

	startTime := time.Now()

	for {
		chunk, err := reader.ReadBytes(byte('\n'))
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from connection:", err)
			}
			break
		}

		buffer.Write(chunk)

		for {
			boundaryIndex := bytes.Index(buffer.Bytes(), []byte(boundary))
			if boundaryIndex == -1 {
				break
			}

			frame := buffer.Bytes()[:boundaryIndex]
			buffer.Next(boundaryIndex + len(boundary))

			headerEnd := bytes.Index(frame, []byte("\r\n\r\n"))
			if headerEnd == -1 {
				continue
			}

			frameData := frame[headerEnd+4:]
			if len(frameData) > 0 {
				if _, err := stdin.Write(frameData); err != nil {
					fmt.Println("Error writing to ffmpeg-linux:", err)
					return
				}
			}
		}

		if time.Since(startTime) > 40*time.Second {
			break
		}

	}

	stdin.Close()
	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for ffmpeg-linux to finish:", err)
	}
}
