package messagequeue

import (
	"fmt"

	"github.com/alperklc/the-zula/service/infrastructure/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type MQ struct {
	logger  zerolog.Logger
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Tag     string
	Done    chan error
}

func New(amqpUri string) (*MQ, error) {
	c := &MQ{
		logger:  logger.Get(),
		Conn:    nil,
		Channel: nil,
		Tag:     "service",
		Done:    make(chan error),
	}

	var err error

	c.logger.Info().Msg("amqp: dialing")

	config := amqp.Config{Properties: amqp.NewConnectionProperties(), SASL: nil}
	c.Conn, err = amqp.DialConfig(amqpUri, config)

	if err != nil {
		c.logger.Error().Msgf("amqp: dialing failed %s", err.Error())
		return nil, fmt.Errorf("dial: %s", err)
	}

	c.logger.Info().Msg("amqp: connected, getting Channel")
	c.Channel, err = c.Conn.Channel()
	if err != nil {
		c.logger.Error().Msgf("amqp: could not get channel %s", err.Error())
		return nil, fmt.Errorf("channel: %s", err)
	}

	exchanges := map[string]ExchangeConfig{
		XCHG_MAIN_EXCHANGE: NewExchange(XCHG_MAIN_EXCHANGE, "direct"),
	}

	queues := map[string]QueueConfig{
		Q_NOTIFICATIONS: NewDurableQueue(Q_NOTIFICATIONS),
		Q_ACTIVITY_LOG:  NewDurableQueue(Q_ACTIVITY_LOG),
		Q_REFERENCES:    NewDurableQueue(Q_REFERENCES),
		Q_SCRAPING:      NewDurableQueue(Q_SCRAPING),
		Q_USER_UPDATES:  NewDurableQueue(Q_USER_UPDATES),
	}

	bindings := []BindingConfig{
		NewBinding(XCHG_MAIN_EXCHANGE, RK_ONLY_LOG, Q_ACTIVITY_LOG),

		NewBinding(XCHG_MAIN_EXCHANGE, RK_NOTIFICA, Q_NOTIFICATIONS),
		NewBinding(XCHG_MAIN_EXCHANGE, RK_NOTIFICA, Q_ACTIVITY_LOG),

		NewBinding(XCHG_MAIN_EXCHANGE, RK_NOTI_REF, Q_NOTIFICATIONS),
		NewBinding(XCHG_MAIN_EXCHANGE, RK_NOTI_REF, Q_REFERENCES),
		NewBinding(XCHG_MAIN_EXCHANGE, RK_NOTI_REF, Q_ACTIVITY_LOG),

		NewBinding(XCHG_MAIN_EXCHANGE, RK_NOTI_SCR, Q_NOTIFICATIONS),
		NewBinding(XCHG_MAIN_EXCHANGE, RK_NOTI_SCR, Q_SCRAPING),
		NewBinding(XCHG_MAIN_EXCHANGE, RK_NOTI_SCR, Q_ACTIVITY_LOG),

		NewBinding(XCHG_MAIN_EXCHANGE, RK_SCR_ONLY, Q_SCRAPING),
		NewBinding(XCHG_MAIN_EXCHANGE, RK_SCR_ONLY, Q_ACTIVITY_LOG),

		NewBinding(XCHG_MAIN_EXCHANGE, RK_NOTI_USR, Q_NOTIFICATIONS),
		NewBinding(XCHG_MAIN_EXCHANGE, RK_NOTI_USR, Q_USER_UPDATES),
		NewBinding(XCHG_MAIN_EXCHANGE, RK_NOTI_USR, Q_ACTIVITY_LOG),
	}

	for _, exchange := range exchanges {
		c.logger.Print(fmt.Sprintf("amqp: declaring Exchange (%s)", exchange.Name))
		if err = c.Channel.ExchangeDeclare(exchange.Name, exchange.Type, exchange.Durable, exchange.AutoDelete, exchange.Internal, exchange.NoWait, exchange.Args); err != nil {
			c.logger.Error().Msgf("amqp: could not declare exchange %s", err.Error())
			return nil, fmt.Errorf("exchange declare: %s", err)
		}
	}

	for _, queue := range queues {
		c.logger.Print(fmt.Sprintf("amqp: declaring Queue (%s)", queue.Name))
		if _, err = c.Channel.QueueDeclare(queue.Name, queue.Durable, queue.AutoDelete, queue.Exclusive, queue.NoWait, queue.Args); err != nil {
			c.logger.Error().Msgf("amqp: could not declare queue %s", err.Error())
			return nil, fmt.Errorf("queue declare: %s", err)
		}
	}

	for _, binding := range bindings {
		c.logger.Info().Msgf("amqp: binding to Exchange ('%s' to %s)", binding.SourceExchange, binding.QueueName)
		if err = c.Channel.QueueBind(binding.QueueName, binding.BindingKey, binding.SourceExchange, binding.NoWait, binding.Args); err != nil {
			return nil, fmt.Errorf("queue bind: %s", err)
		}
	}

	c.logger.Info().Msgf("amqp: queue bound to Exchange, starting to consume %s", c.Tag)
	return c, nil
}
