package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Page struct {
	Title string
	Body  []byte
}
type SaveRequest struct {
	content string `json:"content"`
}

func loadPage(title string) (*Page, error) {
	filename := "templates/" + title + ".html"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil

}

func saveDocument(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var req SaveRequest
	if err := json.Unmarshal(body, &req.content); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	log.Println(req.content)

	err = os.WriteFile("storage/data.txt", []byte(req.content), 0644)
	if err != nil {
		http.Error(w, "Failed to write data to file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Wrote to file succsefully")
}

func loadDocument(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("storage/data.txt")
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		http.Error(w, "Faled to get file information", http.StatusInternalServerError)
		return
	}

	data := make([]byte, info.Size())
	_, err = file.Read(data)
	if err != nil {
		http.Error(w, "Error making file contents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(data)

}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	p, err := loadPage("index")
	if err != nil {
		log.Println("Page not found")
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "%s", p.Body)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/")

	switch path {
	case "save":
		saveDocument(w, r)
	case "load":
		loadDocument(w, r)
	default:
		fmt.Fprintf(w, "API call does not exist")

	}
}

func main() {

	log.Println("Server Opened on port 443")

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/api/", apiHandler)
	http.HandleFunc("/", mainHandler)

	err := http.ListenAndServeTLS("0.0.0.0:443", "certs/cert.pem", "certs/key.pem", nil)

	log.Fatalf("ListenAndServeTLS failed: %v", err)
	//log.Fatal((http.ListenAndServe("localhost:8080", nil)))
}
