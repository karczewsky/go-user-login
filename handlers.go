package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/securecookie"
)

const internalPage = `
<h1>Internal</h1>
<hr>
<small>User: %s</small>
<form method="post" action="/logout">
    <button type="submit">Logout</button>
</form>
`

type user struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func apiHandler(response http.ResponseWriter, request *http.Request) {
	group := &user{
		ID:   1,
		Name: "Adam",
		Role: "admin",
	}

	res1B, err := json.Marshal(group)

	if err != nil {
		fmt.Println(err)
		return
	}

	userName := getUserName(request)
	if userName != "" {
		fmt.Fprintf(response, string(res1B))
	} else {
		http.Redirect(response, request, "/", 302)
	}
}

func internalPageHandler(response http.ResponseWriter, request *http.Request) {
	userName := getUserName(request)
	if userName != "" {
		fmt.Fprintf(response, internalPage, userName)
	} else {
		http.Redirect(response, request, "/", 302)
	}
}

func loginHandler(response http.ResponseWriter, request *http.Request) {
	login := request.FormValue("username")
	pass := request.FormValue("password")

	redirectURL := "/"
	if login != "" && pass != "" {
		// Walidacja danych

		if login == "admin" && pass == "admin" {
			setSession(login, response)
			redirectURL = "/internal"
		}
	}
	http.Redirect(response, request, redirectURL, 302)
}

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}

// Session, Cookie related functions
var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"username": userName,
	}

	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := http.Cookie{
			Name:     "session",
			Value:    encoded,
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(response, &cookie)
		fmt.Printf("Creating cookie -> user: %s", userName)
	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["username"]
		}
	}
	fmt.Printf("LOL %s", userName)
	return userName
}
