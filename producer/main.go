package main

import (
	"fmt"
	"log"
	"time"
	"github.com/streadway/amqp"
	"encoding/gob"
	"encoding/base64"
	"bytes"
	"math/rand"
)

type Message struct{
	Message string
	Value int
	Address int
	Created time.Time
}


func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func ToGOB64(m Message) string {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(m)
	if err != nil {
		fmt.Println("failed gob Encode", err)
	}
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

func simulateValue(current int) int{
	current = current + 1
	return current
}


func init() {
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
	
	rand.Seed(time.Now().UTC().UnixNano())
	id := rand.Intn(1000)
	value := 1
		
	for range time.Tick(500*time.Millisecond) {
		value = simulateValue(value)
	
		m := new(Message)
		m.Message = "Hello!"
		m.Address = id
		m.Created = time.Now()
		m.Value = value
		
		body := ToGOB64(*m)
		
		err := ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		log.Printf(" [x] Sent %s", body)
		failOnError(err, "Failed to publish a message")
	}
		
}
