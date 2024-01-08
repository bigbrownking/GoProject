package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AdminRequest struct {
	Id          string `json:"id"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	PhoneNumber string `json:"number"`
	Address     string `json:"address"`
	Action      string `json:"action"`
}

type ResponseAdmin struct {
	Id          primitive.ObjectID `bson:"_id"`
	Login       string             `json:"login"`
	Password    string             `json:"password"`
	Email       string             `json:"email"`
	PhoneNumber string             `json:"number"`
	Address     string             `json:"address"`
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var requestJSON AdminRequest
	err = json.Unmarshal(body, &requestJSON)
	if err != nil {
		return
	} else {
		//_________________________connect to MongoDb_____________________________________
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		opts := options.Client().ApplyURI("mongodb+srv://Esimgali:LOLRKCjhuCSfTdeY@cluste.vdsc74d.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)
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
		fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
		collection := client.Database("mydb").Collection("users")

		//_______________________________admin actions______________________________________

		if requestJSON.Action == "delete" {
			if requestJSON.Id == "" {
				responseError := ResponseStatus{
					Status:  http.StatusBadRequest,
					Message: "Не указано id",
				}
				sendJSONResponse(w, responseError)
				return
			} else {
				objID, err := primitive.ObjectIDFromHex(requestJSON.Id)
				if err != nil {
					http.Error(w, "Error reading request body", http.StatusInternalServerError)
					return
				}
				collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
				responseError := ResponseStatus{
					Status:  200,
					Message: "Успешно удалено",
				}
				sendJSONResponse(w, responseError)
			}
			return
		}
		if requestJSON.Action == "filter" {
			if requestJSON.Id == "" {
				responseError := ResponseStatus{
					Status:  http.StatusBadRequest,
					Message: "Не указано id",
				}
				sendJSONResponse(w, responseError)
				return
			} else {
				var result ResponseAdmin
				objID, err := primitive.ObjectIDFromHex(requestJSON.Id)
				err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&result)
				responseJSON, err := json.Marshal(result)
				if err != nil {
					http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(200)
				w.Write(responseJSON)
			}
		}
		if requestJSON.Action == "update" {
			if requestJSON.Id == "" {
				responseError := ResponseStatus{
					Status:  http.StatusBadRequest,
					Message: "Не указано id",
				}
				sendJSONResponse(w, responseError)
				return
			} else {
				objID, err := primitive.ObjectIDFromHex(requestJSON.Id)

				update := bson.D{{Key: "$set",
					Value: bson.M{
						"login":    requestJSON.Login,
						"email":    requestJSON.Email,
						"address":  requestJSON.Address,
						"number":   requestJSON.PhoneNumber,
						"password": requestJSON.Password,
					},
				}}

				updateResult, err := collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
				responseJSON, err := json.Marshal(updateResult)
				if err != nil {
					http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(200)
				w.Write(responseJSON)
			}
		}
	}
}

func AdminAll(w http.ResponseWriter, r *http.Request) {
	//_________________________connect to MongoDb_____________________________________
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://Esimgali:LOLRKCjhuCSfTdeY@cluste.vdsc74d.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)
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
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	collection := client.Database("mydb").Collection("users")

	//_______________________________admin actions______________________________________
	var results []*ResponseAdmin
	filter := bson.M{}

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {

		var elem ResponseAdmin
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(context.TODO())
	responseJSON, err := json.Marshal(results)
	if err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	w.Write(responseJSON)
}
