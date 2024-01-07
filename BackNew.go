package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PageVariables struct {
	Title string
}

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

type ResponseStatus struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

var newUserChannel = make(chan struct{}, 1)

func main() {
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

func AdminPage(w http.ResponseWriter, r *http.Request) {
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
				newDate := bson.M{
					"login":    requestJSON.Login,
					"email":    requestJSON.Email,
					"address":  requestJSON.Address,
					"number":   requestJSON.PhoneNumber,
					"password": requestJSON.Password}

				updateResult, err := collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, newDate)
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

	//_______________________________admin actions______________________________________
	var results []*ResponseAdmin
	filter := bson.M{}

	cur, err := collection.Find(context.TODO(), filter)
	fmt.Println(cur)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {

		var elem ResponseAdmin
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(elem)

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

// Обработка данных регистрации
// username := r.FormValue("username")
// password := r.FormValue("password")
// email := r.FormValue("email")
// phone := r.FormValue("phone")
// address := r.FormValue("address")

// // Пример: Вставка данных в коллекцию users
// collection := client.Database("mydb").Collection("users")
// _, err = collection.InsertOne(context.TODO(), map[string]interface{}{
// 	"username": username,
// 	"password": password,
// 	"email":    email,
// 	"phone":    phone,
// 	"address":  address,
// })
// if err != nil {
// 	http.Error(w, err.Error(), http.StatusInternalServerError)
// 	return
// }
// // После успешной регистрации, отправляем уведомление о новом пользователе
// newUserChannel <- struct{}{}

// http.Redirect(w, r, "/admin", http.StatusSeeOther)
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
