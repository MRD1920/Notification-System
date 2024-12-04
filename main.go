package main

import (
	"fmt"

	"github.com/MRD1920/Notification-System/api"
	"github.com/MRD1920/Notification-System/config"
	DB "github.com/MRD1920/Notification-System/db"
	"github.com/MRD1920/Notification-System/rabbitmq"
)

func main() {
	//Load config from .env file
	config.LoadConfig()

	fmt.Println("Starting the server...")
	//Connect to Db and close the connection when the app exits
	fmt.Println("Connecting to DB...")
	DB.InitDBPool()
	defer DB.CloseDBPool()

	//Connect to RabbitMQ and close the connection when the app exits
	fmt.Println("Connecting to RabbitMQ...")

	api.NewServer()
	fmt.Println("Server is running...")

	// Define worker pools
	pools := []*rabbitmq.WorkerPool{
		rabbitmq.NewWorkerPool("notifications_email_high", rabbitmq.NewRateLimiter(50), 3, 5),
		rabbitmq.NewWorkerPool("notifications_sms_high", rabbitmq.NewRateLimiter(10), 3, 5),
		rabbitmq.NewWorkerPool("notifications_push_high", rabbitmq.NewRateLimiter(100), 3, 5),
		rabbitmq.NewWorkerPool("notifications_email_medium", rabbitmq.NewRateLimiter(25), 2, 3),
		rabbitmq.NewWorkerPool("notifications_sms_medium", rabbitmq.NewRateLimiter(6), 2, 3),
		rabbitmq.NewWorkerPool("notifications_push_medium", rabbitmq.NewRateLimiter(75), 2, 3),
		rabbitmq.NewWorkerPool("notifications_email_low", rabbitmq.NewRateLimiter(10), 1, 2),
		rabbitmq.NewWorkerPool("notifications_sms_low", rabbitmq.NewRateLimiter(3), 1, 2),
		rabbitmq.NewWorkerPool("notifications_push_low", rabbitmq.NewRateLimiter(50), 1, 2),
	}
	conn, ch := rabbitmq.ConnectRabbitMQ(pools)
	defer conn.Close()
	defer ch.Close()

	rabbitmq.Scheduler(pools, ch)

}
