package services

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	OS "os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

const (
	serverIP = "192.168.1.10"
	port     = 7000
	boundary = "--123456789000000000000987654321"
)

type VideoService struct {
	IsRecording bool
}

func NewVideoService() *VideoService {
	return &VideoService{
		IsRecording: false,
	}
}

func (v *VideoService) StartRecording(sessionService *SessionService) {

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, port))
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	// use the good binary depends on OS
	ffmpegPath := setPathCheckingOS()

	dir := "tmp/video"

	// use this name to save the video in CLOUDINARY
	videoName := time.Now().Format("2006-01-02T15:04:05")

	createVideoDir(dir)

	cmd := exec.Command(ffmpegPath, "-f", "mjpeg", "-i", "-", "-c:v", "libx264", filepath.Join(dir, videoName+".mp4"))
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

	// display error message depends on OS
	if err := cmd.Start(); err != nil {
		if runtime.GOOS == "windows" {
			fmt.Println("Error starting ffmpeg.exe:", err)
			return
		}
		if runtime.GOOS == "darwin" {
			fmt.Println("Error starting ffmpeg-mac:", err)
			return
		}
		if runtime.GOOS == "linux" {
			fmt.Println("Error starting ffmpeg-linux:", err)
			return
		}
	}

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	reader := bufio.NewReader(conn)
	buffer := bytes.Buffer{}

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
					fmt.Println("Error writing to ffmpeg:", err)
					return
				}
			}
		}
		fmt.Println("Video is recording ...")
		v.IsRecording, _ = sessionService.IsSessionActive()
		if !v.IsRecording {
			fmt.Println("Session stopped - Video recording stopped")
			break
		}

	}

	err = stdin.Close()

	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for ffmpeg to finish:", err)
	}
	// Upload video to Cloudinary
	//err = UploadVideoToCloudinary("http://localhost:8045/upload-video", filepath.Join(dir, videoName+".mp4"), videoName)
	//if err != nil {
	//	fmt.Println("Error uploading video:", err)
	//
	//}
	// delete video from local
	//err = OS.Remove(filepath.Join(dir, videoName+".mp4"))
}

func setPathCheckingOS() string {
	os := runtime.GOOS
	fmmpegDir := "bin/ffmpeg"
	if os == "windows" {
		fmt.Println("Windows OS")
		return filepath.Join(fmmpegDir, "/ffmpeg.exe")
	}
	if os == "linux" {
		fmt.Println("Linux OS")
		return filepath.Join(fmmpegDir, "/ffmpeg-linux")
	}
	if os == "darwin" {
		fmt.Println("Mac OS")
		return filepath.Join(fmmpegDir, "/ffmpeg-mac")
	}
	return ""
}

func createVideoDir(dir string) {
	if _, err := OS.Stat(dir); OS.IsNotExist(err) {
		fmt.Println(dir, "does not exist")
		err := OS.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
	} else {
		fmt.Println("The provided directory named", dir, "exists")
	}
}

// UploadVideoToCloudinary videoUrl : url relative au pkg/other/cloudinary
func UploadVideoToCloudinary(uploadURL string, videoURL string, videoID string) error {
	reqURL := fmt.Sprintf("%s?url=%s&id=%s", uploadURL, videoURL, videoID)
	println("package cloudinary appelé: ", reqURL)
	resp, err := http.Get(reqURL)
	if err != nil {
		return fmt.Errorf("Il y a une erreur dans la requête http: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Il y a une erreur dans l'upload de la vidéo, statut: %s, body: %s", resp.Status, string(body))
	}

	fmt.Println("Les vidéos sont sur Cloudinary")
	return nil

}
