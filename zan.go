package zan

import (
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	// Version is ...
	Version = "v0.0.1"
)

type (
	// Server is struct ...
	Server struct {
		http.Server
		contextPool sync.Pool
		routes      map[string]HandlerFunc
	}
)

// NewServer will create a Server instance and response with a pointer which point to is
func NewServer() *Server {
	s := &Server{contextPool: sync.Pool{}}
	s.contextPool.New = func() interface{} {
		c := Context{}
		return &c
	}
	s.routes = map[string]HandlerFunc{}

	return s
}

// Run server
func (s *Server) Run(addr string) error {
	log.Println("start zan", addr)
	s.Addr = addr
	s.Handler = s
	return s.ListenAndServe()
}

// RunTLS Run server with tls
func (s *Server) RunTLS(addr string, certFile string, keyFile string) error {
	s.Addr = addr
	s.Handler = s
	return s.ListenAndServeTLS(certFile, keyFile)
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	c := s.contextPool.Get().(*Context)
	c.req = req
	c.rw = rw

	s.serveHTTPRequest(c)
}

func (s *Server) serveHTTPRequest(c *Context) {
	handler, ok := s.routes[c.req.URL.RequestURI()+"::"+strings.ToUpper(c.req.Method)]
	if !ok {
		c.rw.WriteHeader(http.StatusNotFound)
		return
	}

	handler(c)
}
