package main

import (
	"context"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	_ "github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

var mqURI = "amqp://user:password@localhost:5672/my_vhost"

func main() {
	logger := watermill.NewStdLogger(false, false)

	// RabbitMQ connection config
	pubConfig := amqp.NewDurablePubSubConfig(mqURI, amqp.GenerateQueueNameTopicName)
	pubConfig.Exchange.Type = "fanout"
	// Create a publisher
	publisher, err := amqp.NewPublisher(pubConfig, logger)
	if err != nil {
		panic(err)
	}
	defer publisher.Close()

	// Create two subscribers (they will receive the same messages)
	subscriber1, err := createSubscriber("subscriber-1", logger)
	if err != nil {
		panic(err)
	}
	defer subscriber1.Close()

	subscriber2, err := createSubscriber("subscriber-2", logger)
	if err != nil {
		panic(err)
	}
	defer subscriber2.Close()

	// Topic name for our events
	topic := "events.notifications"

	// Process messages for subscriber 1
	go processMessages(subscriber1, topic, "Subscriber 1")

	// Process messages for subscriber 2
	go processMessages(subscriber2, topic, "Subscriber 2")

	// Give some time for subscribers to connect
	time.Sleep(time.Second)

	// Publish a few messages
	var i = 0
	for {
		i++
		msg := message.NewMessage(
			watermill.NewUUID(),
			[]byte("Hello, this is message #"+string(rune(i+'0'))),
		)

		if err := publisher.Publish(topic, msg); err != nil {
			panic(err)
		}
		log.Printf("Published message: %s", msg.Payload)
		time.Sleep(time.Millisecond * 100)
	}
}

func createSubscriber(consumerGroup string, logger watermill.LoggerAdapter) (*amqp.Subscriber, error) {
	// For the fanout behavior (each subscriber gets all messages),
	// we use different consumer groups
	return amqp.NewSubscriber(amqp.NewDurablePubSubConfig(mqURI, amqp.GenerateQueueNameTopicNameWithSuffix(consumerGroup)), logger)
}

func processMessages(subscriber *amqp.Subscriber, topic string, name string) {
	messages, err := subscriber.Subscribe(context.Background(), topic)
	if err != nil {
		panic(err)
	}
	for msg := range messages {
		log.Printf("[%s] received message: %s", name, string(msg.Payload))
		msg.Ack()
	}
}
