package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
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

	// Adjusted the URL if the service is in another container
	req, err := http.NewRequest("POST", "http://code-execution-service:8080", bytes.NewBufferString(code))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create request to code execution service: %v", err), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "text/plain")
	httpClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send request to code execution service: %v", err), http.StatusInternalServerError)
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
