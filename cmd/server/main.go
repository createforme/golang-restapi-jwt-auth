package main

import (
	"fmt"
	"net/http"

	"github.com/createforme/golang-restapi-jwt-auth/internal/database"
	transportHttp "github.com/createforme/golang-restapi-jwt-auth/internal/transport/http"
	"github.com/createforme/golang-restapi-jwt-auth/internal/user"
	log "github.com/sirupsen/logrus"

	"github.com/createforme/golang-restapi-jwt-auth/internal/utils"
)

type App struct {
	Name     string
	Version  string
	Hostname string
	Port     string
}

func (app *App) Run() error {
	log.WithFields(
		log.Fields{
			"AppName":    app.Name,
			"AppVersion": app.Version,
			"HostName":   app.Hostname,
		}).Info("Setting up Application")

	db, err := database.NewDatabase()
	if err != nil {
		return err
	}

	userService := user.NewService(db)

	handler := transportHttp.NewHandler(userService)
	handler.SetupRotues()

	utils.LogInfo(fmt.Sprintf("server started on port %s, running on http://%s%s", app.Port, app.Hostname, app.Port))

	if err := http.ListenAndServe(app.Port, handler.Router); err != nil {
		return err
	}

	return nil
}

func main() {
	app := App{
		Name:     "app",
		Version:  "1.0.0",
		Hostname: "localhost",
		Port:     ":4000",
	}

	if err := app.Run(); err != nil {
		log.Error(err)
	}
}
