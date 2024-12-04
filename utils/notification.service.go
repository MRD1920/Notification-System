package utils

import (
	"encoding/json"
	"fmt"

	"github.com/MRD1920/Notification-System/api/service"
	model "github.com/MRD1920/Notification-System/models"
	"github.com/MRD1920/Notification-System/rabbitmq"
)

func ValidateNotification(notification model.Notification) error {
	if notification.Recipient.String() == "" {
		return fmt.Errorf("recipient is required")
	}

	if notification.Priority != "low" && notification.Priority != "medium" && notification.Priority != "high" {
		return fmt.Errorf("priority should be low, medium or high")
	}
	return nil
}

func SaveNotification(notification model.Notification) error {
	connection, channel := rabbitmq.ConnectRabbitMQ([]*rabbitmq.WorkerPool{})
	defer connection.Close()
	defer channel.Close()

	notification_bytes, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	user, err := service.GetUserFromDB(notification.Recipient.String())
	if err != nil {
		return err
	}

	userChannel := user.Preference.GetChannel(notification.Priority)

	queueName := fmt.Sprintf("notifications_%s_%s", userChannel, notification.Priority)

	err = rabbitmq.PublishNotification(channel, "", queueName, string(notification_bytes))
	if err != nil {
		return err
	}

	return nil
}
