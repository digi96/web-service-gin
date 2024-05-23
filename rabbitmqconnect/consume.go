package rabbitmqconnect

import (
	"example/web-service-gin/util"
	"log"
)

func (r *RabbitMQ) Consume() {

	conn, ch := ConnectMQ()
	defer CloseMQ(conn, ch)

	// log.Println(ch)
	q, err := ch.QueueDeclare(
		r.QueueName, // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	util.FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	util.FailOnError(err, "Failed to declare a queue")

	k := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("1.Received a message: %s", d.Body)
			// d.Ack(false)
		}
	}()

	// go func() {
	//  for d := range msgs {
	//   log.Printf("2.Received a message: %s", d.Body)
	//   // d.Ack(false)
	//  }
	// }()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-k
}
