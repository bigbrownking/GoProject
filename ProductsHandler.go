package main

import (
	"net/http"
)

type Products struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
	Desc  string `json:"desc"`
}

func ProductsHandler(w http.ResponseWriter, r *http.Request) {

}
