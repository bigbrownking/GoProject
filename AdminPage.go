package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

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
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer file.Close()
	log.SetOutput(file)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
					logger.WithError(err).Error("Error reading request body")
					http.Error(w, err.Error(), http.StatusInternalServerError)
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
		if requestJSON.Action == "filterLogin" {
			if requestJSON.Login == "" {
				responseError := ResponseStatus{
					Status:  http.StatusBadRequest,
					Message: "Не указан логин",
				}
				sendJSONResponse(w, responseError)
				return
			} else {
				var results []ResponseAdmin
				cursor, err := collection.Find(context.TODO(), bson.M{"login": bson.M{"$regex": primitive.Regex{Pattern: requestJSON.Login, Options: "i"}}})
				if err != nil {
					logger.WithError(err).Error("Error querying daytabase")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer cursor.Close(context.Background())
				for cursor.Next(context.Background()) {
					var result ResponseAdmin
					err := cursor.Decode(&result)
					if err != nil {
						logger.WithError(err).Error("Error decoding database result")
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					results = append(results, result)
				}
				responseJSON, err := json.Marshal(results)
				if err != nil {
					logger.WithError(err).Error("Error encoding JSON response")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(200)
				w.Write(responseJSON)
			}
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
					logger.WithError(err).Error("Error encoding JSON response")
					http.Error(w, err.Error(), http.StatusInternalServerError)
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
					logger.WithError(err).Error("Error encoding JSON response")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(200)
				w.Write(responseJSON)
			}
		}
		if requestJSON.Action == "sort" {
			if requestJSON.Login == "" {
				responseError := ResponseStatus{
					Status:  http.StatusBadRequest,
					Message: "Логин пустой",
				}
				sendJSONResponse(w, responseError)
				return
			}

			sortField := requestJSON.Login
			sortOrder, err := strconv.Atoi(requestJSON.Id)
			sortOptions := bson.D{{Key: sortField, Value: sortOrder}}

			cursor, err := collection.Find(context.TODO(), bson.D{}, options.Find().SetSort(sortOptions))
			if err != nil {
				log.Println(err)
			}
			var results []*ResponseAdmin
			for cursor.Next(context.TODO()) {

				var elem ResponseAdmin
				err := cursor.Decode(&elem)
				if err != nil {
					log.Println(err)
				}

				results = append(results, &elem)
			}
			if err := cursor.Err(); err != nil {
				log.Println(err)
			}
			cursor.Close(context.TODO())
			responseJSON, err := json.Marshal(results)
			if err != nil {
				log.Println("Error encoding JSON response")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)

			w.Write(responseJSON)
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
	log.Println("Pinged your deployment. You successfully connected to MongoDB!")
	collection := client.Database("mydb").Collection("users")

	//_______________________________admin actions______________________________________
	var results []*ResponseAdmin
	filter := bson.M{}

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
	}

	for cur.Next(context.TODO()) {

		var elem ResponseAdmin
		err := cur.Decode(&elem)
		if err != nil {
			log.Println(err)
		}

		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Println(err)
	}
	cur.Close(context.TODO())
	responseJSON, err := json.Marshal(results)
	if err != nil {
		log.Println("Error encoding JSON response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	w.Write(responseJSON)
}
