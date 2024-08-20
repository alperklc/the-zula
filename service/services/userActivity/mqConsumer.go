package userActivityService

import (
	"encoding/json"

	useractivity "github.com/alperklc/the-zula/service/infrastructure/db/userActivity"
	"github.com/rs/zerolog"

	messagequeue "github.com/alperklc/the-zula/service/infrastructure/messageQueue"
	mqpublisher "github.com/alperklc/the-zula/service/infrastructure/messageQueue/publisher"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MqConsumer interface {
	Start()
}

type dataSources struct {
	logger                 zerolog.Logger
	useractivityCollection useractivity.Collection
	mq                     *messagequeue.MQ
}

func NewMqConsumer(l zerolog.Logger, ua useractivity.Collection, mq *messagequeue.MQ) MqConsumer {
	return &dataSources{
		logger:                 l,
		useractivityCollection: ua,
		mq:                     mq,
	}
}

func (ds *dataSources) Start() {
	deliveries, err := ds.mq.Channel.Consume(
		messagequeue.Q_ACTIVITY_LOG, // name
		ds.mq.Tag+"activity_log",    // consumerTag,
		false,                       // noAck
		false,                       // exclusive
		false,                       // noLocal
		false,                       // noWait
		nil,                         // arguments
	)
	if err != nil {
		ds.logger.Fatal().Msg("userActivityService: could not start consuming to the queue")
	}

	forever := make(chan bool)

	go func(dbCollection useractivity.Collection, deliveries <-chan amqp.Delivery, done chan bool) {
		for d := range deliveries {
			incomingMessage := mqpublisher.ActivityMessage{}
			if err := json.Unmarshal(d.Body, &incomingMessage); err != nil {
				ds.logger.Error().Msgf("userActivityService: error parsing incoming message: ", err.Error())
				break
			}

			if _, dbWriteErr := dbCollection.InsertOne(incomingMessage.UserID, incomingMessage.ResourceType, incomingMessage.Action, incomingMessage.ObjectID); dbWriteErr != nil {
				ds.logger.Error().Msgf("userActivityService: error writing incoming message into the database: ", dbWriteErr.Error())
				break
			}

			ds.logger.Info().Msgf("userActivityService: received message %s - tag: %s, length: %d", d.Body, d.DeliveryTag, len(d.Body))
			d.Ack(false)
		}
		ds.logger.Debug().Msg("userActivityService: deliveries channel closed")
		<-done
	}(ds.useractivityCollection, deliveries, forever)
}
