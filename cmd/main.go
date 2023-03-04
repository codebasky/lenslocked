package main

import (
	"fmt"
	"net/http"

	"github.com/codebasky/lenslocked/controllers"
	"github.com/codebasky/lenslocked/model"
	"github.com/codebasky/lenslocked/templates"
	"github.com/codebasky/lenslocked/views"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func main() {
	fmt.Println("Starting web development")

	cfg := model.DefaultPostgresConfig()
	db, err := model.Open(cfg)
	if err != nil {
		fmt.Printf("Error on db connection: %s", err)
		return
	}

	model.Migrate(db)

	userSrv := model.NewUserSrv(db)
	sessionSrv := model.NewSessionSrv(db)

	r := chi.NewRouter()

	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))
	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))
	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))

	signinTmp := views.Must(views.ParseFS(templates.FS,
		"signin.gohtml", "tailwind.gohtml"))
	signupTmp := views.Must(views.ParseFS(templates.FS,
		"signup.gohtml", "tailwind.gohtml"))
	u := controllers.New(signinTmp, signupTmp, userSrv, sessionSrv)
	r.Get("/signin", u.Signin())
	r.Post("/signin", u.ProcessSignIn())
	r.Get("/signup", u.Signup())
	r.Post("/signup", u.ProcessSignup())
	r.Post("/signout", u.ProcessSignout())

	// TODO: auth key should be a config value and secure need to be removed on prod
	CSRF := csrf.Protect([]byte("gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"), csrf.Secure(false))
	http.ListenAndServe(":3000", CSRF(r))
}
