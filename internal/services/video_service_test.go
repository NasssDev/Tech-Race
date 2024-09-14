package services

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

// Unzip TEST //
func TestUnzip_Success(t *testing.T) {
	// Create a temporary zip file
	zipFile, err := os.CreateTemp("", "test*.zip")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(zipFile.Name())

	zipWriter := zip.NewWriter(zipFile)
	_, err = zipWriter.Create("testfile")
	if err != nil {
		t.Fatal(err)
	}
	err = zipWriter.Close()
	if err != nil {
		t.Fatal(err)
	}

	tempDir, err := os.MkdirTemp("", "ffmpeg")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	err = unzip(zipFile.Name(), tempDir)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	_, err = os.Stat(tempDir + "/testfile")
	assert.NoError(t, err)
}

func TestUnzip_NonExistentFile(t *testing.T) {
	err := unzip("nonexistent.zip", ".")
	assert.Error(t, err)
}

func TestUnzip_InvalidZipFile(t *testing.T) {
	tempFile, err := os.CreateTemp("", "test*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("This is not a zip file")
	if err != nil {
		t.Fatal(err)
	}

	err = unzip(tempFile.Name(), ".")
	assert.Error(t, err)
}

func TestUnzip_InvalidDestination(t *testing.T) {
	zipFile, err := os.CreateTemp("", "test*.zip")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(zipFile.Name())

	zipWriter := zip.NewWriter(zipFile)
	_, err = zipWriter.Create("testfile")
	if err != nil {
		t.Fatal(err)
	}
	err = zipWriter.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = unzip(zipFile.Name(), "/invalid/destination")
	assert.Error(t, err)
}

// HandleZipFile TEST //
func TestHandleZipFile_success(t *testing.T) {
	tempDir, err := os.MkdirTemp(".", "ffmpeg")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	zipFile, err := os.CreateTemp(".", "test*.zip")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(zipFile.Name())

	zipWriter := zip.NewWriter(zipFile)
	_, err = zipWriter.Create("ffmpeg")
	if err != nil {
		t.Fatal(err)
	}
	err = zipWriter.Close()
	if err != nil {
		t.Fatal(err)
	}

	zipData, err := os.ReadFile(zipFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	resp := &http.Response{
		Body: io.NopCloser(bytes.NewReader(zipData)),
	}

	_, err = handleZipFile(tempDir, resp)
	assert.NoError(t, err)

	_, err = os.Stat(filepath.Join(tempDir, "ffmpeg"))
	assert.NoError(t, err)
}

func TestHandleZipFile_FailCreateZip(t *testing.T) {
	ffmpegDir := "/invalid/path"

	resp := &http.Response{
		Body: io.NopCloser(bytes.NewReader([]byte("mock zip data"))),
	}

	_, err := handleZipFile(ffmpegDir, resp)
	assert.Error(t, err)
}

func TestHandleZipFile_FailUnzip(t *testing.T) {
	ffmpegDir, err := os.MkdirTemp(".", "ffmpeg")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(ffmpegDir)

	resp := &http.Response{
		Body: io.NopCloser(bytes.NewReader([]byte("invalid zip data"))),
	}

	_, err = handleZipFile(ffmpegDir, resp)
	assert.Error(t, err)
}

// DownloadAndExtractFFMPEG TEST //
func TestDownloadAndExtractFFMPEG_SuccessfulDownloadAndExtractionForLinux(t *testing.T) {
	path, err := DownloadAndExtractFFMPEG("linux")
	assert.NoError(t, err)
	assert.Contains(t, path, "ffmpeg")
}

func TestDownloadAndExtractFFMPEG_SuccessfulDownloadAndExtractionForWindows(t *testing.T) {
	path, err := DownloadAndExtractFFMPEG("windows")
	assert.NoError(t, err)
	assert.Contains(t, path, "ffmpeg")
}

func TestDownloadAndExtractFFMPEG_SuccessfulDownloadAndExtractionForMac(t *testing.T) {
	path, err := DownloadAndExtractFFMPEG("darwin")
	assert.NoError(t, err)
	assert.Contains(t, path, "ffmpeg")
}

func TestDownloadAndExtractFFMPEG_UnsupportedOperatingSystem(t *testing.T) {
	_, err := DownloadAndExtractFFMPEG("unsupported")
	assert.Error(t, err)
	assert.Equal(t, fmt.Errorf("unsupported operating system: %s", "unsupported"), err)
}
