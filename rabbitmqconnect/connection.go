package rabbitmqconnect

import (
	"example/web-service-gin/util"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectMQ() (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672")
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	// defer conn.Close()
	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel")
	// defer ch.Close()
	return conn, ch
}

func CloseMQ(conn *amqp.Connection, channel *amqp.Channel) {
	defer conn.Close()    //rabbit mq close
	defer channel.Close() //rabbit mq channel close
}
