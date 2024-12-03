package main

import (
	"fmt"

	"github.com/MRD1920/Notification-System/api"
	"github.com/MRD1920/Notification-System/config"
	DB "github.com/MRD1920/Notification-System/db"
)

func main() {
	//Load config from .env file
	config.LoadConfig()

	fmt.Println("Starting the server...")
	//Connect to Db and close the connection when the app exits
	fmt.Println("Connecting to DB...")
	DB.InitDBPool()
	defer DB.CloseDBPool()

	api.NewServer()
	fmt.Println("Server is running...")

}
