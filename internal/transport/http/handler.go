package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/go-chi/httprate"
)

type Handler struct {
	Router *chi.Mux
}

// NewHandler -  constructor to create and return a new Handler
func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) SetupRotues() {
	h.Router = chi.NewRouter()

	// logs the start and end of each request, along with some useful data about what was requested,
	// what the response status was, and how long it took to return. When standard output is a TTY,
	// Logger will print in color, otherwise it will print in black and white. Logger prints a request ID if one is provided.
	h.Router.Use(middleware.Logger)

	// clean out double slash mistakes from a user's request path.
	// For example, if a user requests /users//1 or //users////1 will both be treated as: /users/1
	h.Router.Use(middleware.CleanPath)

	// automatically route undefined HEAD requests to GET handlers.
	h.Router.Use(middleware.GetHead)

	// recovers from panics, logs the panic (and a backtrace),
	// returns a HTTP 500 (Internal Server Error) status if possible. Recoverer prints a request ID if one is provided.
	h.Router.Use(middleware.Recoverer)

	// Enable httprate request limiter of 100 requests per minute.
	//
	// rate-limiting is bound to the request IP address via the LimitByIP middleware handler.
	//
	// To have a single rate-limiter for all requests, use httprate.LimitAll(..).
	h.Router.Use(httprate.LimitByIP(100, 1*time.Minute))

	h.Router.Route("/api/v1", func(r chi.Router) {
		r.Get("/", h.TestRoute)
	})

}

func (h *Handler) TestRoute(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Test Routes"))
	return
}
