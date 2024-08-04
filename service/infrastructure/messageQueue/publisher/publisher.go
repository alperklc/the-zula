package mqpublisher

import (
	"encoding/json"
	"fmt"

	messagequeue "github.com/alperklc/the-zula/service/infrastructure/messageQueue"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessagePublisher interface {
	Publish(body ActivityMessage) error
}

type dataSources struct {
	channel *amqp.Channel
}

func New(channel *amqp.Channel) MessagePublisher {
	return &dataSources{
		channel: channel,
	}
}

func (data *dataSources) Publish(message ActivityMessage) error {
	body, marshalingErr := json.Marshal(message)
	if marshalingErr != nil {
		return fmt.Errorf("Publish failed: %s", marshalingErr)
	}
	msg := amqp.Publishing{
		Timestamp:       message.Timestamp,
		Headers:         amqp.Table{},
		ContentType:     "application/json",
		ContentEncoding: "",
		Body:            body,
		DeliveryMode:    1, // 1=non-persistent, 2=persistent
		Priority:        0, // 0-9
	}

	key := "*"
	if message.RoutingKey != nil {
		key = *message.RoutingKey
	}

	if err := data.channel.Publish(
		messagequeue.XCHG_MAIN_EXCHANGE, // publish to an exchange
		key,                             // routing to 0 or more queues
		false,                           // mandatory
		false,                           // immediate
		msg,
	); err != nil {
		return fmt.Errorf("Exchange Publish: %s", err)
	}

	return nil
}
