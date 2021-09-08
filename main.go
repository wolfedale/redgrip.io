package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", index).Methods("GET")

	// Choose the folder to serve
	cssDir := "/css/"
	router.PathPrefix(cssDir).Handler(http.StripPrefix(cssDir, http.FileServer(http.Dir("."+cssDir))))

	imgDir := "/img/"
	router.PathPrefix(imgDir).Handler(http.StripPrefix(imgDir, http.FileServer(http.Dir("."+imgDir))))

	jsDir := "/js/"
	router.PathPrefix(jsDir).Handler(http.StripPrefix(jsDir, http.FileServer(http.Dir("."+jsDir))))

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	log.Println("Listening...")
	err := http.ListenAndServe(":8080", loggedRouter)
	if err != nil {
		log.Fatal(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	render(w, "index.html", nil)
}

func render(w http.ResponseWriter, filename string, data interface{}) {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		log.Println(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
	}
}
