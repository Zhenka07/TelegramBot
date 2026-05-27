.PHONY: run build setup

run:
	go run main.go

build:
	go build -o bot_app main.go

setup:
	cp .env.example .env
	@echo "Setup complete! Please add your TELEGRAM_API_KEY to the .env file."