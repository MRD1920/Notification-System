package rabbitmq

import (
	"encoding/json"
	"log"
	"time"

	"github.com/MRD1920/Notification-System/api/service"
	model "github.com/MRD1920/Notification-System/models"
	"github.com/rabbitmq/amqp091-go"
	"golang.org/x/time/rate"
)

func Worker(queueName string, rateLimit *rate.Limiter, channel *amqp091.Channel) {
	msgs, err := channel.Consume(
		queueName, //queue
		"",        //consumer
		false,     // auto-ack, set to false for manual ack
		false,     //exclusive
		false,     //no-local
		false,     //no-wait
		nil,       //args

	)

	failOnError(err, "Failed to register a consumer")

	for d := range msgs {
		if rateLimit.Allow() {
			if processNotification(d.Body) {
				//d.Ack(false) s a message acknowledgment in RabbitMQ where:
				// Ack tells RabbitMQ that the message was successfully processed
				// The false parameter means "single message acknowledgment" (not multiple)
				// Parameters
				// true = Multiple acknowledgment (acknowledges all outstanding messages)
				// false = Single acknowledgment (acknowledges only this message)
				d.Ack(false)
			} else {
				handleFailure(d.Body)
			}
		}
	}

}

func processNotification(body []byte) bool {
	var unmarshaledNotification model.Notification

	err := json.Unmarshal(body, &unmarshaledNotification)
	if err != nil {
		log.Println("Failed to unmarshal notification: ", err)
		return false
	}

	var user model.User

	user, err = service.GetUserFromDB(unmarshaledNotification.Recipient.String())
	if err != nil {
		log.Println("Failed to get user from DB: ", err)
		return false
	}

	time.Sleep(3 * time.Second)

	channel := user.Preference.GetChannel(unmarshaledNotification.Priority)
	println("notifications_" + channel + "_" + unmarshaledNotification.Priority)

	return true

}

func handleFailure(notification []byte) {
	// TODO: Log the failure to PostgreSQL after 3 retries
	log.Printf("Failed to process notification: %s", string(notification))
}

// Define a rate limiter with X requests per minute
func NewRateLimiter(requestsPerMinute int) *rate.Limiter {
	return rate.NewLimiter(rate.Every((time.Minute)/time.Duration(requestsPerMinute)), 1)
}
