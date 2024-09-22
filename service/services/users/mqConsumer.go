package usersService

import (
	"encoding/json"

	"github.com/rs/zerolog"

	messagequeue "github.com/alperklc/the-zula/service/infrastructure/messageQueue"
	mqpublisher "github.com/alperklc/the-zula/service/infrastructure/messageQueue/publisher"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MqConsumer interface {
	Start()
}

type dataSources struct {
	logger zerolog.Logger
	users  UsersService
	mq     *messagequeue.MQ
}

func NewMqConsumer(l zerolog.Logger, u UsersService, mq *messagequeue.MQ) MqConsumer {
	return &dataSources{
		logger: l,
		users:  u,
		mq:     mq,
	}
}

func (ds *dataSources) Start() {
	deliveries, err := ds.mq.Channel.Consume(
		messagequeue.Q_USER_UPDATES, // name
		ds.mq.Tag+"user_updates",    // consumerTag,
		false,                       // noAck
		false,                       // exclusive
		false,                       // noLocal
		false,                       // noWait
		nil,                         // arguments
	)
	if err != nil {
		ds.logger.Fatal().Msg("usersService: could not start consuming to the queue")
	}

	forever := make(chan bool)

	go func(users UsersService, deliveries <-chan amqp.Delivery, done chan bool) {
		for d := range deliveries {
			incomingMessage := mqpublisher.ActivityMessage{}
			if err := json.Unmarshal(d.Body, &incomingMessage); err != nil {
				ds.logger.Error().Msgf("usersService: error parsing incoming message: %s", err.Error())
				break
			}

			if err := users.RefreshUserInCache(incomingMessage.ObjectID); err != nil {
				ds.logger.Error().Msgf("usersService: error invalidating the user in cache: %s", err.Error())
				break
			}

			ds.logger.Info().Msgf("usersService: received message %s - tag: %d, length: %d", d.Body, d.DeliveryTag, len(d.Body))
			// todo: ensure the message received & processed by all instances of the the zula backend
			d.Ack(false)
		}
		ds.logger.Debug().Msg("usersService: deliveries channel closed")
		<-done
	}(ds.users, deliveries, forever)
}
