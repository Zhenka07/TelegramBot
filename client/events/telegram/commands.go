package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/Zhenka07/TelegramBot/storage"
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

func (p *Processor) AddPage(ChatID int, PageUrl string, username string) error {
	page := &storage.Page{
		URL:      PageUrl,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(context.Background(), page)
	if err != nil {
		return fmt.Errorf("Can't save this page %w", err)
	}

	if isExists {
		return p.tg.SendMessage(ChatID, msgAlreadyExists)
	}

	if err := p.storage.Save(context.Background(), page); err != nil {
		return fmt.Errorf("Can't save this page %w", err)
	}

	if err := p.tg.SendMessage(ChatID, msgSaved); err != nil {
		return fmt.Errorf("Can't save this page %w", err)
	}
	return nil
}

func (p *Processor) SendRandom(ChatId int, username string) error {
	page, err := p.storage.PickRandom(context.Background(), username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPage) {
		return fmt.Errorf("Can't send random page %w", err)
	}
	if errors.Is(err, storage.ErrNoSavedPage) {
		return p.tg.SendMessage(ChatId, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(ChatId, page.URL); err != nil {
		return fmt.Errorf("Can't send random page %w", err)
	}

	return p.storage.Remove(context.Background(), page)
}

func (p *Processor) SendHelp(ChatId int) error {
	return p.tg.SendMessage(ChatId, msgHelp)
}

func (p *Processor) SendHello(ChatId int) error {
	return p.tg.SendMessage(ChatId, msgHello)
}

func IsAddCmd(text string) bool {
	return IsURL(text)
}

func IsURL(text string) bool {
	u, err := url.Parse(text)
	if err != nil || u.Host == "" {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}
