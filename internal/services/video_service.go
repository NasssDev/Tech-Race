package services

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xi2/xz"
	"hetic/tech-race/internal/config"
	"hetic/tech-race/internal/models"
	"io"
	"log"
	"net"
	"net/http"
	OS "os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var cfg = config.LoadStreamInfo()

type VideoService struct {
	IsRecording bool
	CurrentOS   string
}

func NewVideoService(currentOS string) *VideoService {
	return &VideoService{
		IsRecording: false,
		CurrentOS:   currentOS,
	}
}

type UploadService struct {
	db models.DatabaseInterface
}

func NewUploadService(db models.DatabaseInterface) *UploadService {
	return &UploadService{db: db}
}

type RecordingData struct {
	VideoName           string
	PathToSuppressVideo string
}

func (v *VideoService) StartRecording(sessionService *SessionService) (*RecordingData, error) {
	// use this name to save the video in CLOUDINARY
	videoName := time.Now().Format("2006-01-02T15:04:05")

	conn, err := net.Dial("tcp", cfg.RelayAddress)
	if err != nil {
		return nil, fmt.Errorf("error connecting: %w", err)
	}
	defer conn.Close()

	// retrieve the good binary depends on OS
	ffmpegPath, err := DownloadAndExtractFFMPEG(v.CurrentOS)
	if err != nil {
		return nil, fmt.Errorf("error during download and extract binary: %w", err)
	}

	// path to save videos in tmp/
	// TODO : Understand why "tmp/video" is relative to the content root (Tech-Race/)
	dir := "tmp/video"

	err = createVideoDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error when creating video dir : %w", err)
	}

	cmd := exec.Command(ffmpegPath, "-f", "mjpeg", "-i", "-", "-c:v", "libx264", filepath.Join(dir, videoName+".mp4"))
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("error creating stdin pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("error creating stderr pipe: %w", err)
	}

	// display error message depends on OS
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("error when starting ffmpeg: %w", err)
	}

	//TODO : Decrypt and understand all the code below
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
			boundaryIndex := bytes.Index(buffer.Bytes(), []byte(cfg.StreamBoundary))
			if boundaryIndex == -1 {
				break
			}

			frame := buffer.Bytes()[:boundaryIndex]
			buffer.Next(boundaryIndex + len(cfg.StreamBoundary))

			headerEnd := bytes.Index(frame, []byte("\r\n\r\n"))
			if headerEnd == -1 {
				continue
			}

			frameData := frame[headerEnd+4:]
			if len(frameData) > 0 {
				if _, err := stdin.Write(frameData); err != nil {
					fmt.Println("Error writing to ffmpeg:", err)
					return nil, fmt.Errorf("erreur dans la création de la frameData: %w", err)
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
	if err != nil {
		return nil, fmt.Errorf("error closing stdin: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("error waiting for ffmpeg to finish: %w", err)
	}

	return &RecordingData{
		VideoName:           videoName,
		PathToSuppressVideo: filepath.Join(dir, videoName+".mp4"),
	}, nil
}

func createVideoDir(dir string) error {
	if _, err := OS.Stat(dir); OS.IsNotExist(err) {
		fmt.Println(dir, "does not exist")
		err := OS.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return err
		}
	} else {
		fmt.Println("The provided directory named", dir, "exists")
	}
	return nil
}

// UploadVideoToCloudinary videoUrl : url relative au pkg/other/cloudinary
func UploadVideoToCloudinary(uploadURL string, videoURL string, videoID string) (models.AssetData, error) {

	//test : upload-video?url=../../../tmp/video/2024-07-11T16:29:13.mp4&id=2024-07-11T16:29:13
	url := fmt.Sprintf("%s?url=%s&id=%s", uploadURL, videoURL, videoID)
	println("package cloudinary appelé: ", url)

	cloudinaryClient := http.Client{
		Timeout: 0, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "techrace-cloudinary")

	res, getErr := cloudinaryClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	assetData := models.AssetData{}
	jsonErr := json.Unmarshal(body, &assetData)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	println("json response is ok")

	fmt.Println(assetData.Data.Data.URL)

	return assetData, nil
}

func (u *UploadService) InsertVideo(videoPath string) error {

	sessionID, sesserr := u.db.GetLastSessionID()
	if sesserr != nil {
		println("problem getting session id", sesserr)
		return sesserr
	}

	data := models.Video{VideoURL: videoPath, IDSession: sessionID}
	err := u.db.InsertVideoData(data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	println("insertion de la video en bdd")
	return nil
}

func DownloadAndExtractFFMPEG(currentOS string) (string, error) {
	const (
		windowsLinuxUrl = "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz"
		macOSUrl        = "https://evermeet.cx/ffmpeg/ffmpeg-4.3.1.zip"
	)

	ffmpegDir := "../../bin"

	var url string

	switch currentOS {
	case "windows":
		url = windowsLinuxUrl
	case "linux":
		url = windowsLinuxUrl
	case "darwin":
		url = macOSUrl
	default:
		return "", fmt.Errorf("unsupported operating system: %s", currentOS)
	}

	// download ffmpeg file from the url
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	// Create the ffmpeg directory if it does not exist
	err = OS.MkdirAll(ffmpegDir, 0755)
	if err != nil {
		return "", err
	}

	if currentOS == "darwin" {
		//Handle .zip files for MacOS
		ffmpegDir, err = handleZipFile(ffmpegDir, resp)
		fmt.Println("FFMPEG DIR (after handleZipFile) : ", ffmpegDir)
		return ffmpegDir, nil
	}

	//Handle .tar.xz files for Windows and Linux
	ffmpegDir, err = handleTarXZFile(ffmpegDir, resp)
	fmt.Println("FFMPEG DIR (after handleTarXZFile) : ", ffmpegDir)

	return filepath.Join(ffmpegDir), nil
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// r.File is a slice of *File
	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			if err := OS.MkdirAll(fpath, OS.ModePerm); err != nil {
				return err
			}
			continue
		}

		if err := OS.MkdirAll(filepath.Dir(fpath), OS.ModePerm); err != nil {
			return err
		}

		outFile, err := OS.OpenFile(fpath, OS.O_WRONLY|OS.O_CREATE|OS.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		rc.Close()
		outFile.Close()

		if err != nil {
			return err
		}

	}
	// delete zip file since we dont need it anymore
	if err = OS.Remove(filepath.Join(src)); err != nil {
		return err
	}
	return nil
}

func handleZipFile(ffmpegDir string, resp *http.Response) (string, error) {
	zipFile, err := OS.Create(filepath.Join(ffmpegDir, "ffmpeg.zip"))
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	if _, err = io.Copy(zipFile, resp.Body); err != nil {
		return "", err
	}
	// Unzip the file
	err = unzip(filepath.Join(ffmpegDir, "ffmpeg.zip"), ffmpegDir)
	if err != nil {
		return "", err
	}

	return filepath.Join(ffmpegDir, "ffmpeg"), nil
}

func handleTarXZFile(ffmpegDir string, resp *http.Response) (string, error) {
	xzr, err := xz.NewReader(resp.Body, 0)
	if err != nil {
		return "", err
	}

	tr := tar.NewReader(xzr)

	for {
		header, err := tr.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return "", err
		}

		target := filepath.Join(ffmpegDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := OS.Stat(target); OS.IsNotExist(err) {
				if err := OS.MkdirAll(target, 0755); err != nil {
					return "", err
				}
			}
		case tar.TypeReg:
			file, err := OS.OpenFile(target, OS.O_CREATE|OS.O_RDWR, OS.FileMode(header.Mode))
			if err != nil {
				return "", err
			}

			if _, err := io.Copy(file, tr); err != nil {
				return "", err
			}

			file.Close()
		}
		parts := strings.Split(target, "/")
		if len(parts) > 0 && parts[len(parts)-1] == "ffmpeg" {
			ffmpegDir = target
		}
	}
	return ffmpegDir, nil
}
