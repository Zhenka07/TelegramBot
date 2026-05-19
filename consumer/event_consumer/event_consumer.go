package eventconsumer

import (
	"log"
	"time"

	"github.com/Zhenka07/TelegramBot/client/events"
)

type Consumer struct {
	Processor events.Processor
	Fetcher   events.Fetcher
	batchSize int
}

func New(processor events.Processor, fetcher events.Fetcher, BatchSize int) Consumer {
	return Consumer{
		Processor: processor,
		Fetcher:   fetcher,
		batchSize: BatchSize,
	}
}

func (c Consumer) Start() error {
	for {
		events, err := c.Fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERROR] consumer: ", err.Error())
			continue
		}

		if len(events) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		if err := c.handleEvents(events); err != nil {
			log.Println(err)
			continue
		}

	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event %s", event.Text)

		if err := c.Processor.Process(event); err != nil {
			log.Printf("can't handle event '%s': %s", event.Text, err.Error())
			continue
		}
	}
	return nil
}
