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
		Token:  token,
		Sender: sender,
		Supplier: &LongPollingSupplier{
			Sender:         sender,
			Offset:         0,
			Timeout:        30,
			Limit:          100,
			AllowedUpdates: &[]string{},
		},

		updateHandlers:  make(map[string]HandlerFunc),
		commandHandlers: make(map[string]HandlerFunc),

		bufSize:    0,
		workerPool: runtime.NumCPU(),
	}

	return &bot
}

type Bot struct {
	// configurable
	Token    string
	Sender   RequestSender
	Supplier UpdateSupplier

	// only through methods
	updateHandlers  map[string]HandlerFunc
	commandHandlers map[string]HandlerFunc
	bufSize         int
	workerPool      int

	// runtime
	chUpdate chan Update
	ctx      context.Context
	cancel   context.CancelFunc
}

func (b *Bot) Handle(t string, handler HandlerFunc) *Bot {
	if b.updateHandlers == nil {
		b.updateHandlers = make(map[string]HandlerFunc)
	}

	if _, ok := allUpdTypes[t]; ok {
		b.updateHandlers[t] = handler
	}

	return b
}

func (b *Bot) HandleCommand(cmd string, handler HandlerFunc) *Bot {
	if b.commandHandlers == nil {
		b.commandHandlers = make(map[string]HandlerFunc)
	}

	if !strings.HasPrefix(cmd, "/") {
		cmd = "/" + cmd
	}

	b.commandHandlers[cmd] = handler

	return b
}

func (b *Bot) WithChannelSize(l int) *Bot {
	if l >= 0 {
		b.bufSize = l
	}

	return b
}

func (b *Bot) WithWorkerPool(l int) *Bot {
	if l > 0 {
		b.workerPool = l
	}

	return b
}

// TODO: simplify it

func (b *Bot) Serve() error {
	b.init()

	wh, err := b.getWebhookInfo()
	if err != nil {
		// FIXME: there is a better way to handle this
		return fmt.Errorf("getting webhook info: %w", err)
	}

	if _, ok := b.Supplier.(*LongPollingSupplier); ok && wh.URL != "" {
		return fmt.Errorf("can't use long polling when webhook is set; use deleteWebhook before running long polling bot")
	}

	if _, ok := b.Supplier.(*WebhookSupplier); ok && wh.URL == "" {
		if err = b.setWebhook(); err != nil {
			return fmt.Errorf("setting webhook: %w", err)
		}
	}

	for upd := range b.updateHandlers {
		b.Supplier.AllowUpdate(upd)
	}

	defer b.Shutdown()

	for range b.workerPool {
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

	if b.updateHandlers == nil {
		b.updateHandlers = make(map[string]HandlerFunc)
	}
	if b.commandHandlers == nil {
		b.commandHandlers = make(map[string]HandlerFunc)
	}

	if b.bufSize < 0 {
		b.bufSize = 0
	}

	if b.workerPool <= 0 {
		b.workerPool = runtime.NumCPU()
	}

	ctx, cancel := context.WithCancel(context.Background())
	b.ctx = ctx
	b.cancel = cancel

	b.chUpdate = make(chan Update, b.bufSize)
}

func (b *Bot) getWebhookInfo() (*WebhookInfo, error) {
	r, err := b.Sender.Send(GetWebhookInfo)
	if err != nil {
		return nil, fmt.Errorf("requesting for webhook info: %w", err)
	}

	var wh WebhookInfo
	err = r.BindResult(&wh)
	if err != nil {
		return nil, fmt.Errorf("reading API response: %w", err)
	}

	return &wh, nil
}

func (b *Bot) setWebhook() error {
	ws := b.Supplier.(*WebhookSupplier)
	var r *APIResponse

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
		return fmt.Errorf("sending request: %w", err)
	}

	var whOk bool
	if err = r.BindResult(&whOk); err != nil {
		return fmt.Errorf("reading API response: %w", err)
	}

	if !whOk {
		err = r.GetError()
		return fmt.Errorf("failed to set webhook: %w", err)
	}

	return nil
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
				ctx:     b.ctx,
			}

			if upd.Message != nil && upd.Message.IsCommand() {
				cmd, _ := upd.Message.GetCommand()
				cmdHandler, ok := b.commandHandlers[cmd]

				if ok {
					cmdHandler(ctx)
				}

				return
			}

			handler, ok := b.updateHandlers[ctx.updType]
			if ok {
				handler(ctx)
			}

		}
	}
}
