package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/huytd/mss/source"
)

type (
	Map map[string]interface{}
)

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

func streamFunc(w http.ResponseWriter, r *http.Request) {
	sourceUrl := r.FormValue("url")
	if !strings.Contains(sourceUrl, "http://") {
		sourceUrl = "http://" + sourceUrl
	}
	log.Println("Requested: ", sourceUrl)
	targetUrl := source.GetURL(sourceUrl)
	log.Println("Found: ", targetUrl)

	data, err := json.Marshal(Map{
		"url": targetUrl,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func searchFunc(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("query")
	log.Print("Search for: ", query)
	searchTerm := strings.Replace(query, " ", "+", -1)
	searchContent := source.GetContent("http://search.chiasenhac.vn/search.php?s=" + searchTerm)
	matches := source.ParseRegExAll(searchContent, `\<div\ class\=\"tenbh\"\>\s*\<p\>\<a\ href\=\"(.*)\"\ class.*\>(.*)\<\/a\>\<\/p>\s*\<p\>(.*)\<\/p>\s*\<\/div\>`)
	data, err := json.Marshal(Map{
		"content": matches,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func main() {
	port := envString("PORT", "3333")

	fs := http.FileServer(http.Dir("www"))

	http.Handle("/", fs)
	http.HandleFunc("/stream", streamFunc)
	http.HandleFunc("/search", searchFunc)

	log.Println("Server running at http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
