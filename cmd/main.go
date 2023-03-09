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
	mcfg := model.DefaultEmailConfig()
	emailSrv := model.NewEmailService(mcfg)

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
	forgotPwd := views.Must(views.ParseFS(templates.FS,
		"forgot_password.gohtml", "tailwind.gohtml",
	))
	resetPwd := views.Must(views.ParseFS(templates.FS,
		"reset_password.gohtml", "tailwind.gohtml",
	))
	cye := views.Must(views.ParseFS(templates.FS,
		"check_your_email.gohtml", "tailwind.gohtml",
	))
	u := controllers.New(signinTmp, signupTmp, forgotPwd, cye, resetPwd,
		userSrv, sessionSrv, emailSrv, &model.PasswordResetService{})
	r.Get("/signin", u.Signin)
	r.Post("/signin", u.ProcessSignIn)
	r.Get("/signup", u.Signup)
	r.Post("/signup", u.ProcessSignup)
	r.Post("/signout", u.ProcessSignout)
	r.Get("/users/me", u.CurrentUser)
	r.Get("/forgot-pw", u.ForgotPassword)
	r.Post("/forgot-pw", u.ProcessForgotPassword)
	r.Get("/reset-pw", u.ResetPassword)
	r.Post("/reset-pw", u.ProcessResetPassword)

	// TODO: auth key should be a config value and secure need to be removed on prod
	CSRF := csrf.Protect([]byte("gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"), csrf.Secure(false))
	http.ListenAndServe(":3000", CSRF(r))
}
