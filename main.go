package main

import (
	"context"
	"log"

	tgEvent "github.com/Zhenka07/TelegramBot/client/events/telegram"
	tgClient "github.com/Zhenka07/TelegramBot/client/telegram"
	"github.com/Zhenka07/TelegramBot/config"
	eConsumer "github.com/Zhenka07/TelegramBot/consumer/event_consumer"
	bdStorage "github.com/Zhenka07/TelegramBot/storage/sqlite"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error in config file %w", err)
	}

	log.Print("service started")

	tgClient := tgClient.New(config.TelegramHost, config.TelegramAPIKey)

	st, err := bdStorage.New(config.SQLitePath)
	if err != nil {
		log.Fatal("can't start database", err)
	}
	st.Init(context.Background())

	processor, fetcher := tgEvent.New(tgClient, st)
	consumer := eConsumer.New(processor, fetcher, config.BatchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}
