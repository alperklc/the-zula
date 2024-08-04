package messagequeue

import amqp "github.com/rabbitmq/amqp091-go"

const (
	// exchange
	XCHG_MAIN_EXCHANGE = "main"

	// queue
	Q_ACTIVITY_LOG  = "activity-log"
	Q_NOTIFICATIONS = "notifications"
	Q_REFERENCES    = "references"
	Q_SCRAPING      = "scraping"
	Q_USER_UPDATES  = "user-updates"

	// routing key
	RK_ONLY_LOG = "log"
	RK_NOTIFICA = "notification"
	RK_NOTI_REF = "notification-reference"
	RK_NOTI_SCR = "notification-scraping"
	RK_SCR_ONLY = "scraping"
	RK_NOTI_USR = "notification-userupdate"
)

type ExchangeConfig struct {
	Name       string
	Type       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

type QueueConfig struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

type BindingConfig struct {
	QueueName      string
	BindingKey     string
	SourceExchange string
	NoWait         bool
	Args           amqp.Table
}

func NewExchange(name, xchgType string) ExchangeConfig {
	return ExchangeConfig{
		Name:       name,
		Type:       xchgType,
		Durable:    true,
		AutoDelete: true,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	}
}

func NewDurableQueue(name string) QueueConfig {
	return QueueConfig{
		Name:       name,
		Durable:    true,
		AutoDelete: true,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	}
}

func NewBinding(source, key, queue string) BindingConfig {
	return BindingConfig{
		QueueName:      queue,
		BindingKey:     key,
		SourceExchange: source,
		NoWait:         false,
		Args:           nil,
	}
}
