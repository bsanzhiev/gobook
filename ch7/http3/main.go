package main

import (
	"fmt"
	"log"
	"net/http"
)

type dollars float64

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

func main() {
	db := database{"shoes": 50, "socks": 5}
	// mux := http.NewServeMux()
	// mux.Handle("/list", http.HandlerFunc(db.list))
	// mux.Handle("/price", http.HandlerFunc(db.price))
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	log.Fatal(http.ListenAndServe("localhost:8008", nil))
}

type database map[string]dollars

func (db database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	// price, ok := db[item]
	if price, ok := db[item]; ok {
		fmt.Fprintf(w, "%s\n", price)

	} else {
		w.WriteHeader(http.StatusNotFound) //404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}
