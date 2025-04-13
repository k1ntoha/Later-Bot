package telegram

import (
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/k1ntoha/LaterBot/lib/e"
	"github.com/k1ntoha/LaterBot/storage"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("New message %s , fro, %s", text, username)

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}

	return nil
}

func (p *Processor) sendRandom(chatId int, username string) (err error) {
	defer func() { err = e.Wrap("Failed to send random message from list", err) }()

	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatId, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatId, page.URL); err != nil {
		return err
	}
	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatId int) error {
	return p.tg.SendMessage(chatId, msgHelp)
}

func (p *Processor) sendHello(chatId int) error {
	return p.tg.SendMessage(chatId, msgHello)
}

func (p *Processor) savePage(chatId int, pageUrl string, username string) (err error) {
	defer func() { err = e.Wrap("Failed to save page", err) }()
	page := &storage.Page{
		URL:      pageUrl,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessage(chatId, msgAlreadyExists)
	}
	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatId, msgSaved); err != nil {
		return err
	}
	return nil
}

func isAddCmd(text string) bool {
	return isUrl(text)

}

func isUrl(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
