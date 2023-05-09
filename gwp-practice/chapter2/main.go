package main

import (
	"data"
	"net/http"
	"time"
)

func main() {
	p("ChitChat", version(), "started at", config.Address)

	mux := http.NewServeMux()
	file := http.FileServer(http.Dir(config.Static))
	mux.Handle("/static/", http.StripPrefix("/static/", file))

	mux.HandleFunc("/", index)
	mux.HandleFunc("/err", err)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/signup", signup)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/signup_account", signupAccount)
	mux.HandleFunc("/authenticate", authenticate)

	mux.HandleFunc("/thread/new", newThread)
	mux.HandleFunc("/thread/create", createThread)
	mux.HandleFunc("/thread/post", postThread)
	mux.HandleFunc("/thread/read", readThread)

	server := &http.Server{
		Addr:           config.Address,
		Handler:        mux,
		ReadTimeout:    time.Duration(config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	data.Init()
	server.ListenAndServe()
}

//func index(w http.ResponseWriter, r *http.Request) {
//	files := []string{
//		"templates/layout.html",
//		"templates/navbar.html",
//		"templates/index.html",
//	}
//	templates := template.Must(template.ParseFiles(files...))
//	threads, err := data.Threads()
//	if err == nil {
//		templates.ExecuteTemplate(w, "layout", threads)
//	}
//}
