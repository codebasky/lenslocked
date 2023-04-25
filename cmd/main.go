package main

import (
	"fmt"
	"net/http"

	"github.com/codebasky/lenslocked/config"
	"github.com/codebasky/lenslocked/controllers"
	"github.com/codebasky/lenslocked/model"
	"github.com/codebasky/lenslocked/templates"
	"github.com/codebasky/lenslocked/views"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func main() {
	fmt.Println("Starting web development!!")

	//TODO: Need to keep db and other passwords in vault/other secret service
	cfg := config.LoadConfig()
	fmt.Printf("Using the configuration: %+v\n", cfg)
	db, err := model.Open(cfg.DBCfg)
	if err != nil {
		fmt.Printf("Error on db connection: %s", err)
		return
	}

	model.Migrate(db)

	userSrv := model.NewUserSrv(db)
	sessionSrv := model.NewSessionSrv(db)
	emailSrv := model.NewEmailService(cfg.SMTPCfg)
	pwdSrv := model.NewPwdResetService(db)

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
		userSrv, sessionSrv, emailSrv, pwdSrv)
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
	CSRF := csrf.Protect([]byte(cfg.ServerCfg.AuthKey), csrf.Secure(false))
	http.ListenAndServe(cfg.ServerCfg.Address, CSRF(r))
}
