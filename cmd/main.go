package main

import (
	"fmt"
	"net/http"

	"github.com/codebasky/lenslocked/controllers"
	"github.com/codebasky/lenslocked/views"
)

func main() {
	fmt.Println("Starting web development")
	http.HandleFunc("/", controllers.StaticHandler(views.Must(views.Parse("home.gohtml"))))
	http.HandleFunc("/contact", controllers.StaticHandler(views.Must(views.Parse("contact.gohtml"))))
	http.HandleFunc("/faq", controllers.StaticHandler(views.Must(views.Parse("faq.gohtml"))))
	http.ListenAndServe(":3000", nil)
}
