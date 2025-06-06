package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"

	"github.com/k1ntoha/LaterBot/lib/e"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

var ErrNoSavedPages = errors.New("no saved pages")

func (p Page) Hash() (string, error) {
	h := sha1.New()
	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("Failed to Hash", err)
	}
	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("Failed to Hash", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
