package handler

import (
	"fmt"
	"github.com/u2takey/ffmpeg-go"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// FetchAndConvertVideoHandler fetches the video stream from a test URL and converts it to MP4
func FetchAndConvertVideoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//https://test-streams.mux.dev/x36xhzz/x36xhzz.m3u8
		// Fetch the stream from the test URL
		resp, err := http.Get("http://192.168.36.10:7000/")
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching video stream: %v", err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Define the directory where you want to save the files
		outputDir := "videos"

		// Ensure the directory exists
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			http.Error(w, fmt.Sprintf("Error creating directory: %v", err), http.StatusInternalServerError)
			return
		}

		// Create the raw video file in the specified directory
		rawFilePath := filepath.Join(outputDir, "output.jpeg")
		rawFile, err := os.Create(rawFilePath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating raw file: %v", err), http.StatusInternalServerError)
			return
		}
		defer rawFile.Close()

		// Copy the stream to the file
		_, err = io.Copy(rawFile, resp.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error copying stream to file: %v", err), http.StatusInternalServerError)
			return
		}

		// Ensure the file is closed before converting
		rawFile.Close()

		// Define the MP4 file path
		mp4FilePath := filepath.Join(outputDir, "output.mp4")

		// Convert the raw H.264 file to MP4 using ffmpeg-go
		err = ffmpeg_go.Input(rawFilePath).
			Output(mp4FilePath, ffmpeg_go.KwArgs{"c:v": "libx264"}).
			OverWriteOutput().
			Run()

		if err != nil {
			http.Error(w, fmt.Sprintf("Error converting to MP4: %v", err), http.StatusInternalServerError)
			return
		}

		// Verify if the MP4 file was created successfully
		if _, err := os.Stat(mp4FilePath); os.IsNotExist(err) {
			http.Error(w, "MP4 file not found after conversion", http.StatusInternalServerError)
			return
		}

		// Provide the MP4 file as a downloadable response
		w.Header().Set("Content-Disposition", "attachment; filename=output.mp4")
		w.Header().Set("Content-Type", "video/mp4")
		http.ServeFile(w, r, mp4FilePath)
	} else {
		http.Error(w, "Only GET method is supported", http.StatusMethodNotAllowed)
	}
}

func MakeFramesHandler(w http.ResponseWriter, r *http.Request) {
	streamURL := "http://192.168.36.10:7000/"
	outputDir := "./frames/"

	// Create output directory if it doesn't exist
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.Mkdir(outputDir, 0755)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating output directory: %v", err), http.StatusInternalServerError)
			return
		}
	}

	numFrames := 150 // Adjust as needed

	for i := 0; i < numFrames; i++ {
		frameURL := fmt.Sprintf("%sframe_%d.jpg", streamURL, i)

		// Fetch frame from stream
		err := fetchFrame(frameURL, outputDir, i)
		if err != nil {
			fmt.Printf("Error fetching frame %d: %v\n", i, err)
			http.Error(w, fmt.Sprintf("Error fetching frame %d: %v", i, err), http.StatusInternalServerError)
			return
		}

		// Adjust sleep duration to match frame rate
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("Finished capturing frames.")
}

func fetchFrame(url, outputDir string, index int) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Create output file
	fileName := fmt.Sprintf("%s/frame_%04d.jpg", outputDir, index)
	out, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Saved frame %d\n", index)
	return nil
}
