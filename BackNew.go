package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PageVariables struct {
	Title string
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Message string `json:"write"`
}

type ResponseLogin struct {
	Status int    `json:"status"`
	Login  string `json:"login"`
	Id     string `json:"id"`
}

type adminRequest struct {
	Message string `json:"message"`
}

type ResponseStatus struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

var newUserChannel = make(chan struct{}, 1)

func main() {
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/register", RegisterHandler)
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		var requestJSON LoginRequest
		err = json.Unmarshal(body, &requestJSON)
		if err != nil || requestJSON.Login == "" {
			responseError := ResponseStatus{
				Status:  http.StatusBadRequest,
				Message: "Некорректное JSON-сообщение",
			}
			sendJSONResponse(w, responseError)
			return
		}

		fmt.Fprintf(os.Stdout, "Received POST request with message: %s\n", requestJSON)

		//_________________________connect to MongoDb_____________________________________
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		opts := options.Client().ApplyURI("mongodb+srv://Esimgali:LOLRKCjhuCSfTdeY@cluste.vdsc74d.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)
		// Create a new client and connect to the server
		client, err := mongo.Connect(context.TODO(), opts)
		if err != nil {
			panic(err)
		}
		defer func() {
			if err = client.Disconnect(context.TODO()); err != nil {
				panic(err)
			}
		}()
		// Send a ping to confirm a successful connection
		if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
			panic(err)
		}
		fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
		collection := client.Database("mydb").Collection("users")

		fmt.Println(requestJSON)
		filter := bson.D{{"login", "okfnkvnfkvdn"}}
		fmt.Println(filter)
		var result ResponseLogin
		err = collection.FindOne(context.TODO(), filter).Decode(&result)
		fmt.Println()

		//___________________________send success response_________________________________________
		response := ResponseLogin{
			Status: http.StatusOK,
			Login:  "fjnvfjkv",
			Id:     "150",
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.Status)
		w.Write(responseJSON)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		var requestJSON RegisterRequest
		err = json.Unmarshal(body, &requestJSON)
		if err != nil || requestJSON.Message == "" {
			responseError := ResponseStatus{
				Status:  http.StatusBadRequest,
				Message: "Некорректное JSON-сообщение",
			}
			sendJSONResponse(w, responseError)
			return
		}

		fmt.Printf("Received POST request with message: %s\n", requestJSON.Message)

		response := ResponseStatus{
			Status:  http.StatusOK,
			Message: "Данные успешно приняты",
		}

		sendJSONResponse(w, response)

		// Подключение к MongoDB
		clientOptions := options.Client().ApplyURI("mongodb+srv://Esimgali:LOLRKCjhuCSfTdeY@cluste.vdsc74d.mongodb.net/?retryWrites=true&w=majority")
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
