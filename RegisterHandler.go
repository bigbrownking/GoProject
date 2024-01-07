package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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
}

type ResponseRegister struct {
	Status int         `json:"status"`
	Login  string      `json:"login"`
	Id     interface{} `json:"id"`
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
		if err != nil || requestJSON.Login == "" || requestJSON.Password == "" || requestJSON.Email == "" || requestJSON.PhoneNumber == "" {
			responseError := ResponseStatus{
				Status:  http.StatusBadRequest,
				Message: "Некорректное JSON-сообщение",
			}
			sendJSONResponse(w, responseError)
			return
		}

		fmt.Printf("Received POST request with message: %s\n", requestJSON)

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

		//________________________Нахождение и insert_____________________________
		var result ResponseRegister
		err = collection.FindOne(context.TODO(), requestJSON.Login).Decode(&result)

		var response ResponseRegister
		if result.Login == "" {
			fmt.Println("Not found Login")
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
			http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.Status)
		w.Write(responseJSON)

	} else {
		// Handle other HTTP methods as needed
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
