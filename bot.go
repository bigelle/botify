package botify

import (
	"context"
	"fmt"
	"net/url"
	"runtime"
)

type Bot struct {
	// configurable
	Handlers   map[UpdateType]HandlerFunc
	Sender     RequestSender
	Supplier   UpdateSupplier
	BufSize    int
	WorkerPool int

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

			exact, ok := b.Handlers[ctx.updType]
			if ok {
				exact(ctx)
			}

			all, ok := b.Handlers[UpdateTypeAll]
			if ok {
				all(ctx)
			}
		}
	}
}

func (b *Bot) Handle(t UpdateType, handler HandlerFunc) {
	if b.Handlers == nil {
		b.Handlers = make(map[UpdateType]HandlerFunc)
	}
	b.Handlers[t] = handler
}

func (b *Bot) Serve() error {
	r, err := b.Sender.Send(GetWebhookInfo)
	if err != nil {
		return fmt.Errorf("requesting for webhook info: %w", err)
	}

	var wh WebhookInfo
	err = r.BindResult(&wh)
	if err != nil {
		return fmt.Errorf("reading API response: %w", err)
	}

	if _, ok := b.Supplier.(*LongPollingSupplier); ok && wh.URL != "" {
		return fmt.Errorf("can't use long polling when webhook is set; use deleteWebhook before running long polling bot")
	}

	if ws, ok := b.Supplier.(*WebhookSupplier); ok && wh.URL == "" {
		whURL := ws.Domain + ws.Path
		_, err := url.Parse(whURL)
		if err != nil {
			return fmt.Errorf("invalid webhook URL: %w", err)
		}

		swh := SetWebhook{
			URL:                whURL,
			Certificate:        ws.Certificate,
			IPAddress:          ws.IPAddress,
			MaxConnections:     ws.MaxConnections,
			AllowedUpdates:     ws.AllowedUpdates,
			DropPendingUpdates: ws.DropPendingUpdates,
			SecretToken:        ws.SecretToken,
		}

		r, err = b.Sender.Send(&swh)
		if err != nil {
			return fmt.Errorf("setting webhook: %w", err)
		}

		var whOk bool
		if err = r.BindResult(&whOk); err != nil {
			return fmt.Errorf("reading API response: %w", err)
		}

		if !whOk {
			err = r.GetError()
			return fmt.Errorf("failed to set webhook: %w", err)
		}
	}

	for upd := range b.Handlers {
		b.Supplier.AllowUpdate(upd)
	}

	ctx, cancel := context.WithCancel(context.Background())
	b.ctx = ctx
	b.cancel = cancel

	b.chUpdate = make(chan Update, b.BufSize)

	if b.WorkerPool == 0 {
		b.WorkerPool = runtime.NumCPU()
	}

	for range b.WorkerPool {
		go b.work()
	}

	return b.Supplier.GetUpdates(b.ctx, b.chUpdate)
}

func DefaultLongPollingBot(token string) *Bot {
	sender := DefaultRequestSender(token)

	bot := Bot{
		Handlers: make(map[UpdateType]HandlerFunc),
		Sender:   sender,
		Supplier: &LongPollingSupplier{
			Sender:         sender,
			Offset:         0,
			Timeout:        30,
			Limit:          100,
			AllowedUpdates: &[]string{},
		},
		BufSize:    0,
		WorkerPool: runtime.NumCPU(),
	}

	return &bot
}
