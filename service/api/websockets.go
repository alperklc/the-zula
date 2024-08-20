package api

import (
	"encoding/json"

	messagequeue "github.com/alperklc/the-zula/service/infrastructure/messageQueue"
	mqpublisher "github.com/alperklc/the-zula/service/infrastructure/messageQueue/publisher"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type Notifier interface {
	SendNotification()
}

type resources struct {
	deliveries <-chan amqp.Delivery
	logger     zerolog.Logger
	clientHub  Hub
}

func NewNotifier(l zerolog.Logger, m *messagequeue.MQ, h Hub) Notifier {
	d, err := m.Channel.Consume(
		messagequeue.Q_NOTIFICATIONS, // name
		m.Tag+"notifications",        // consumerTag,
		false,                        // noAck
		false,                        // exclusive
		false,                        // noLocal
		false,                        // noWait
		nil,                          // arguments
	)
	if err != nil {
		l.Fatal().Msg("notifications: could not start consuming to the queue")
	}

	return &resources{
		deliveries: d,
		logger:     l,
		clientHub:  h,
	}
}

func (r *resources) SendNotification() {
	forever := make(chan bool)

	go func(deliveries <-chan amqp.Delivery, done chan bool) {

		for d := range deliveries {
			incomingMessage := mqpublisher.ActivityMessage{}
			if err := json.Unmarshal(d.Body, &incomingMessage); err != nil {
				r.logger.Error().Msgf("notifications: error parsing incoming message: %s", err.Error())
				break
			}
			r.logger.Debug().Msgf("notifications: received message %s - tag: %d, length: %d", d.Body, d.DeliveryTag, len(d.Body))
			EmitToSpecificClient(&r.clientHub, SocketEventStruct{EventName: "msg", EventPayload: incomingMessage}, incomingMessage.ClientID)

			d.Ack(false)
		}
		r.logger.Debug().Msg("notifications: deliveries channel closed")
		<-done
	}(r.deliveries, forever)
}
