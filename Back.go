package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PageVariables struct {
	Title string
}

var newUserChannel = make(chan struct{}, 1)

func main() {
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/register", RegisterHandler)
	http.ListenAndServe(":8080", nil)
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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Подключение к MongoDB
		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Disconnect(context.TODO())
		// Проверка соединения
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatal(err)
		}
		// Обработка данных регистрации
		username := r.FormValue("username")
		password := r.FormValue("password")
		email := r.FormValue("email")
		phone := r.FormValue("phone")
		address := r.FormValue("address")

		// Пример: Вставка данных в коллекцию users
		collection := client.Database("mydb").Collection("users")
		_, err = collection.InsertOne(context.TODO(), map[string]interface{}{
			"username": username,
			"password": password,
			"email":    email,
			"phone":    phone,
			"address":  address,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// После успешной регистрации, отправляем уведомление о новом пользователе
		newUserChannel <- struct{}{}

		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		// Process registration data here
		// username := r.FormValue("username")
		// password := r.FormValue("password")
		// confirm_password := r.FormValue("confirm_password")
		// email := r.FormValue("email")
		// phone := r.FormValue("phone")
		// address := r.FormValue("address")

		// You can handle the registration data as needed (e.g., store it in a database)

		// Redirect to the admin page after registration
		// http.Redirect(w, r, "/admin", http.StatusSeeOther)
	} else {
		// Handle other HTTP methods as needed
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func AdminPage(w http.ResponseWriter, r *http.Request) {
	// Handle admin page logic here
	fmt.Fprintf(w, "This is the admin page.")
	// Отправляем заголовок, чтобы установить соединение в режиме long polling
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Отправляем пустое событие для установки соединения
	fmt.Fprintf(w, "event: keep-alive\ndata: \n\n")
	w.(http.Flusher).Flush()

	// Ожидаем уведомления о новом пользователе и отправляем его на клиент
	for {
		select {
		case <-r.Context().Done():
			return
		case <-newUserChannel:
			users := getUsersFromDB() // Получите список пользователей из базы данных
			sendDataToClient(w, users)
		}
	}
}
func sendDataToClient(w http.ResponseWriter, data interface{}) {
	// Отправляем данные на клиент в формате Server-Sent Events
	fmt.Fprintf(w, "data: %v\n\n", data)
	w.(http.Flusher).Flush()
}

func getUsersFromDB() []map[string]interface{} {
	// Подключение к MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	// Получение данных из коллекции users
	collection := client.Database("mydb").Collection("users")
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	// Обработка данных
	var users []map[string]interface{}
	for cursor.Next(context.TODO()) {
		var user map[string]interface{}
		if err := cursor.Decode(&user); err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	return users
}
