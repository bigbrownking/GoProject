package main

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

func CartHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Handle the case where we're getting the cart
		// Connect to your MongoDB client
		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://myAtlasDBUser:111@myatlasclusteredu.z25a02h.mongodb.net/?retryWrites=true&w=majority&appName=myAtlasClusterEDU"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = client.Connect(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer client.Disconnect(ctx)

		// Get the 'carts' collection
		cartsCollection := client.Database("mydb").Collection("carts")

		// Get the cart from the database
		var cart Cart
		err = cartsCollection.FindOne(ctx, bson.M{"userId": CurrentUser.Id}).Decode(&cart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(cart)
	case http.MethodPost:
		// Handle the case where we're adding an item to the cart
		var newItem CartItem
		err := json.NewDecoder(r.Body).Decode(&newItem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Connect to your MongoDB client
		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://myAtlasDBUser:111@myatlasclusteredu.z25a02h.mongodb.net/?retryWrites=true&w=majority&appName=myAtlasClusterEDU"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = client.Connect(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer client.Disconnect(ctx)

		// Get the 'carts' collection
		cartsCollection := client.Database("mydb").Collection("carts")

		// Add the new item to the user's cart in the database
		_, err = cartsCollection.UpdateOne(
			ctx,
			bson.M{"userId": CurrentUser.Id},
			bson.M{"$push": bson.M{"items": newItem}},
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	case http.MethodDelete:
		// Handle the case where we're removing an item from the cart
		var itemToRemove CartItem
		err := json.NewDecoder(r.Body).Decode(&itemToRemove)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Connect to your MongoDB client
		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://myAtlasDBUser:111@myatlasclusteredu.z25a02h.mongodb.net/?retryWrites=true&w=majority&appName=myAtlasClusterEDU"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = client.Connect(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer client.Disconnect(ctx)

		// Get the 'carts' collection
		cartsCollection := client.Database("mydb").Collection("carts")

		// Remove the item from the user's cart in the database
		_, err = cartsCollection.UpdateOne(
			ctx,
			bson.M{"userId": CurrentUser.Id},
			bson.M{"$pull": bson.M{"items": bson.M{"productId": itemToRemove.ProductID}}},
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	default:
		// If we don't recognize the method, return an error
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
