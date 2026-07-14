package server

import (
	"TravelBackend/handlers"
	"TravelBackend/middleware"
	"context"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

type Server struct {
	router http.Handler
	server *http.Server
}

const (
	readTimeout       = 5 * time.Minute
	readHeaderTimeout = 30 * time.Second
	writeTimeout      = 5 * time.Minute
)

func public(h func(w http.ResponseWriter, r *http.Request, db *sqlx.DB), db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r, db)
	}
}

func registerRiderRoutes(mux *http.ServeMux, db *sqlx.DB) {
	mux.HandleFunc("/v1/riders/register", public(handlers.RegisterRider, db))
	mux.HandleFunc("/v1/riders/login", public(handlers.LoginRider, db))
	mux.HandleFunc("/v1/riders/delete", public(handlers.DeleteRider, db))
}

func registerDriverRoutes(mux *http.ServeMux, db *sqlx.DB) {
	mux.HandleFunc("/v1/drivers/register", public(handlers.RegisterDriver, db))
	mux.HandleFunc("/v1/drivers/login", public(handlers.LoginDriver, db))
	mux.HandleFunc("/v1/drivers/delete", public(handlers.DeleteDriver, db))
}

func SetupRoutes(db *sqlx.DB) *Server {
	mux := http.NewServeMux()

	registerRiderRoutes(mux, db)
	registerDriverRoutes(mux, db)

	mux.HandleFunc("/v1/drivers/status", middleware.Authenticate(middleware.RequireRole("driver",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.UpdateDriverStatus(w, r, db)
		},
	)))

	mux.HandleFunc("/v1/drivers/location", middleware.Authenticate(middleware.RequireRole("driver",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.UpdateDriverLocation(w, r, db)
		},
	)))

	mux.HandleFunc("/v1/drivers/rides/pending", middleware.Authenticate(middleware.RequireRole("driver",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.GetPendingRides(w, r, db)
		},
	)))

	mux.HandleFunc("/v1/rides/request", middleware.Authenticate(middleware.RequireRole("rider",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.RequestRide(w, r, db)
		},
	)))

	mux.HandleFunc("/v1/rides/{id}/accept", middleware.Authenticate(middleware.RequireRole("driver",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.AcceptRide(w, r, db)
		},
	)))
	mux.HandleFunc("/v1/rides/{id}/reject", middleware.Authenticate(middleware.RequireRole("driver",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.RejectRide(w, r, db)
		},
	)))

	var handler http.Handler = mux
	handler = middleware.Logging(handler.ServeHTTP)

	return &Server{
		router: handler,
	}
}

func (s *Server) Start(addr string) {
	s.server = &http.Server{
		Addr:              addr,
		Handler:           s.router,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
	}

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
}

// to gracefully shut down the server

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
