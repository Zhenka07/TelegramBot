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

var (
	ErrURLTooLong    = errors.New("url is too long (max 2048 characters)")
	ErrInvalidScheme = errors.New("url scheme must be http or https")
	ErrInvalidHost   = errors.New("url host is invalid or missing")
	ErrForbiddenHost = errors.New("access to local/private hosts is forbidden")
)

func ValidateURL(rawURL string) error {
	rawURL = strings.TrimSpace(rawURL)
	if len(rawURL) > 2048 {
		return ErrURLTooLong
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return ErrInvalidHost
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return ErrInvalidScheme
	}

	hostname := u.Hostname()
	if hostname == "" {
		return ErrInvalidHost
	}

	forbidden := map[string]bool{
		"localhost": true,
		"127.0.0.1": true,
		"0.0.0.0":   true,
		"::1":       true,
	}
	if forbidden[strings.ToLower(hostname)] {
		return ErrForbiddenHost
	}

	return nil
}

func IsURL(text string) bool {
	return ValidateURL(text) == nil
}


