package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/internal", internalPageHandler)
	// Obsluga statycznych plikow
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.Handle("/", http.FileServer(http.Dir("./views/")))
	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/logout", logoutHandler)

	// Ustawienie middleware - logger
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	http.ListenAndServe(":8080", loggedRouter)
}
