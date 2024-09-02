package services

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xi2/xz"
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

const (
	serverIP = "192.168.87.10"
	port     = 7000
	boundary = "--123456789000000000000987654321"
)

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

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, port))
	if err != nil {
		fmt.Println("Error connecting:", err)
		return nil, fmt.Errorf("il y a une erreur dans la requête  de la vidéo au serveur mqtt: %w", err)
	}
	defer conn.Close()

	// use the good binary depends on OS
	ffmpegPath, err := DownloadAndExtractFFMPEG(v.CurrentOS)
	if err != nil {
		fmt.Println("Error during download and extract binary :", err)
		return nil, fmt.Errorf("Error during download and extract binary: %w", err)
	}

	// Chemin du fichier temporaire : relatif à la racine du package cloudinary
	// TODO : Understand why "tmp/video" is relative to the content root (path)
	dir := "tmp/video"

	err = createVideoDir(dir)
	if err != nil {
		fmt.Println("Error when creating video dir :", err)
		return nil, fmt.Errorf("Error when creating video dir : %w", err)
	}

	cmd := exec.Command(ffmpegPath, "-f", "mjpeg", "-i", "-", "-c:v", "libx264", filepath.Join(dir, videoName+".mp4"))
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("Error creating stdin pipe:", err)
		return nil, fmt.Errorf("il y a une erreur lors de la transformation en libx264: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("Error creating stderr pipe:", err)
		return nil, fmt.Errorf("il y a une erreur dans la création de stderr pipe: %w", err)
	}

	// display error message depends on OS
	// TODO : Remove this because it's not necessary and it is deprecated + find a unique message for all OS
	if err := cmd.Start(); err != nil {
		if v.CurrentOS == "windows" {
			fmt.Println("Error starting ffmpeg.exe:", err)
			return nil, fmt.Errorf("erreur dans le lancement de ffmpeg.exe sur windows: %w", err)
		}
		if v.CurrentOS == "darwin" {
			fmt.Println("Error starting ffmpeg-mac:", err)
			return nil, fmt.Errorf("erreur dans le lancement de ffmpeg.exe sur darwin (mac): %w", err)
		}
		if v.CurrentOS == "linux" {
			fmt.Println("Error starting ffmpeg-linux:", err)
			return nil, fmt.Errorf("erreur dans le lancement de ffmpeg.exe sur linux: %w", err)
		}
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
		fmt.Println("Error closing stdin:", err)
		return nil, fmt.Errorf("erreur dans la fermeture du stdin: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for ffmpeg to finish:", err)
	}
	// Upload video to Cloudinary

	// delete video from local
	//err = OS.Remove(filepath.Join(dir, videoName+".mp4"))

	return &RecordingData{
		VideoName:           videoName,
		PathToSuppressVideo: filepath.Join(dir, videoName+".mp4"),
	}, nil

}

// TODO : delete this i dont think it is still usefull
func setPathCheckingOS(currentOS string) string {
	os := currentOS
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
		// Handle ZIP file for macOS
		zipFile, err := OS.Create(filepath.Join(ffmpegDir, "ffmpeg.zip"))
		if err != nil {
			return "", err
		}
		defer zipFile.Close()

		_, err = io.Copy(zipFile, resp.Body)
		if err != nil {
			return "", err
		}

		// Unzip the file
		err = unzip(filepath.Join(ffmpegDir, "ffmpeg.zip"), ffmpegDir)
		if err != nil {
			return "", err
		}

		return filepath.Join(ffmpegDir, "ffmpeg"), nil

		//ffmpegDir, _ = handleZipFile(ffmpegDir, resp)
		//fmt.Println("FFMPEG DIR : ", ffmpegDir)
		//return ffmpegDir, nil
	}

	//Handle .tar.xz files for Windows and Linux
	xzr, err := xz.NewReader(resp.Body, 0)
	if err != nil {
		return "", err
	}

	tr := tar.NewReader(xzr)
	if err != nil {
		return "", err
	}

	for {
		header, err := tr.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return "", err
		}

		target := filepath.Join(ffmpegDir, header.Name)
		fmt.Println("Extracting TARGET :", target)

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := OS.Stat(target); OS.IsNotExist(err) {
				if err := OS.MkdirAll(target, 0755); err != nil {
					return "", err
				}
			}
		case tar.TypeReg:
			f, err := OS.OpenFile(target, OS.O_CREATE|OS.O_RDWR, OS.FileMode(header.Mode))
			if err != nil {
				return "", err
			}

			if _, err := io.Copy(f, tr); err != nil {
				return "", err
			}

			f.Close()
		}
		parts := strings.Split(target, "/")
		if len(parts) > 0 && parts[len(parts)-1] == "ffmpeg" {
			ffmpegDir = target
		}
	}

	return filepath.Join(ffmpegDir), nil
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

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

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

func handleZipFile(ffmpegDir string, resp *http.Response) (string, error) {
	zipFile, err := OS.Create(filepath.Join(ffmpegDir, "ffmpeg.zip"))
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	_, err = io.Copy(zipFile, resp.Body)
	if err != nil {
		fmt.Println("Error file ICIII:", err)
		return "", err
	}

	err = unzip(filepath.Join(ffmpegDir, "ffmpeg.zip"), ffmpegDir)
	if err != nil {
		fmt.Println("Error unzipping file ICIII:", err)
		return "", err
	}

	return filepath.Join(ffmpegDir, "ffmpeg"), nil
}
