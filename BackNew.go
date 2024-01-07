package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type PageVariables struct {
	Title string
}

type ResponseStatus struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

var newUserChannel = make(chan struct{}, 1)

func main() {
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/admin", AdminPage)
	http.HandleFunc("/admin/all", AdminAll)
	fmt.Println("Server listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	pageVariables := PageVariables{
		Title: "Registration Page",
	}

	tmpl, err := template.ParseFiles("Registration.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageVariables)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func sendJSONResponse(w http.ResponseWriter, response ResponseStatus) {
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	w.Write(responseJSON)
}
