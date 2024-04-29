package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"
)

func compileAndRunCode(w http.ResponseWriter, r *http.Request) {
	// Read the code from the request
	code, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Write the code to a temp file
	tmpFile, err := os.CreateTemp("", "code-*.cpp")
	if err != nil {
		http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
		return
	}
	defer tmpFile.Close()
	os.WriteFile(tmpFile.Name(), code, 0666)

	// Compile and run the code
	cmd := exec.Command("/ctok/build/ctok", tmpFile.Name())
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, "Compilation or execution failed: "+string(output), http.StatusInternalServerError)
		return
	}

	// Write the output back to the client
	w.Write(output)
}

func main() {
	http.HandleFunc("/", compileAndRunCode)
	http.ListenAndServe(":8080", nil)
}
