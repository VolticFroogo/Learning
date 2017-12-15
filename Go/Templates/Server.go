package main

import (
    "log"
    "html/template"
    "github.com/gorilla/mux"
    "net/http"
)

type PageVariables struct {
	Name string
    Age int
}

func indexPage(w http.ResponseWriter, r *http.Request) {
    indexPageVariables := PageVariables {
        Name: "Harry",
        Age: 13,
    }

    t, err := template.ParseFiles("Websites/index.html") //Parse the index HTML page
    if err != nil { // If there's an error
        log.Print("Template parsing error: ", err) // Log it
  	}

    err = t.Execute(w, indexPageVariables) // Execute temmplate with variables
    if err != nil { // If there's an error
        log.Print("Template execution error: ", err) // Log it
  	}
}

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/", indexPage)

    log.Printf("Server started listening.")
    http.ListenAndServe(":3737", r)
}