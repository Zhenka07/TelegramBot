package main

import (
	"log"
	tgEvent "main/client/events/telegram"
	tgClient "main/client/telegram"
	eConsumer "main/consumer/event_consumer"
	"main/storage/files"
	"os"

	"github.com/joho/godotenv"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	key := os.Getenv("TELEGRAM_API_KEY")
	if key == "" {
		log.Fatal("TELEGRAM_API_KEY is not set in .env")
	}

	log.Print("service started")

	tg_client := tgClient.New(tgBotHost, key)
	processor := tgEvent.New(tg_client, files.New(storagePath))
	consumer := eConsumer.New(processor, processor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}
