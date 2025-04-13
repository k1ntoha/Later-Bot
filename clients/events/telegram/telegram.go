package telegram

import (
	"errors"

	"github.com/k1ntoha/LaterBot/clients/events"
	"github.com/k1ntoha/LaterBot/clients/telegram"
	"github.com/k1ntoha/LaterBot/lib/e"
	"github.com/k1ntoha/LaterBot/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

var (
	ErrUnknownEvent    = errors.New("unknown event type")
	ErrUnknownMetaType = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Update(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("Failed to Fetch ", err)
	}
	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))
	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1
	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("Failed to Process message", ErrUnknownEvent)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("Failed to process message", err)
	}
	if err := p.doCmd(event.Text, meta.ChatId, meta.Username); err != nil {
		return e.Wrap("failed to handle message", err)
	}
	panic("unimplemented")
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("Failed to get Meta", ErrUnknownMetaType)
	}
	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}
	if updType == events.Message {
		res.Meta = Meta{
			ChatId:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}
	return res
}

type Meta struct {
	ChatId   int
	Username string
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}
