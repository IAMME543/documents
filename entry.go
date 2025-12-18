package main

import ("fmt";
		 "os";
		"log";
		"net/http")

type Page struct {
	Title string
	Body []byte
}



func loadPage(title string) (*Page, error) {
	filename := "public/" + title + ".html"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
    p, err := loadPage("index")
	if err != nil {
		log.Fatal("Page not found")
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
    fmt.Fprintf(w, "%s", p.Body)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServeTLS("0.0.0.0:443", "cert.pem", "key.pem", nil))

}