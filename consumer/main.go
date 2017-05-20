package main

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type Message struct {
	Message string
	Value   int
	Address int
	Created time.Time
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// go binary decoder
func FromGOB64(str string) Message {
	m := Message{}
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		fmt.Println("failed base64 Decode", err)
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(&m)
	if err != nil {
		fmt.Println("failed gob Decode", err)
	}
	return m
}

func init() {
	log.Printf("Init")
	gob.Register(Message{})
}

func main() {
	conn, err := amqp.Dial("amqp://rabbitmq:rabbitmq@rabbit1:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	devices := make(map[int][]string)

	cluster := gocql.NewCluster("cassandra-1")
	cluster.Keyspace = "demo"
	cluster.Consistency = gocql.LocalOne
	session, _ := cluster.CreateSession()
	
	forever := make(chan bool)

	canShow := true

	go func() {
		for d := range msgs {

			m := FromGOB64(string(d.Body))
			//log.Printf("Message: %v\n", m)
			err := session.Query("INSERT INTO messages (address, message, value, created) VALUES (?, ?, ?, ?)", m.Address, m.Message, m.Value, m.Created).Exec()
			if err != nil {
				log.Fatal(err)
			}
			
			devices[m.Address] = append(devices[m.Address], m.Message)

			if canShow {
				timeout := make(chan bool, 1)
				go func() {
					canShow = false
					log.Printf("--------------------------\n")
					time.Sleep(3 * time.Second)
					for k, v := range devices {
						log.Printf("-Device %d Received %d messages\n", k, len(v))
					}
					canShow = true
					timeout <- true
				}()
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
