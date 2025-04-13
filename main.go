package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/k1ntoha/LaterBot/clients/telegram"
)

const (
	host = "api.telegram.org"
)

func main() {

	tgClient := telegram.New(host, mustToken())

	//fetcher = fetcher.New()

	//processor = processor.New(tgClient)

	//consumer.Start(fetcher, processor)
	fmt.Println(tgClient)
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)
	flag.Parse()

	if *token == "" {
		log.Fatal("Invalid token")
	}
	return *token
}
