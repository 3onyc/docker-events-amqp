package main

import (
	"encoding/json"
	"flag"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/streadway/amqp"
	"log"
	"time"
)

var (
	DOCKER_URL    = flag.String("docker", "", "Docker URL")
	AMQP_URL      = flag.String("amqp", "", "AMQP URL")
	EXCHANGE_NAME = flag.String("exchange", "docker.events", "Exchange Name")
)

func eventHandler(
	dockerEvents chan *docker.APIEvents,
	amqpWrapper *AmqpWrapper,
) {
	for {
		select {
		case event := <-dockerEvents:
			log.Printf("DEBUG [eventHandler] received (%s) %s", event.Status, event.ID)

			if event == docker.EOFEvent {
				break
			}

			encoded, err := json.Marshal(event)
			if err != nil {
				log.Println("ERROR [eventHandler]", err)
				continue
			}

			amqpWrapper.Publish(&amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				Timestamp:    time.Now(),
				ContentType:  "text/json",
				Body:         encoded,
			})
		}
	}
}

func parseArgs() {
	flag.Parse()

	if *DOCKER_URL == "" {
		log.Fatal("Please specify Docker URL using -docker")
	}

	if *AMQP_URL == "" {
		log.Fatal("Please specify AMQP URL using -amqp")
	}

}

func main() {
	parseArgs()

	dockerClient, err := connectDocker(*DOCKER_URL)
	if err != nil {
		log.Fatal("ERROR [main::docker]", err)
	}

	amqpWrapper := NewAmqpWrapper(*AMQP_URL)
	if err := amqpWrapper.Connect(); err != nil {
		log.Fatal("ERROR [main::amqp]", err)
	}

	if err := start(amqpWrapper, dockerClient); err != nil {
		log.Fatal("ERROR [main]", err)
	}
}

func start(amqpWrapper *AmqpWrapper, dockerClient *docker.Client) error {

	eventChan := make(chan *docker.APIEvents, 100)
	if err := dockerClient.AddEventListener(eventChan); err != nil {
		return err
	}

	for {
		log.Println("INFO [start] Starting event loop...")
		eventHandler(eventChan, amqpWrapper)
		time.Sleep(1 * time.Second)
	}

	return nil
}
