package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
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

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1000000) // Grab the multipart form
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	//get the *fileheaders
	files := r.MultipartForm.File["file"]

	for i := range files { // loop through the files one by one
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		// In production use a random file name associated with DB
		out, err := os.Create("tmp/" + files[i].Filename) // In Linux use "/tmp/program"

		defer out.Close()
		if err != nil {
			fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
			return
		}

		_, err = io.Copy(out, file) // file not files[i] !

		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		fmt.Fprintf(w, "Files uploaded successfully: ")
		fmt.Fprintf(w, files[i].Filename+"\n\n")
		/* 	Here you'd upload the files to S3 or something like that.
		Then you could save that in a MySQL DB. */
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/upload", uploadHandler)

	fmt.Println("Server began listening.")
	http.ListenAndServe(":3737", nil)
}
