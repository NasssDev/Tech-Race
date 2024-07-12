package services

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"testing"
)

var getGOOS = func() string {
	return runtime.GOOS
}

var httpGet = http.Get

func TestDownloadAndExtractFFMPEG_SuccessfulDownloadAndExtraction(t *testing.T) {
	path, err := DownloadAndExtractFFMPEG()
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	fmt.Println("PATHH :", path)

	//if _, err := os.Stat(path); os.IsNotExist(err) {
	//	t.Errorf("Expected file to exist at path: %s, but it does not", path)
	//}
	// assert equal to ../../bin/ffmpeg-7.0.1-amd64-static/ffmpeg
	assert.Equal(t, path, "../../bin/ffmpeg-7.0.1-amd64-static/ffmpeg")
}

func TestDownloadAndExtractFFMPEG_UnsupportedOperatingSystem(t *testing.T) {
	// Temporarily change the operating system to an unsupported one
	oldGetGOOS := getGOOS
	getGOOS = func() string { return "unsupported" }
	defer func() { getGOOS = oldGetGOOS }()

	_, err := DownloadAndExtractFFMPEG()
	if err == nil {
		t.Errorf("Expected error due to unsupported operating system, but got none")
	}
}

func TestDownloadAndExtractFFMPEG_FailedDownload(t *testing.T) {
	// Temporarily change the operating system to an unsupported one
	oldGetGOOS := getGOOS
	getGOOS = func() string { return "windows" }
	defer func() { getGOOS = oldGetGOOS }()

	// Mock http.Get to return an error
	oldHttpGet := httpGet
	httpGet = func(url string) (*http.Response, error) {
		return nil, errors.New("mocked error")
	}
	defer func() { httpGet = oldHttpGet }()

	_, err := DownloadAndExtractFFMPEG()
	if err == nil {
		t.Errorf("Expected error due to failed download, but got none")
	}
}

func TestDownloadAndExtractFFMPEG_FailedExtraction(t *testing.T) {
	// Temporarily change the operating system to an unsupported one
	oldGetGOOS := getGOOS
	getGOOS = func() string { return "windows" }
	defer func() { getGOOS = oldGetGOOS }()

	// Mock http.Get to return a response with a body that cannot be extracted
	oldHttpGet := httpGet
	httpGet = func(url string) (*http.Response, error) {
		return &http.Response{
			Body: ioutil.NopCloser(strings.NewReader("not a tarball")),
		}, nil
	}
	defer func() { httpGet = oldHttpGet }()

	_, err := DownloadAndExtractFFMPEG()
	if err == nil {
		t.Errorf("Expected error due to failed extraction, but got none")
	}
}
