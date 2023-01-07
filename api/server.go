package api

import (
	"recommendation-service/config"
	"recommendation-service/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
)

// Server serves HTTP requests for our banking service
type Server struct {
	config      config.Config
	locationsCB *gobreaker.CircuitBreaker
	router      *gin.Engine
}

// NewServer creates a new HTTP server and set up routing
func NewServer(config config.Config, locationsCB *gobreaker.CircuitBreaker) (*Server, error) {

	gin.SetMode(config.GinMode)
	router := gin.Default()
	router.Use(middleware.Logger(config.LogitAddress))

	server := &Server{
		config:      config,
		locationsCB: locationsCB,
	}

	// Setup routing for server
	v1 := router.Group("v1")
	{
		v1.POST("/locations", server.GetPopularLocations)
		v1.GET("/weather", server.GetLocationWeather)
	}

	// Setup health check routes
	health := router.Group("health")
	{
		health.GET("/live", server.Live)
		health.GET("/ready", server.Ready)
	}

	server.router = router
	return server, nil
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
