package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var (
	projectName = "RepoSync"
	secret      = os.Getenv("GITHUB_SECRET")
	repoDir     = os.Getenv("REPO_DIR")
	listenAddr  = os.Getenv("LISTEN_ADDR")
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "invalid content type", http.StatusUnsupportedMediaType)
		return
	}

	sig := r.Header.Get("X-Hub-Signature-256")
	if sig == "" {
		http.Error(w, "missing signature", http.StatusBadRequest)
		return
	}

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read error", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(sig)) {
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}

	cmd := exec.Command("git", "-C", repoDir, "pull")
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Printf("[%s] git pull failed: %v\n%s", projectName, err, out)
		http.Error(w, "deploy failed", http.StatusInternalServerError)
		return
	}

	log.Printf("[%s] successfully pulled latest changes in %s", projectName, repoDir)
	fmt.Fprint(w, "OK")
}

func main() {
	if secret == "" || repoDir == "" {
		log.Fatal("GITHUB_SECRET and REPO_DIR environment variables must be set")
	}

	cmd := exec.Command("git", "-C", repoDir, "rev-parse", "--is-inside-work-tree")
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("REPO_DIR %s is not a valid git repository: %v\n%s", repoDir, err, out)
	}

	if listenAddr == "" {
		listenAddr = ":8080"
	}

	http.HandleFunc("/hook", handler)
	log.Printf("[%s] listening on %s", projectName, listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
