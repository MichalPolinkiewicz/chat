package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(res, req)
}

func main() {
	r := newRoom()
	http.Handle("/", &templateHandler{
		filename: "chat.html",
	})
	http.Handle("/room", r)
	go r.run()

	if err := http.ListenAndServe(":8083", nil); err != nil {
		log.Fatal("Listen and Serve error", err)
	}

}
