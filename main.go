package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"time"
)

type Config struct {
	URLs []string `json:"urls"`
}

type Status struct {
	URL        string
	Online     bool
	StatusCode int
}

var (
	config Config
	status []Status
	tpl    *template.Template
)

func loadConfig() {
	data, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data, &config)
}

func checkURLs() {
	status = status[:0]
	for _, url := range config.URLs {
		resp, err := http.Get(url)
		online := err == nil && resp.StatusCode < 400
		status = append(status, Status{URL: url, Online: online, StatusCode: resp.StatusCode})
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, status)
}

func main() {
	loadConfig()
	tpl = template.Must(template.ParseFiles("index.html"))
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for {
			checkURLs()
			<-ticker.C
		}
	}()
	http.HandleFunc("/", handler)

	// Handle static content
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", nil)
}
