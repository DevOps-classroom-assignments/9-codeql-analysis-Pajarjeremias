package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
)

const allowedDir = "./safe-files"

func main() {
	http.HandleFunc("/readfile", readFileHandler)
	http.HandleFunc("/exec", execHandler)

	fmt.Println("Listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func readFileHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")

	// Resolve the absolute path of the allowed directory
	absAllowedDir, err := filepath.Abs(allowedDir)
	if err != nil {
		http.Error(w, "Server error", 500)
		return
	}

	// Resolve the absolute path of the requested file
	absPath, err := filepath.Abs(filepath.Join(allowedDir, filename))
	if err != nil || !strings.HasPrefix(absPath, absAllowedDir) {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		http.Error(w, "File not found", 404)
		return
	}
	w.Write(data)
}

func execHandler(w http.ResponseWriter, r *http.Request) {
	cmd := r.URL.Query().Get("cmd")

	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		http.Error(w, "Command failed", 500)
		return
	}
	w.Write(out)
}
