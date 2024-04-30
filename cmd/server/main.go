package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	serverAddress           = ":8080"
	codeExecutionServiceURL = "http://code-execution-service:8080"
	requestTimeout          = 10 * time.Second
	contentType             = "text/plain"
)

func main() {
	http.HandleFunc("/execute", executeCodeHandler)
	log.Println("Server started on", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, nil))
}

// executeCodeHandler handles the /execute endpoint.
func executeCodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is accepted", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	response, err := executeCode(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Output: %s", response)
}

// executeCode sends the code to the code execution service and returns the output.
func executeCode(code []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, codeExecutionServiceURL, bytes.NewBuffer(code))
	if err != nil {
		return nil, fmt.Errorf("failed to create request to code execution service: %v", err)
	}

	req.Header.Set("Content-Type", contentType)
	httpClient := &http.Client{Timeout: requestTimeout}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to code execution service: %v", err)
	}
	defer resp.Body.Close()

	output, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return output, nil
}
