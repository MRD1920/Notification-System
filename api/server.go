package api

import (
	"fmt"
	"os"

	"github.com/MRD1920/Notification-System/api/controllers"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	GinEngine *gin.Engine
	db        *pgxpool.Pool
}

func NewServer() {

	app := &Server{
		GinEngine: gin.Default(),
	}
	app.GinEngine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	app.GinEngine.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	app.GinEngine.GET("/users", controllers.GetUsers)
	app.GinEngine.GET("/users/get/:id", controllers.GetUser)
	app.GinEngine.POST("/users/create", controllers.CreateUser)
	app.GinEngine.DELETE("/users/delete/:id", controllers.DeleteUser)

	app.GinEngine.POST("/notify", controllers.CreateNotification)

	fmt.Println("Server is running on port: ", os.Getenv("PORT"))

	app.GinEngine.Run(":" + os.Getenv("PORT"))
}
