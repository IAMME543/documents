package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

type Page struct {
	Title string
	Body  []byte
}
type SaveRequest struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
type LoadRequest struct {
	Id int64 `json:"id"`
}
type IndexList struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

func loadPage(title string) (*Page, error) {
	filename := "templates/" + title + ".html"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil

}

func createDocument(w http.ResponseWriter, r *http.Request) {
	id := addToDB("")

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatInt(id, 10)))
}

func updateDocument(w http.ResponseWriter, r *http.Request) {

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
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	log.Printf("id %d: title %s", req.Id, req.Title)
	updateDB(req.Id, req.Title, w)

	filename := fmt.Sprintf("storage/%d.json", req.Id)

	err = os.WriteFile(filename, body, 0644)
	if err != nil {
		http.Error(w, "Failed to write data to file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Updated file succsefully")
}

func loadDocument(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var req LoadRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	filename := fmt.Sprintf("storage/%d.json", req.Id)

	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(data)

}
func listDocument(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite", "storage/database/index.db")
	if err != nil {
		http.Error(w, "SQL open failed", http.StatusInternalServerError)
		panic(err)
	}

	rows, err := db.Query(`SELECT id, title FROM "index"`)
	if err != nil {
		http.Error(w, "SQL query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var indexlist []IndexList

	for rows.Next() {
		var item IndexList
		if err := rows.Scan(&item.Id, &item.Title); err != nil {
			http.Error(w, "SQL scan failed", http.StatusInternalServerError)
			return
		}
		indexlist = append(indexlist, item)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Row iteration failed", http.StatusInternalServerError)
		return
	}

	log.Println(indexlist)

	w.Header().Set("Content-Type", "applicaiton/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(indexlist)

	if err != nil {
		http.Error(w, "JSON encoder failed", http.StatusInternalServerError)
		panic(err)
	}

}

func parseArchive(w http.ResponseWriter, r *http.Request) {

	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	id := parts[1]

	log.Println(id)

	filename := fmt.Sprintf("storage/%s.json", id)

	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	var req SaveRequest
	if err := json.Unmarshal(data, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprintln(w, "<!doctype html>")
	fmt.Fprintln(w, "<html><head>")
	fmt.Fprintf(w, "<title>%s</title>\n", html.EscapeString(req.Title))
	fmt.Fprintf(w, `<link rel="stylesheet" href="/static/style.css" />`)
	fmt.Fprintln(w, "</head><body>")
	fmt.Fprintf(w, `<p>This is a read only archive version of this document<p>
	<hr>
	<div class="headercontainer">
		<button> </button>
		<div class="spacer"></div>
		<h1>%s</h1>
	</div>`, html.EscapeString(req.Title))

	// Very simple paragraph splitting
	for _, para := range strings.Split(req.Content, "\n\n") {
		fmt.Fprintf(w, "<p>%s</p>\n", html.EscapeString(para))
	}

	fmt.Fprintln(w, "</body></html>")

}

func initDB() {
	db, err := sql.Open("sqlite", "storage/database/index.db")

	if err != nil {
		log.Fatalf("open failed: %#v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("ping failed: %#v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS "index" (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL
		);
	`)
	if err != nil {
		panic(err)
	}
	defer db.Close()

}

func addToDB(title string) int64 {
	db, err := sql.Open("sqlite", "storage/database/index.db")
	if err != nil {
		panic(err)
	}
	res, err := db.Exec(
		`INSERT INTO "index" (title) VALUES (?)`,
		title)

	if err != nil {
		panic(err)
	}

	id, _ := res.LastInsertId()
	defer db.Close()
	return id
}
func updateDB(id int64, title string, w http.ResponseWriter) {
	db, err := sql.Open("sqlite", "storage/database/index.db")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`UPDATE "index" SET title = ? WHERE id = ?`, title, id)
	defer db.Close()

	if err != nil {
		http.Error(w, "Erorr updating DB", http.StatusInternalServerError)
	}

}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	host := strings.Split(r.Host, ":")[0]

	switch host {
	case "masondoesthings.com":
		p, err := loadPage("home")
		if err != nil {
			log.Println("Page not found")
			http.Error(w, "Page not found", http.StatusNotFound)
			return
		}
		fmt.Fprintf(w, "%s", p.Body)
	case "atypingsite.masondoesthings.com":
		var pageName string = "nil"
		path := strings.Trim(r.URL.Path, "/")

		switch {
		case path == "":
			pageName = "index"
		case strings.HasPrefix(path, "editing"):
			pageName = "editing"
		case strings.HasPrefix(path, "archive"):
			parseArchive(w, r)
		}
		if pageName != "nil" {
			p, err := loadPage(pageName)
			if err != nil {
				log.Println("Page not found")
				http.Error(w, "Page not found", http.StatusNotFound)
				return
			}
			fmt.Fprintf(w, "%s", p.Body)
		}

	default:
		http.NotFound(w, r)
		return
	}

}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/")

	switch path {
	case "create":
		createDocument(w, r)
	case "update":
		updateDocument(w, r)
	case "load":
		loadDocument(w, r)
	case "list":
		listDocument(w, r)
	default:
		http.Error(w, "API call does not exist", http.StatusMethodNotAllowed)
		log.Printf("API call not found")

	}
}

func main() {

	log.Println("Server Opened on port 443")

	initDB()

	fsstatic := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fsstatic))

	http.HandleFunc("/api/", apiHandler)
	http.HandleFunc("/", mainHandler)

	err := http.ListenAndServeTLS("0.0.0.0:443", "certs/cert.pem", "certs/key.pem", nil)

	log.Fatalf("ListenAndServeTLS failed: %v", err)

	// log.Fatal((http.ListenAndServe("localhost:8080", nil)))
}
