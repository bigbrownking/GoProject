package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type PageVariables struct {
	Title string
}

type ResponseStatus struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

var newUserChannel = make(chan struct{}, 1)

var limiter = rate.NewLimiter(rate.Limit(10), 5)

var logger = logrus.New()

func main() {
	logger.SetFormatter(&logrus.JSONFormatter{})
	http.HandleFunc("/", rateLimit(HomePage))
	http.HandleFunc("/Script.js", rateLimit(jsFile))
	http.HandleFunc("/loginPage", rateLimit(LoginPage))
	http.HandleFunc("/AdminPage", rateLimit(AdminPage))
	http.HandleFunc("/ProductsPage", rateLimit(ProductsPage))
	http.HandleFunc("/CartPage", rateLimit(CartPage))
	http.HandleFunc("/login", rateLimit(LoginHandler))       //get
	http.HandleFunc("/Products", rateLimit(ProductsHandler)) //hz
	http.HandleFunc("/Cart", rateLimit(CartHandler))         //hz
	http.HandleFunc("/register", rateLimit(RegisterHandler)) //post
	http.HandleFunc("/admin", rateLimit(AdminHandler))       //post
	http.HandleFunc("/admin/all", rateLimit(AdminAll))       //get
	logger.Info("Server listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.WithError(err).Error("Error starting server")
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}

func rateLimit(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
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
		logger.WithError(err).Error("Error parsing template file")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageVariables)
	if err != nil {
		logger.WithError(err).Error("Error executing template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CartPage(w http.ResponseWriter, r *http.Request) {
	pageVariables := PageVariables{
		Title: "Cart",
	}

	tmpl, err := template.ParseFiles("./front/Cart.html")
	if err != nil {
		logger.WithError(err).Error("Error executing template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageVariables)
	if err != nil {
		logger.WithError(err).Error("Error parsing template file")
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
		logger.WithError(err).Error("Error parsing template file")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageVariables)
	if err != nil {
		logger.WithError(err).Error("Error executing template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	pageVariables := PageVariables{
		Title: "Login Page",
	}

	tmpl, err := template.ParseFiles("./front/LoginPage.html")
	if err != nil {
		logger.WithError(err).Error("Error parsing template file")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageVariables)
	if err != nil {
		logger.WithError(err).Error("Error executing template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	pageVariables := PageVariables{
		Title: "Registration Page",
	}

	tmpl, err := template.ParseFiles("./front/Registration.html")
	if err != nil {
		logger.WithError(err).Error("Error parsing template file")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, pageVariables)
	if err != nil {
		logger.WithError(err).Error("Error executing template")
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
