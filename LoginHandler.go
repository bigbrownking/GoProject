package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ResponseLogin struct {
	Status int                `json:"status"`
	Login  string             `json:"login"`
	Id     primitive.ObjectID `bson:"_id"`
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

		//______________________________find login and password______________________
		var result ResponseLogin
		err = collection.FindOne(context.TODO(), requestJSON).Decode(&result)
		fmt.Println(result.Id)

		//___________________________send success response_________________________________________
		response := ResponseLogin{
			Status: http.StatusOK,
			Login:  result.Login,
			Id:     result.Id,
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