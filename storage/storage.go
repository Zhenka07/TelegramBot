package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
)

var ErrNoSavedPage = errors.New("You don't have saved page")

type Storage interface {
	Save(p *Page) error
	PickRandom(user_name string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", fmt.Errorf("Can't create hash %w", err)
	}

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", fmt.Errorf("Can't create hash %w", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
