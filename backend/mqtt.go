package backend

import (
	"github.com/mchmarny/ws-collector/common"
	"log"

	mqtt "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

func NewMQTTBackend() *MQTTBackend {
	return &MQTTBackend{}
}

type MQTTBackend struct {
	id       string
	uri      string
	topic    string
	username string
	password string
	qos      int
}

func (q *MQTTBackend) Config(uri string, args map[string]string) {
	q.id = args["id"]
	q.uri = args["uri"]
	q.topic = args["topic"]
	q.username = args["username"]
	q.password = args["password"]
	q.qos = common.ParseInt(args["qos"], 0)
}

func (q *MQTTBackend) Start(in <-chan *common.Message) {

	// TODO: Add support for SSL
	opts := mqtt.NewClientOptions().
		AddBroker(q.uri).
		SetClientId(q.id).
		SetCleanSession(true)

	if len(q.username) > 0 && len(q.password) > 0 {
		opts.SetUsername(q.username)
		opts.SetPassword(q.password)
	}

	client := mqtt.NewClient(opts)
	_, err := client.Start()
	if err != nil {
		log.Panicln(err)
	}
	defer client.Disconnect(250)

	for {
		select {
		case msg := <-in:
			receipt := client.Publish(mqtt.QoS(q.qos), q.topic, msg.ToBytes())
			<-receipt
		default:
			//nothing to do here.
		}
	}

}
