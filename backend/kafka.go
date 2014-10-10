package backend

import (
	"github.com/Shopify/sarama"
	"github.com/mchmarny/ws-collector/common"
	"log"
	"time"
)

func NewKafkaBackend() *KafkaBackend {
	return &KafkaBackend{}
}

type KafkaBackend struct {
	id              string
	uri             string
	topic           string
	partition       int
	bufferTimeInSec time.Duration
	timeoutInSec    time.Duration
}

func (q *KafkaBackend) Config(uri string, args map[string]string) {
	q.uri = uri
	q.topic = args["topic"]
	q.partition = common.ParseInt(args["partition"], 0)
	q.bufferTimeInSec = time.Duration(
		common.ParseInt(args["buffer_time_in_sec"], 3)) * time.Second
	q.timeoutInSec = time.Duration(
		common.ParseInt(args["timeout_in_sec"], 3)) * time.Second
}

func (q *KafkaBackend) Start(in <-chan *common.Message) {

	pConfig := sarama.NewProducerConfig()
	pConfig.MaxBufferTime = q.bufferTimeInSec * time.Second

	producerClient, err := sarama.NewClient(q.id, []string{q.uri},
		sarama.NewClientConfig())
	if err != nil {
		log.Panicln(err)
	}

	defer producerClient.Close()

	producer, err := sarama.NewProducer(producerClient, pConfig)
	if err != nil {
		log.Panicln(err)
	}

	defer producer.Close()

	for {
		select {
		case msg := <-in:
			err := producer.SendMessage(q.topic, nil,
				sarama.StringEncoder(msg.ToBytes()))
			if err != nil {
				log.Fatalln(err)
			}
		default:
			//nothing to do here.
		}
	}

}
