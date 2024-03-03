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
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/time/rate"
)

type PageVariables struct {
	Title string
}

type ResponseStatus struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ResponsePosts struct {
	Img   string             `json:"img"`
	Decs  string             `json:"Decs"`
	Price string             `json:"Price"`
	Name  string             `json:"Name"`
	Id    primitive.ObjectID `bson:"_id"`
}

var newUserChannel = make(chan struct{}, 1)

var limiter = rate.NewLimiter(1000, 1)

var CurrentUser ResponseLogin

var logger = logrus.New()

func main() {
	logger.SetFormatter(&logrus.JSONFormatter{})
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer file.Close()
	log.SetOutput(file)
	http.HandleFunc("/", rateLimit(HomePage))
	http.HandleFunc("/Script.js", rateLimit(jsFile))
	http.HandleFunc("/loginPage", rateLimit(LoginPage))
	http.HandleFunc("/AdminPage", rateLimit(AdminPage))
	http.HandleFunc("/main", rateLimit(Main))
	http.HandleFunc("/ProductsPage", rateLimit(ProductsPage))
	http.HandleFunc("/CartPage", rateLimit(CartPage))
	http.HandleFunc("/login", rateLimit(LoginHandler))
	http.HandleFunc("/Products", rateLimit(ProductsHandler))
	http.HandleFunc("/Cart", rateLimit(CartHandler))
	http.HandleFunc("/register", rateLimit(RegisterHandler))
	http.HandleFunc("/admin", rateLimit(AdminHandler))
	http.HandleFunc("/admin/all", rateLimit(AdminAll))
	http.HandleFunc("/auth", rateLimit(SendVerificationCodeEmail))
	http.HandleFunc("/getPosts", rateLimit(getAllPosts))
	http.HandleFunc("/isLogin", rateLimit(isLogin))
	http.HandleFunc("/profile", rateLimit(Profile))
	http.HandleFunc("/logOut", rateLimit(LogOut))

	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}
	log.Println("Server listening on port: 8080")
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("Error starting server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println("Error shutting down server")
	}

	log.Println("Server stopped")
}

func rateLimit(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			log.Println("Rate limit exceeded")
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		h(w, r)
	}
}

func ProductsPage(w http.ResponseWriter, r *http.Request) {
	pageVariables := PageVariables{
		Title: "Products",
	}

	tmpl, err := template.ParseFiles("./front/Products.html")
	if err != nil {
		log.Println("Error parsing template file")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageVariables)
	if err != nil {
		log.Println("Error executing template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CartPage(w http.ResponseWriter, r *http.Request) {
	pageVariables := PageVariables{
		Title: "Cart",
	}

	tmpl, err := template.ParseFiles("./front/Cart.html")
	if err != nil {
		log.Println("Error executing template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageVariables)
	if err != nil {
		log.Println("Error parsing template file")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func jsFile(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("./front/Script.js")
	if err != nil {
		http.Error(w, "Couldn't read file", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Write(data)
}

func AdminPage(w http.ResponseWriter, r *http.Request) {
	pageVariables := PageVariables{
		Title: "Admin Page",
	}
	tmpl, err := template.ParseFiles("./front/AdminPage.html")
	if err != nil {
		log.Println("Error parsing template file")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageVariables)
	if err != nil {
		log.Println("Error executing template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Main(w http.ResponseWriter, r *http.Request) {
	pageVariables := PageVariables{
		Title: "Main",
	}
	tmpl, err := template.ParseFiles("./front/Main.html")
	if err != nil {
		log.Println("Error parsing template file")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageVariables)
	if err != nil {
		log.Println("Error executing template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	pageVariables := PageVariables{
		Title: "Login Page",
	}

	tmpl, err := template.ParseFiles("./front/LoginPage.html")
	if err != nil {
		log.Println("Error parsing template file")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageVariables)
	if err != nil {
		log.Println("Error executing template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	pageVariables := PageVariables{
		Title: "Registration Page",
	}

	tmpl, err := template.ParseFiles("./front/Registration.html")
	if err != nil {
		log.Println("Error parsing template file")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageVariables)
	if err != nil {
		log.Println("Error executing template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Profile(w http.ResponseWriter, r *http.Request) {
	pageVariables := PageVariables{
		Title: "Profile",
	}

	tmpl, err := template.ParseFiles("./front/profile.html")
	if err != nil {
		log.Println("Error parsing template file")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageVariables)
	if err != nil {
		log.Println("Error executing template")
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

func getAllPosts(w http.ResponseWriter, r *http.Request) {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer file.Close()
	log.SetOutput(file)
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
	collection := client.Database("mydb").Collection("posts")

	//_______________________________posts get actions______________________________________
	var results []*ResponsePosts
	filter := bson.M{}

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
	}

	for cur.Next(context.TODO()) {

		var elem ResponsePosts
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
	//___________________________send success response_________________________________________

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	w.Write(responseJSON)
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer file.Close()
	log.SetOutput(file)
	//______________________________find login and password______________________
	CurrentUser = ResponseLogin{}
	var response ResponseLogin
	if CurrentUser.Login == "" {
		response = ResponseLogin{
			Status: 200,
			Cards:  nil,
		}
	} else {
		response = ResponseLogin{
			Status: 500,
			Cards:  nil,
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
