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

func protected(role string, h func(w http.ResponseWriter, r *http.Request, db *sqlx.DB), db *sqlx.DB) http.HandlerFunc {
	return middleware.Authenticate(middleware.RequireRole(role,
		func(w http.ResponseWriter, r *http.Request) {
			h(w, r, db)
		},
	))
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

	mux.HandleFunc("/v1/drivers/status", protected("driver", handlers.UpdateDriverStatus, db))
	//mux.HandleFunc("/v1/drivers/location", protected("driver", handlers.UpdateDriverLocation, db))
	//mux.HandleFunc("/v1/drivers/location", protected("driver", handlers.GetDriverLocation, db))
	mux.HandleFunc("/v1/drivers/location", protected("driver", handlers.DriverLocationHandler, db))
	mux.HandleFunc("/v1/drivers/rides/pending", protected("driver", handlers.GetPendingRides, db))
	mux.HandleFunc("/v1/drivers/rides", protected("driver", handlers.GetAllRides, db))

}

func registerRideRoutes(mux *http.ServeMux, db *sqlx.DB) {
	mux.HandleFunc("/v1/rides/request", protected("rider", handlers.RequestRide, db))
	mux.HandleFunc("/v1/rides/{id}/accept", protected("driver", handlers.AcceptRide, db))
	mux.HandleFunc("/v1/rides/{id}/reject", protected("driver", handlers.RejectRide, db))
	mux.HandleFunc("/v1/rides/{id}/complete", protected("driver", handlers.RideCompleted, db))
}

func SetupRoutes(db *sqlx.DB) *Server {
	mux := http.NewServeMux()

	registerRiderRoutes(mux, db)
	registerDriverRoutes(mux, db)
	registerRideRoutes(mux, db)

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
