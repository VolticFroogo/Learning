// Client/Server AJAX JSON Communication using golang web-server and JQuery
// Visit: http://127.0.0.1:8080
package main

import (
	"log"
	"encoding/json"
	"html/template"
	"net/http"
)

type Data struct {
	Response, From, Quote string
}

// Default Request Handler
func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		log.Print("Template parsing error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		log.Print("Template execution error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AJAX Request Handler
func ajaxHandler(w http.ResponseWriter, r *http.Request) {
	var d Data // Create struct to store data
	err := json.NewDecoder(r.Body).Decode(&d) // Decode response to struct
	if err != nil {
		log.Print("JSON decoding error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	n := Data {
		Response: "Success.",
		From: d.From,
		Quote: d.Quote,
	}
	a, err := json.Marshal(n)
	if err != nil {
		log.Print("JSON responding error: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(a) // Write to response
}

func main() {
	http.HandleFunc("/", indexPageHandler)
	http.HandleFunc("/ajax", ajaxHandler)

	log.Println("Server started listening.")

	http.ListenAndServe(":3737", nil)
}