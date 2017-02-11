package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/stalehard/chat/config"
)

var indexTemplate = template.Must(template.ParseFiles("template/index.html"))

func serveStatic(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	indexTemplate.Execute(w, r.Host)
}

func main()  {
	conf, err := config.ReadConfig("config/config.json")
	if err != nil {
		log.Fatal("Read config error:", err)
	}

	hub := newHub()
	go hub.run()


	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r, conf.Limit)
	})

	log.Println("Server start on port:", conf.Port)
	http.HandleFunc("/", serveStatic)
	log.Fatal(http.ListenAndServe(conf.Port, nil))
}
