package main

import (
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-mail/mail"
	"github.com/gorilla/mux"
)

var rxEmail = regexp.MustCompile(".+@.+\\..+")

type Message struct {
	Name    string
	Email   string
	Content string
	Errors  map[string]string
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/", send).Methods("POST")
	router.HandleFunc("/confirmation", confirmation).Methods("GET")

	// Choose the folder to serve
	cssDir := "/css/"
	router.PathPrefix(cssDir).Handler(http.StripPrefix(cssDir, http.FileServer(http.Dir("."+cssDir))))

	imgDir := "/img/"
	router.PathPrefix(imgDir).Handler(http.StripPrefix(imgDir, http.FileServer(http.Dir("."+imgDir))))

	jsDir := "/js/"
	router.PathPrefix(jsDir).Handler(http.StripPrefix(jsDir, http.FileServer(http.Dir("."+jsDir))))

	log.Println("Listening...")
	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatal(err)
	}
}

func (msg *Message) Validate() bool {
	msg.Errors = make(map[string]string)

	match := rxEmail.Match([]byte(msg.Email))
	if !match {
		msg.Errors["Email"] = "Please enter a valid email address"
	}

	if strings.TrimSpace(msg.Content) == "" {
		msg.Errors["Content"] = "Please enter a message"
	}

	if strings.TrimSpace(msg.Name) == "" {
		msg.Errors["Name"] = "Please enter your name"
	}

	return len(msg.Errors) == 0
}

func (msg *Message) Deliver() error {
	data := "Name: " + msg.Name + " " + "Content: " + msg.Content + " " + "Email: " + msg.Email
	email := mail.NewMessage()
	email.SetHeader("To", "pawel.grzesik@protonmail.com")
	email.SetHeader("From", "redgrip.firma@gmail.com")
	email.SetHeader("Reply-To", msg.Email)
	email.SetHeader("Subject", "RedGrip.io contact from website")
	email.SetBody("text/plain", data)

	username := "redgrip.firma@gmail.com"
	password := "l1q9CYkghAXZzake"

	return mail.NewDialer("smtp.gmail.com", 587, username, password).DialAndSend(email)
}

func index(w http.ResponseWriter, r *http.Request) {
	render(w, "index.html", nil)
}

func confirmation(w http.ResponseWriter, r *http.Request) {
	render(w, "confirmation.html", nil)
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

func send(w http.ResponseWriter, r *http.Request) {
	// Step 1: Validate form
	msg := &Message{
		Name:    r.PostFormValue("name"),
		Email:   r.PostFormValue("email"),
		Content: r.PostFormValue("content"),
	}

	if !msg.Validate() {
		render(w, "home.html", msg)
		return
	}

	// Step 2: Send message in an email
	log.Println("Name:", msg.Name, "Email:", msg.Email, "Content:", msg.Content)
	if err := msg.Deliver(); err != nil {
		log.Println(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
		return
	}

	// Step 3: Redirect to confirmation page
	http.Redirect(w, r, "/confirmation", http.StatusSeeOther)
}
