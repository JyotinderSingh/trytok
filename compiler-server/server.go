package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func compileAndRunCode(w http.ResponseWriter, r *http.Request) {
	// Read the code from the request
	code, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Write the code to a temp file
	tmpFile, err := os.CreateTemp("", "code-*.tok")

	if err != nil {
		http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
		return
	}
	os.WriteFile(tmpFile.Name(), code, 0666)
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	absPath, err := filepath.Abs(tmpFile.Name())
	if err != nil {
		http.Error(w, "Failed to get absolute path", http.StatusInternalServerError)
		return
	}

	cmd := exec.Command("docker", "run", "--platform", "linux/x86_64", "--rm", "-i",
		"-v", absPath+":"+absPath, "jyotindersingh/ctok",
		"/ctok", absPath)
	log.Println("Running command: ", cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, "Compilation or execution failed: "+string(output), http.StatusInternalServerError)
		return
	}

	// Write the output back to the client
	w.Write([]byte("Output:" + string(output)))
}

func main() {
	http.HandleFunc("/", compileAndRunCode)
	http.ListenAndServe(":8080", nil)
}
