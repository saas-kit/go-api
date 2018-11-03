package http

import (
	"net/http"

	"github.com/xujiajun/gorouter"
)

// NewServer func returns a new instance of the Server structure
func NewServer(cnf serverConfig) *Server {
	mux := gorouter.New()
	return &Server{router: mux, config: cnf}
}

// ServerConfig interface
type serverConfig interface {
	ListenAddr() string
	APIVersion() string
}

// Server structure
type Server struct {
	config serverConfig
	router *gorouter.Router
}

// ListenAndServe function to handle requests on incoming connections via http
func (s *Server) ListenAndServe() error {
	s.router.GET("/health", s.healthcheckHandler)
	apiGroup := s.router.Group(s.config.APIVersion())
	setupRoutes(apiGroup)
	return http.ListenAndServe(s.config.ListenAddr(), s.router)
}

func (s *Server) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
