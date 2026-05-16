package main

import (
	"context"
	"log"
	tgEvent "main/client/events/telegram"
	tgClient "main/client/telegram"
	eConsumer "main/consumer/event_consumer"
	bdStorage "main/storage/sqlite"
	"os"

	"github.com/joho/godotenv"
)

const (
	tgBotHost      = "api.telegram.org"
	storagePath    = "user_data/local"
	storageSqlPath = "user_data/sqlite/storage.db"
	batchSize      = 100
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

	tgClient := tgClient.New(tgBotHost, key)
	st, err := bdStorage.New(storageSqlPath)
	if err != nil {
		log.Fatal("can't start database", err)
	}
	st.Init(context.Background())

	processor := tgEvent.New(tgClient, st)
	consumer := eConsumer.New(processor, processor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}
