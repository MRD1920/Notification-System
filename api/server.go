package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	GinEngine *gin.Engine
	db        *pgxpool.Pool
}

func NewServer(db *pgxpool.Pool) *Server {
	ginServer := gin.Default()
	server := &Server{
		GinEngine: ginServer,
		db:        db,
	}
	return server
}
func (s *Server) StartServer(port string) {

	s.GinEngine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	s.GinEngine.GET("/health", func(ctx *gin.Context) {
		//I want to just send the response 200
		ctx.JSON(http.StatusOK, gin.H{})
	})
	s.GinEngine.Run()

}
