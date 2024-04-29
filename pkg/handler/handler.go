package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/JyotinderSingh/trytok/pkg/docker"
	"github.com/docker/docker/api/types/container"
)

// ExecuteCode handles the /execute endpoint.
func ExecuteCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is accepted", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	code := string(body) // Assume the body contains the code directly.

	// Container setup (assuming docker is some package managing Docker operations)
	containerId, cli, err := docker.CreateContainer()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create container: %v", err), http.StatusInternalServerError)
		return
	}
	defer cli.ContainerStop(r.Context(), containerId, container.StopOptions{})     // Assuming nil is a valid option
	defer cli.ContainerRemove(r.Context(), containerId, container.RemoveOptions{}) // Assuming nil is a valid option

	// Adjusted the URL if the service is in another container
	req, err := http.NewRequest("POST", "http://localhost:5500", bytes.NewBufferString(code))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create request to container: %v", err), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "text/plain")
	httpClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send request to container: %v", err), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	output, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read response body: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Output: %s", output)
}
