package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"main/storage"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, ChatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("Get out new command '%s' from '%s'", text, username)

	if IsAddCmd(text) {
		return p.AddPage(ChatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.SendRandom(ChatID, username)
	case HelpCmd:
		return p.SendHelp(ChatID)
	case StartCmd:
		return p.SendHello(ChatID)
	default:
		return p.tg.SendMessage(ChatID, msgUnknownCommand)
	}

}

func (p *Processor) AddPage(ChatId int, page_url string, username string) error {
	page := &storage.Page{
		URL:      page_url,
		UserName: username,
	}

	is_exists, err := p.storage.IsExists(page, context.Background())
	if err != nil {
		return fmt.Errorf("Can't save this page %w", err)
	}

	if is_exists {
		return p.tg.SendMessage(ChatId, msgAlreadyExists)
	}

	if err := p.storage.Save(page, context.Background()); err != nil {
		return fmt.Errorf("Can't save this page %w", err)
	}

	if err := p.tg.SendMessage(ChatId, msgSaved); err != nil {
		return fmt.Errorf("Can't save this page %w", err)
	}
	return nil
}

func (p *Processor) SendRandom(chatId int, username string) error {
	page, err := p.storage.PickRandom(username, context.Background())
	if err != nil && !errors.Is(err, storage.ErrNoSavedPage) {
		return fmt.Errorf("Can't send random page %w", err)
	}
	if errors.Is(err, storage.ErrNoSavedPage) {
		return p.tg.SendMessage(chatId, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatId, page.URL); err != nil {
		return fmt.Errorf("Can't send random page %w", err)
	}

	return p.storage.Remove(page, context.Background())
}

func (p *Processor) SendHelp(chatId int) error {
	return p.tg.SendMessage(chatId, msgHelp)
}

func (p *Processor) SendHello(chatId int) error {
	return p.tg.SendMessage(chatId, msgHello)
}

func IsAddCmd(text string) bool {
	return IsUrl(text)
}

func IsUrl(text string) bool {
	url, err := url.Parse(text)
	return err == nil && url.Host != ""
}
