package http

import (
	"encoding/json"
	"net/http"
	"time"

	user "github.com/createforme/golang-restapi-jwt-auth/internal/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/go-chi/httprate"
	log "github.com/sirupsen/logrus"
)

// Handler - store pointer to our comment service
type Handler struct {
	Router      *chi.Mux
	ServiceUser *user.Service
}

// Response - an object to store responses from our api
type Response struct {
	Message string
	Error   string
}

// NewHandler - return a pointer to a handler
func NewHandler(userservice *user.Service) *Handler {
	return &Handler{
		ServiceUser: userservice,
	}
}

// LoggingMiddleware - a handy middleware function that logs out incoming requests
func LogginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(
			log.Fields{
				"Method": r.Method,
				"Path":   r.URL.Path,
				"Host":   r.RemoteAddr,
			}).
			Info("handled request")
		next.ServeHTTP(w, r)
	})
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

	// RedirectSlashes is a middleware that will match request paths with a trailing slash
	// and redirect to the same path, less the trailing slash.
	h.Router.Use(middleware.RedirectSlashes)

	// automatically route undefined HEAD requests to GET handlers.
	h.Router.Use(middleware.GetHead)

	// Throttle is a middleware that limits number of currently processed requests at a time
	// across all users. Note: Throttle is not a rate-limiter per user, instead it just puts a
	// ceiling on the number of currentl in-flight requests being processed from the point
	// from where the Throttle middleware is mounted.
	h.Router.Use(middleware.Throttle(15))

	// ThrottleBacklog is a middleware that limits number of currently processed requests
	// at a time and provides a backlog for holding a finite number of pending requests
	h.Router.Use(middleware.ThrottleBacklog(10, 50, time.Second*10))

	// timeout middleware
	h.Router.Use(middleware.Timeout(time.Second * 60))

	// recovers from panics, logs the panic (and a backtrace),
	// returns a HTTP 500 (Internal Server Error) status if possible. Recoverer prints a request ID if one is provided.
	h.Router.Use(middleware.Recoverer)

	// RealIP is a middleware that sets a http.Request's RemoteAddr to the results of parsing either
	// the X-Real-IP header or the X-Forwarded-For header (in that order).
	h.Router.Use(middleware.RealIP)

	// Enable httprate request limiter of 100 requests per minute.
	//
	// rate-limiting is bound to the request IP address via the LimitByIP middleware handler.
	//
	// To have a single rate-limiter for all requests, use httprate.LimitAll(..).
	h.Router.Use(httprate.LimitByIP(100, 1*time.Minute))

	h.Router.Route("/api/v2", func(r chi.Router) {

		r.Use(cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: []string{"https://*", "http://*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", h.CreateUser)
			r.Post("/login", h.AuthUser)
		})

		r.Route("/user", func(r chi.Router) {
			r.Use(AuthMiddleware)
			r.Get("/me", h.CurrentUser)
			r.Get("/{username}", h.GetUser)
		})

		/* handle errors */

		h.Router.NotFound(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": "route not found"})
		})

		h.Router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]interface{}{"error": "method is not valid"})
		})
	})
}

// handle ok responses
func sendOkResponse(w http.ResponseWriter, resp interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp)
}

// handle error responses
func sendErrorResponse(w http.ResponseWriter, message string, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)

	if err := json.NewEncoder(w).Encode(Response{
		Message: message,
		Error:   err.Error(),
	}); err != nil {
		panic(err)
	}
}
