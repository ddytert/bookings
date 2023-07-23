package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/ddytert/bookings/pkg/config"
	"github.com/ddytert/bookings/pkg/handlers"
	"github.com/ddytert/bookings/pkg/render"
	"log"
	"net/http"
	"time"
)

const portNUmber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreteTemplateCache()
	if err != nil {
		log.Fatal(err)
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	fmt.Println("Starting Webserver on port", portNUmber)

	srv := &http.Server{
		Addr:    portNUmber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}
