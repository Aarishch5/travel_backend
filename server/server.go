package server

import (
	"TravelBackend/handlers"
	"context"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

type Server struct {
	router *http.ServeMux
	server *http.Server
}

const (
	readTimeout       = 5 * time.Minute
	readHeaderTimeout = 30 * time.Second
	writeTimeout      = 5 * time.Minute
)

func SetupRoutes(db *sqlx.DB) *Server {
	mux := http.NewServeMux()
	// driver's apis
	mux.HandleFunc("/create-driver", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateDriver(w, r, db)
	})
	mux.HandleFunc("/delete-driver", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteDriver(w, r, db)
	})

	// rider's apis
	mux.HandleFunc("/create-rider", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateRider(w, r, db)
	})
	mux.HandleFunc("/delete-rider", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteRider(w, r, db)
	})

	return &Server{
		router: mux,
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

// to shut down the server

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
