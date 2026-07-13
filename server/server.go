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

func SetupRoutes(db *sqlx.DB) *Server {
	mux := http.NewServeMux()
	// rider's auth
	mux.HandleFunc("/v1/riders/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterRider(w, r, db)
	})
	mux.HandleFunc("/v1/riders/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginRider(w, r, db)
	})
	mux.HandleFunc("/v1/riders/delete", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteRider(w, r, db)
	})

	// driver's auth
	mux.HandleFunc("/v1/drivers/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterDriver(w, r, db)
	})
	mux.HandleFunc("/v1/drivers/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginDriver(w, r, db)
	})
	mux.HandleFunc("/v1/drivers/delete", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteDriver(w, r, db)
	})

	mux.HandleFunc("/v1/drivers/status", middleware.Authenticate(middleware.RequireRole("driver",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.UpdateDriverStatus(w, r, db)
		},
	)))

	// driver location update
	mux.HandleFunc("/v1/drivers/location", middleware.Authenticate(middleware.RequireRole("driver",
		func(w http.ResponseWriter, r *http.Request) {
			handlers.UpdateDriverLocation(w, r, db)
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
