package http

import (
	"encoding/json"
	"net/http"

	user "github.com/createforme/golang-restapi-jwt-auth/internal/user"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Handler - store pointer to our comment service
type Handler struct {
	Router      *mux.Router
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

// SetupRoutes - sets up all the routes for our application
func (h *Handler) SetupRotues() {
	log.Info("Setting up routes")

	// initicate new gorilla mox router
	h.Router = mux.NewRouter()
	h.Router.Use(LogginMiddleware)
	h.Router.Use(CORSMiddleware)

	//  authenticated routes
	authRoutes := h.Router.Methods(http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodOptions).Subrouter()

	// Services realted to user
	//authRoutes.HandleFunc("/api/v1/user/create", h.CreateUser).Methods(http.MethodPost)
	authRoutes.HandleFunc("/api/v1/user/me", h.CurrentUser).Methods(http.MethodGet, http.MethodOptions)

	authRoutes.Use(AuthMiddleware)

	// just made this rouer unauth just for local testing.
	h.Router.HandleFunc("/api/v1/user/create", h.CreateUser).Methods(http.MethodPost)

	// users
	h.Router.HandleFunc("/api/v1/user/{username}", h.GetUser).Methods("GET")
	h.Router.HandleFunc("/api/v1/user/auth", h.AuthUser).Methods(http.MethodPost)

	h.Router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(Response{
			Message: "Api is Running OK",
		}); err != nil {
			log.Fatal(err)
			panic(err)
		}
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
