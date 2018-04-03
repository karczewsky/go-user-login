package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("./views/")))
	r.HandleFunc("/internal", internalPageHandler)
	// Obsluga statycznych plikow
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.HandleFunc("/login", loginHandler)

	// Ustawienie middleware - logger
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	http.ListenAndServe(":8080", loggedRouter)
}

func internalPageHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "XD LOL")
}

func loginHandler(response http.ResponseWriter, request *http.Request) {
	login := request.FormValue("username")
	pass := request.FormValue("password")

	redirectURL := "/"
	if login != "" && pass != "" {
		// Walidacja danych

		if login == "admin" && pass == "admin" {
			redirectURL = "/internal"
		}
	}
	http.Redirect(response, request, redirectURL, 302)
}
