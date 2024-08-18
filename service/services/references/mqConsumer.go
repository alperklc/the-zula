package referencesService

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"

	"github.com/alperklc/the-zula/service/infrastructure/db/notes"
	messagequeue "github.com/alperklc/the-zula/service/infrastructure/messageQueue"
	mqpublisher "github.com/alperklc/the-zula/service/infrastructure/messageQueue/publisher"
)

type MqConsumer interface {
	Start()
}

type dataSources struct {
	logger     zerolog.Logger
	mq         *messagequeue.MQ
	references ReferencesService
}

func NewMqConsumer(l zerolog.Logger, r ReferencesService, mq *messagequeue.MQ) MqConsumer {
	return &dataSources{
		logger:     l,
		references: r,
		mq:         mq,
	}
}

func (ds *dataSources) Start() {
	deliveries, err := ds.mq.Channel.Consume(
		messagequeue.Q_REFERENCES, // name
		ds.mq.Tag+"references",    // consumerTag,
		false,                     // noAck
		false,                     // exclusive
		false,                     // noLocal
		false,                     // noWait
		nil,                       // arguments
	)
	if err != nil {
		ds.logger.Fatal().Msg("referencesService: could not start consuming to the queue")
	}
	forever := make(chan bool)

	go func(references ReferencesService, deliveries <-chan amqp.Delivery, done chan bool) {
		for d := range deliveries {
			incomingMessage := mqpublisher.ActivityMessage{}
			if err := json.Unmarshal(d.Body, &incomingMessage); err != nil {
				ds.logger.Error().Msgf("referencesService: error parsing incoming message: %s", err.Error())
				break
			}

			if incomingMessage.ResourceType != mqpublisher.ResourceTypeNote {
				ds.logger.Debug().Msg("referencesService: message is not relevant, ignoring")
				break
			}

			fmt.Println(incomingMessage)

			if incomingMessage.Action == mqpublisher.ActionCreate {
				note := notes.NoteDocument{}
				errUnmarshal := json.Unmarshal(*incomingMessage.Object, &note)
				if errUnmarshal != nil {
					ds.logger.Error().Msgf("referencesService: error parsing incoming note: %s", errUnmarshal.Error())
					break
				}

				references.UpsertReferencesOfNote(incomingMessage.ObjectID, note.Content)
			} else if incomingMessage.Action == mqpublisher.ActionUpdate {
				var update map[string]interface{}
				errUnmarshal := json.Unmarshal(*incomingMessage.Object, &update)
				if errUnmarshal != nil {
					ds.logger.Error().Msgf("referencesService: error parsing incoming note: %s", errUnmarshal.Error())
					break
				}

				content, contentChanged := update["content"]
				if contentChanged {
					references.UpsertReferencesOfNote(incomingMessage.ObjectID, content.(string))
				}
			} else if incomingMessage.Action == mqpublisher.ActionDelete {
				references.DeleteReferencesOfNote(incomingMessage.ObjectID)
			}

			ds.logger.Info().Msgf("referencesService: received message %s - tag: %d, length: %d", d.Body, d.DeliveryTag, len(d.Body))
		}
		ds.logger.Debug().Msg("referencesService: deliveries channel closed")
		<-done
	}(ds.references, deliveries, forever)
}
