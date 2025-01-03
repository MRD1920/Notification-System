package rabbitmq

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func ConnectRabbitMQ(pools []*WorkerPool) (*amqp.Connection, *amqp.Channel) {
	connString := fmt.Sprintf("amqp://%s:%s@%s:%s/", os.Getenv("RABBITMQ_USERNAME"), os.Getenv("RABBITMQ_PASSWORD"), os.Getenv("RABBITMQ_HOST"), os.Getenv("RABBITMQ_PORT"))

	//Create connection to RabbitMQ instance
	conn, err := amqp.Dial(connString)
	failOnError(err, "Failed to connect to RabbitMQ")

	//Create a channel to communicate with RabbitMQ
	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	for _, pool := range pools {
		_, err := channel.QueueDeclare(
			pool.QueueName, // name
			false,          // durable
			false,          // delete when unused
			false,          // exclusive
			false,          // no-wait
			nil,            // arguments
		)
		failOnError(err, "Failed to declare a queue")
	}
	return conn, channel
}
