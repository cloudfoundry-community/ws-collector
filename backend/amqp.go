package backend

import (
	"github.com/mchmarny/ws-collector/common"
	"github.com/streadway/amqp"
	"log"
)

func NewAMQPBackend() *AMQPBackend {
	return &AMQPBackend{}
}

type AMQPBackend struct {
	uri          string
	exchange     string
	exchangeType string
}

func (q *AMQPBackend) Config(uri string, args map[string]string) {
	//TODO: Implements
	q.uri = uri
	q.exchange = args["exchange"]
	q.exchangeType = args["exchange_type"]
}

func (q *AMQPBackend) Start(in <-chan *common.Message) {

	conn, err := amqp.Dial(q.uri)
	if err != nil {
		log.Panicf("error on dial: %v", err)
	}
	defer conn.Close()

	chn, err := conn.Channel()
	if err != nil {
		log.Panicf("channel error: %v", err)
	}

	if err := chn.ExchangeDeclare(
		q.exchange,     // name
		q.exchangeType, // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // noWait
		nil,            // arguments
	); err != nil {
		log.Panicf("error on exchange declare: %v", err)
	}

	for {
		select {
		case msg := <-in:
			log.Printf("publishing: %v", msg)
			if err := chn.Publish(
				q.exchange,     // publish to an exchange
				q.exchangeType, // routing to 0 or more queues
				true,           // mandatory
				false,          // immediate
				amqp.Publishing{
					Headers:         amqp.Table{},
					ContentType:     "text/plain",
					ContentEncoding: "",
					Body:            msg.ToBytes(),
					DeliveryMode:    amqp.Transient,
					Priority:        0, // 0-9
				},
			); err != nil {
				log.Fatalf("error on pub: %s", err)
			}
		default:
			//nothing to do here.
		}
	}
}
