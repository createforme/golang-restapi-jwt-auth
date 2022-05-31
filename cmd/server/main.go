package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	transportHttp "github.com/createforme/simple-api-golang.git/internal/transport/http"
	"github.com/createforme/simple-api-golang.git/internal/utils"
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

	handler := transportHttp.NewHandler()
	handler.SetupRotues()

	utils.LogInfo(fmt.Sprintf("server started on port %s, running on http://%s%s", app.Port, app.Hostname, app.Port))

	if err := http.ListenAndServe(app.Port, handler.Router); err != nil {
		return err
	}

	return nil
}

func main() {
	app := App{
		Name:     "sso.client.example.fossnsbm.org",
		Version:  "1.0.0",
		Hostname: "localhost",
		Port:     ":4000",
	}

	if err := app.Run(); err != nil {
		log.Error(err)
	}
}
