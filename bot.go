package botify

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
)

type Bot struct {
	// configurable
	handlers   map[UpdateType]HandlerFunc
	sender     RequestSender
	supplier   UpdateSupplier
	bufSize    int
	workerPool int

	// runtime
	chUpdate   chan Update
	ctx        context.Context
	cancel     context.CancelFunc
	isOnline   bool
	hasWebhook bool
}

func (b *Bot) work() {
	for {
		select {
		case <-b.ctx.Done():
			return

		case upd := <-b.chUpdate:
			ctx := Context{
				bot:     b,
				updType: upd.UpdateType(),
				upd:     &upd,
			}

			exact, ok := b.handlers[ctx.updType]
			if ok {
				exact(ctx)
			}

			all, ok := b.handlers[UpdateTypeAll]
			if ok {
				all(ctx)
			}
		}
	}
}

func (b *Bot) Handle(t UpdateType, handler HandlerFunc) {
	b.handlers[t] = handler
}

func (b *Bot) RequestSender() RequestSender {
	return b.sender
}

func (b *Bot) WithRequestSender(sender RequestSender) {
	b.sender = sender
	lps, ok := b.supplier.(*LongPollingSupplier)
	if ok {
		lps.Sender = sender
		b.supplier = lps
	}
}

func (b *Bot) UpdateSupplier() UpdateSupplier {
	return b.supplier
}

func (b *Bot) WithUpdateSupplier(supp UpdateSupplier) {
	b.supplier = supp
}

func (b *Bot) Serve() error {
	r, err := b.sender.Send(GetWebhookInfo)
	if err != nil {
		return fmt.Errorf("requesting for webhook info: %w", err)
	}

	var wh WebhookInfo
	err = r.BindResult(&wh)
	if err != nil {
		return err
	}

	if _, ok := b.supplier.(*LongPollingSupplier); ok && wh.URL != "" {
		return fmt.Errorf("can't use long polling when webhook is set; use deleteWebhook before running long polling bot")
	}

	ctx, cancel := context.WithCancel(context.Background())
	b.ctx = ctx
	b.cancel = cancel

	b.chUpdate = make(chan Update, b.bufSize)

	for range b.workerPool {
		go b.work()
	}
	return b.supplier.GetUpdates(b.ctx, b.chUpdate)
}

func DefaultLongPollingBot(token string) *Bot {
	sender := DefaultRequestSender{
		Client:   http.DefaultClient,
		APIToken: token,
		APIHost:  "https://api.telegram.org/",
	}
	bot := Bot{
		handlers: make(map[UpdateType]HandlerFunc),
		sender:   &sender,
		supplier: &LongPollingSupplier{
			Sender:         &sender,
			Offset:         0,
			Timeout:        30,
			Limit:          100,
			AllowedUpdates: &[]string{},
		},
		bufSize:    0,
		workerPool: runtime.NumCPU(),
	}
	return &bot
}
