package botify

import (
	"context"
	"fmt"
	"net/url"
	"runtime"
	"strings"
)

func DefaultBot(token string) *Bot {
	sender := DefaultRequestSender(token)

	bot := Bot{
		Token:    token,
		Handlers: make(map[string]HandlerFunc),
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

type Bot struct {
	// configurable
	Token       string
	Handlers    map[string]HandlerFunc
	Sender      RequestSender
	Supplier    UpdateSupplier

	BufSize    int
	WorkerPool int

	// runtime
	chUpdate chan Update
	ctx      context.Context
	cancel   context.CancelFunc
}

func (b *Bot) Handle(t string, handler HandlerFunc) {
	if b.Handlers == nil {
		b.Handlers = make(map[string]HandlerFunc)
	}

	b.Handlers[t] = handler
}

func (b *Bot) HandleCommand(cmd string, handler HandlerFunc) {
	if b.Handlers == nil {
		b.Handlers = make(map[string]HandlerFunc)
	}
	
	if !strings.HasPrefix(cmd, "/") {
		cmd = "/" + cmd
	}

	b.Handlers[cmd] = handler
}

// TODO: simplify it
func (b *Bot) Serve() error {
	b.init()

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
		whURL := ws.WebhookURL()

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

	defer b.Shutdown()

	for range b.WorkerPool {
		go b.work()
	}

	return b.Supplier.GetUpdates(b.ctx, b.chUpdate)
}

// TODO: make it more graceful
func (b *Bot) Shutdown() error {
	b.cancel()
	close(b.chUpdate)

	resp, err := b.Sender.Send(Close)
	if err != nil {
		return err
	}

	if !resp.IsSuccessful() {
		// i mean do i care if it was closed too early?
		if resp.ErrorCode != 429 {
			return resp.GetError()
		}
	}

	return nil
}

func (b *Bot) init() {
	if b.Token == "" {
		panic("API token must not be empty")
	}

	if b.Sender == nil {
		b.Sender = DefaultRequestSender(b.Token)
	}

	if b.Supplier == nil {
		b.Supplier = &LongPollingSupplier{
			Sender: b.Sender,

			AllowedUpdates: &[]string{},
			Timeout:        30,
			Offset:         0,
			Limit:          100,
		}
	}

	if b.Handlers == nil {
		b.Handlers = make(map[string]HandlerFunc)
	}

	if b.BufSize < 0 {
		b.BufSize = 0
	}

	if b.WorkerPool <= 0 {
		b.WorkerPool = runtime.NumCPU()
	}

	ctx, cancel := context.WithCancel(context.Background())
	b.ctx = ctx
	b.cancel = cancel

	b.chUpdate = make(chan Update, b.BufSize)
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

			handler, ok := b.Handlers[ctx.updType]
			if ok {
				handler(ctx)
			}

		}
	}
}
