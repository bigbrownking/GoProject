package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RegisterRequest struct {
	Login       string `json:"login"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	PhoneNumber string `json:"number"`
	Address     string `json:"address"`
	IsAdmin     bool   `json:"isAdmin"`
}

type ResponseRegister struct {
	Status int         `json:"status"`
	Login  string      `json:"login"`
	Id     interface{} `json:"id"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Failed to open log file:", err)
	}
	defer file.Close()
	log.SetOutput(file)
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading request body")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var requestJSON RegisterRequest
		err = json.Unmarshal(body, &requestJSON)
		if err != nil || requestJSON.Login == "" || requestJSON.Password == "" || requestJSON.Email == "" || requestJSON.PhoneNumber == "" {
			responseError := ResponseStatus{
				Status:  http.StatusBadRequest,
				Message: "Некорректное JSON-сообщение",
			}
			sendJSONResponse(w, responseError)
			return
		}

		log.Printf("Received POST request with message: %s\n", requestJSON)

		//_________________________connect to MongoDb_____________________________________
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		opts := options.Client().ApplyURI("mongodb+srv://myAtlasDBUser:111@myatlasclusteredu.z25a02h.mongodb.net/?retryWrites=true&w=majority&appName=myAtlasClusterEDU").SetServerAPIOptions(serverAPI)
		client, err := mongo.Connect(context.TODO(), opts)
		if err != nil {
			panic(err)
		}
		defer func() {
			if err = client.Disconnect(context.TODO()); err != nil {
				panic(err)
			}
		}()
		if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
			panic(err)
		}
		log.Println("Pinged your deployment. You successfully connected to MongoDB!")
		collection := client.Database("mydb").Collection("users")

		//________________________Find and insert_____________________________
		var result ResponseRegister
		err = collection.FindOne(context.TODO(), bson.M{"login": requestJSON.Login}).Decode(&result)
		log.Println(result)
		var response ResponseRegister
		if result.Login != "" {
			log.Println("Login exist")
			response = ResponseRegister{
				Status: 505,
				Login:  requestJSON.Login,
			}
		} else {
			insertResult, err := collection.InsertOne(context.TODO(), requestJSON)
			if err != nil {
				log.Fatal(err)
			}
			response = ResponseRegister{
				Status: http.StatusOK,
				Login:  requestJSON.Login,
				Id:     insertResult.InsertedID,
			}
		}

		//___________________________send success response_________________________________________

		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Println("Error encoding JSON response")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.Status)
		w.Write(responseJSON)

	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
