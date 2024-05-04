package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	serverPort       = ":8080"
	compilerImage    = "jyotindersingh/ctok"
	maxCodeSize      = 1024 * 1024     // 1 MB
	executionTimeout = 5 * time.Second // 5 seconds
)

func main() {
	http.HandleFunc("/", compileAndRunCodeHandler)
	log.Fatal(http.ListenAndServe(serverPort, nil))
}

func compileAndRunCodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > maxCodeSize {
		http.Error(w, "Code size is too large", http.StatusBadRequest)
		return
	}

	code, err := io.ReadAll(r.Body)
	if err != nil {
		handleError(w, "Failed to read request body", err)
		return
	}

	output, err := compileAndRunCode(code)
	if err != nil {
		handleError(w, "Compilation or execution failed: "+string(output), err)
		return
	}

	w.Write(output)
}

func compileAndRunCode(code []byte) ([]byte, error) {
	tmpFile, err := os.CreateTemp("", "code-*.tok")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())

	if err := os.WriteFile(tmpFile.Name(), code, 0666); err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(tmpFile.Name())
	if err != nil {
		return nil, err
	}

	// Create a context with the defined timeout
	ctx, cancel := context.WithTimeout(context.Background(), executionTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "run", "--platform", "linux/x86_64",
		"--rm", "-i", "-m", "65m", "--cpus", "0.1",
		"-v", absPath+":"+absPath, compilerImage,
		"/ctok", absPath)
	log.Printf("Running command: %s; with code: %s", cmd.String(), code)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, err
	}

	return output, nil
}

func handleError(w http.ResponseWriter, msg string, err error) {
	log.Printf("%s: %v", msg, err)
	http.Error(w, msg+": "+err.Error(), http.StatusInternalServerError)
}
