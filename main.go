package main

import (
	"flag"
	"log"

	"github.com/k1ntoha/LaterBot/clients/events/telegram"
	tgClient "github.com/k1ntoha/LaterBot/clients/telegram"
	event_consumer "github.com/k1ntoha/LaterBot/consumer/event-consumer"
	"github.com/k1ntoha/LaterBot/storage/files"
)

const (
	host        = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {

	eventsProcessor := telegram.New(
		tgClient.New(host, mustToken()),
		files.New(storagePath),
	)

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
	log.Print("Service started...")
}

func mustToken() string {
	token := flag.String(
		"token",
		"",
		"token for access to telegram bot",
	)
	flag.Parse()

	if *token == "" {
		log.Fatal("Invalid token")
	}
	return *token
}
