package main

import (
	"github.com/streadway/amqp"
	"log"
	"time"
)

type AmqpWrapper struct {
	url        string
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewAmqpWrapper(url string) *AmqpWrapper {
	return &AmqpWrapper{
		url: url,
	}
}

func (a *AmqpWrapper) Connect() error {
	connection, err := amqp.Dial(a.url)
	if err != nil {
		return err
	}

	channel, err := connection.Channel()
	if err != nil {
		return err
	}

	a.connection = connection
	a.channel = channel

	if err := a.defineExchange(); err != nil {
		return err
	}

	return nil
}

func (a *AmqpWrapper) reconnect() {
	for {

		if err := a.Connect(); err != nil {
			log.Println("ERR [AmqpWrapper::Reconnect]", err)
			time.Sleep(1 * time.Second)
		}
	}
}

func (a *AmqpWrapper) defineExchange() error {
	return a.channel.ExchangeDeclare(*EXCHANGE_NAME, "fanout", true, false, false, false, nil)
}

// Auto-reconnects and auto-retries
func (a *AmqpWrapper) Publish(message *amqp.Publishing) {
	for {
		log.Println("DEBUG [AmqpWrapper::Publish] Publishing message")

		if err := a.channel.Publish(*EXCHANGE_NAME, "", false, false, *message); err != nil {
			log.Println("WARN [AmqpWrapper::Publish]", err)
			a.reconnect()
		} else {
			break
		}
	}
}
