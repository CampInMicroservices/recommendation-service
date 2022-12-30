package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) Live(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "UP"})
}

func (server *Server) Ready(ctx *gin.Context) {

	// Check connection outside REST API.

	// err := server.store.PingDB()
	// if err != nil {
	// 	ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": "DOWN"})
	// 	ctx.Abort()
	// 	return
	// }

	ctx.JSON(http.StatusOK, gin.H{"status": "UP"})
}
