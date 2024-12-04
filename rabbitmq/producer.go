package rabbitmq

import "github.com/rabbitmq/amqp091-go"

func PublishNotification(channel *amqp091.Channel, exchange string, queueName string, body string) error {
	q, err := channel.QueueDeclare(
		queueName, //name
		false,     //durable
		false,     //autoDelete
		false,     //exclusive
		false,     //nowait
		nil,       //args

	)
	if err != nil {
		return err
	}

	err = channel.Publish(
		exchange, //exchange
		q.Name,   // **queue name is routing key**
		false,    //mandatory
		false,    //immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})

	if err != nil {
		return err
	}

	return nil
}
