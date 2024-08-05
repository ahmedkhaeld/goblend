package HTTP

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

type Module interface {
	Register(*gin.Engine)
}

type Server struct {
	engine    *gin.Engine
	modules   []Module
	rateLimit ratelimit.Limiter
}

func NewServer(rps int) *Server {
	return &Server{
		engine:    gin.Default(),
		rateLimit: ratelimit.New(rps),
	}
}

func (s *Server) RegisterModule(module Module) {
	s.modules = append(s.modules, module)
}

func (s *Server) SetupAndRun(address string) error {

	// Add global middlewares
	s.engine.Use(requestid.New())
	s.engine.Use(s.leakBucket())

	// Register all modules
	for _, module := range s.modules {
		module.Register(s.engine)
	}

	// Run the server
	return s.engine.Run(address)
}

func (s *Server) leakBucket() gin.HandlerFunc {
	return func(c *gin.Context) {
		s.rateLimit.Take()
		c.Next()
	}
}
