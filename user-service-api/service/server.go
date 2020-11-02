package service

import (
	"github.com/codegangsta/negroni"

	"context"
	"net/http"
	"time"

	serviceDomain "keon.com/CitadelAllianceLobbyServer/user-service-api/service/domain"
)

// Request JSON body limit is set at 5MB (currently not enforced)
const BodyLimitBytes uint32 = 1048576 * 5

// Server type
type Server struct {
	negroni    *negroni.Negroni
	router     *Router
	httpServer *http.Server
	timeout    time.Duration
}

// Config type
type Config struct {
}

// ServerOptions for running the server
type ServerOptions struct {
	Timeout         time.Duration
	ShutdownHandler func()
	CertPath        string
	KeyPath         string
}

// NewServer Returns a new Server object
func NewServer(options *Config) *Server {

	// set up server and middlewares
	n := negroni.Classic()

	s := &Server{n, nil, nil, 0}

	return s
}

func (s *Server) UseMiddleware(middleware serviceDomain.IMiddleware) *Server {
	// next convert it into negroni style handlerfunc
	s.negroni.Use(negroni.HandlerFunc(middleware.Handler))
	return s
}

func (s *Server) UseRouter(router *Router) *Server {
	// add router and clear mux.context values at the end of request life-times
	s.negroni.UseHandler(router)
	return s
}

func (s *Server) Run(address string, options ServerOptions) *Server {
	s.timeout = options.Timeout
	s.httpServer = &http.Server{
		IdleTimeout: options.Timeout,
		Addr:        address,
		Handler:     s.negroni,
	}

	s.httpServer.RegisterOnShutdown(options.ShutdownHandler)

	err := s.httpServer.ListenAndServeTLS(options.CertPath, options.KeyPath)
	if err != nil {
		panic(err)
	}
	return s
}

func (s *Server) Stop() {
	s.httpServer.Shutdown(context.Background())
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) *Server {
	s.negroni.ServeHTTP(w, r)
	return s
}
