package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

type mailingRequest struct {
	Email   string `json:"email"`
	Text    string `json:"text"`
	Subject string `json:"subject"`
}

type ResponseMailing struct {
	Email   string `json:"email"`
	Text    string `json:"text"`
	Subject string `json:"subject"`
	Status  int    `json:"status"`
}

func SendMailengTextEmail(w http.ResponseWriter, r *http.Request) {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Failed to open log file:", err)
	}
	defer file.Close()
	log.SetOutput(file)
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	} else {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading request body")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var requestJSON mailingRequest
		err = json.Unmarshal(body, &requestJSON)
		if err != nil || requestJSON.Email == "" || requestJSON.Text == "" {
			responseError := ResponseStatus{
				Status:  http.StatusBadRequest,
				Message: "Некорректное JSON-сообщение",
			}
			sendJSONResponse(w, responseError)
			return
		}

		log.Printf("Received POST request with message: %s\n", requestJSON)
		from := "esimgalikhamitov2005@gmail.com"
		password := "oauc fsxn vnxd paxx"
		SMTPHost := "smtp.gmail.com"
		port := 587
		msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n %s", from, requestJSON.Email, requestJSON.Subject, requestJSON.Text)

		auth := smtp.PlainAuth("", from, password, SMTPHost)

		errLast := smtp.SendMail(fmt.Sprintf("%s:%d", SMTPHost, port), auth, from, []string{requestJSON.Email}, []byte(msg))
		var response ResponseMailing
		log.Println(errLast)

		if errLast != nil {
			response = ResponseMailing{
				Status:  505,
				Email:   requestJSON.Email,
				Subject: requestJSON.Subject,
				Text:    requestJSON.Text,
			}
		} else {
			response = ResponseMailing{
				Status:  http.StatusOK,
				Email:   requestJSON.Email,
				Subject: requestJSON.Subject,
				Text:    requestJSON.Text,
			}
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Println("Error encoding JSON response")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.Status)
		w.Write(responseJSON)
	}
}
