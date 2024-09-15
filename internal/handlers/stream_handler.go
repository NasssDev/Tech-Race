package handlers

import (
	"bufio"
	"fmt"
	"hetic/tech-race/pkg/util"
	"net/http"
	"os/exec"
	"sync"
)

var (
	cmd *exec.Cmd
	mu  sync.Mutex
)

func StartStream() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		if cmd != nil {
			mu.Unlock()
			util.RenderJson(w, http.StatusBadRequest, map[string]string{"status": "error", "message": "Stream is already running"})
			return
		}
		cmd = exec.Command("go", "run", "cmd/stream/main.go")
		mu.Unlock()

		go func() {
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				fmt.Println("Error when starting the stream:", err)
				return
			}

			stderr, err := cmd.StderrPipe()
			if err != nil {
				fmt.Println("Error when starting the stream:", err)
				return
			}

			if err := cmd.Start(); err != nil {
				fmt.Println("Error when starting the stream:", err)
				return
			}

			go func() {
				scanner := bufio.NewScanner(stderr)
				for scanner.Scan() {
					fmt.Println(scanner.Text())
				}
			}()

			go func() {
				scanner := bufio.NewScanner(stdout)
				for scanner.Scan() {
					fmt.Println(scanner.Text())
				}
			}()

			if err := cmd.Wait(); err != nil {
				fmt.Println("Error when starting the stream:", err)
				return
			}

			mu.Lock()
			cmd = nil
			mu.Unlock()
		}()

		util.RenderJson(w, http.StatusOK, map[string]string{"status": "success", "message": "Recording stream is on", "stream": "active"})
	}
}

func StopStream() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if cmd == nil {
			util.RenderJson(w, http.StatusBadRequest, map[string]string{"status": "error", "message": "No stream is running"})
			return
		}
		if err := cmd.Process.Kill(); err != nil {
			util.RenderJson(w, http.StatusInternalServerError, map[string]string{"status": "error", "message": "Failed to stop the stream"})
			return
		}
		cmd = nil
		util.RenderJson(w, http.StatusOK, map[string]string{"status": "success", "message": "Stream has been stopped"})
	}
}
