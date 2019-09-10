package main

import (
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

const LISTEN string = "consume"
const SEND string = "publish"

func main() {
	queueName := flag.String("queue", "", "target queue name")
	mode := flag.String("mode", SEND, "test mode (consume or publish)")
	username := flag.String("username", "guest", "set username (default is guest)")
	password := flag.String("password", "guest", "set password (default is guest)")
	data := flag.String("data", "", "set data to send (must in JSON format)")

	flag.Parse()

	url := flag.Arg(len(flag.Args()) - 1)

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", *username, *password, url))

	if err != nil {
		log.Fatalf("Error make connection : %s", err)
	}

	channel, err := conn.Channel()

	if err != nil {
		log.Fatalf("Error create channel : %s", err)
	}

	if *mode == LISTEN {
		consumer, err := channel.Consume(*queueName, "", true, false, false, false, nil)

		if err != nil {
			log.Fatalf("Error create consumer : %s", err)
		}

		wait := make(chan bool)

		go func() {
			for msg := range consumer {
				fmt.Printf("Data : %s\n\n", msg.Body)
			}
			wait <- true
		}()

		<-wait
	} else {
		err = channel.Publish("", *queueName, false, false, amqp.Publishing{
			Body: []byte(*data),
		})

		if err != nil {
			log.Fatalf("Error send data : %s", err)
		}
	}
}
