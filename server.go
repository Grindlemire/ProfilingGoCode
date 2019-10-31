package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime/trace"
	"time"

	log "github.com/cihub/seelog"
	"github.com/gorilla/mux"
	"github.com/pkg/profile"
)

// Server ...
type Server struct {
	Profiling bool
	Tracing   bool
	Done      chan struct{}
	Profile   interface {
		Stop()
	}
}

// LaunchServer ...
func LaunchServer() {
	s := Server{
		Profiling: false,
		Done:      make(chan struct{}),
		Profile:   nil,
	}

	r := mux.NewRouter()
	r.Methods("GET").Path("/profiling").HandlerFunc(s.setProfiling)
	r.Methods("GET").Path("/cpuintensive/start").HandlerFunc(s.runAlgorithm)
	r.Methods("GET").Path("/cpuintensive/stop").HandlerFunc(s.stopAlgorithm)
	server := &http.Server{
		Addr:    "localhost:1414",
		Handler: r,
	}

	log.Info("Starting server")
	err := server.ListenAndServe()
	if err != nil {
		log.Error("Error starting server")
		exit(1)
	}
}

func (s *Server) setTracing(w http.ResponseWriter, r *http.Request) {
	if !s.Tracing {
		s.Tracing = true
		f, err := os.Create("./out.trace")
		if err != nil {
			log.Error("Error creating file: ", err)
			exit(1)
		}
		trace.Start(f)
	} else {
		s.Tracing = false
		trace.Stop()
	}
}

func (s *Server) setProfiling(w http.ResponseWriter, r *http.Request) {
	if !s.Profiling {
		s.Profiling = true
		s.Profile = profile.Start(profile.CPUProfile, profile.ProfilePath("./"))
	} else {
		s.Profiling = false
		s.Profile.Stop()
	}

	log.Infof("Profiling is now: %t", s.Profiling)

	profilingStr := fmt.Sprintf("%t", s.Profiling)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(profilingStr))
}

func (s *Server) runAlgorithm(w http.ResponseWriter, r *http.Request) {
	log.Info("running algorithm")

	go func() {
		for {
			executeAlgorithm(opts)
			select {
			case <-s.Done:
				log.Info("Stopping code execution")
				return
			default:
			}
			log.Info("Running algorithm again")
			time.Sleep(500 * time.Millisecond)
		}
	}()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Started Algorithm"))
}

func (s *Server) stopAlgorithm(w http.ResponseWriter, r *http.Request) {
	log.Info("Stopping algorithm")
	s.Done <- struct{}{}
}
