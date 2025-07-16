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
		commandHandlers: new(commandRegistry),

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
	commandHandlers *commandRegistry
	bufSize         int
	workerPool      int

	// runtime
	chUpdate chan Update
	ctx      context.Context
	cancel   context.CancelFunc

	initErr error
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

// Usage:
//
//	HandleCommand("/foo", "foo description", fooHandler) // create a bot command for default scope and handle it with fooHandler
//	HandleCommand("/foo", "foo desciption", fooHandler, BotCommandScopeAllPrivateChats) // create a bot command for only private chats and handle it with fooHandler
func (b *Bot) HandleCommand(cmd, desc string, handler HandlerFunc, scopes ...BotCommandScope) *Bot {
	if cmd == "" {
		b.initErr = fmt.Errorf("cmd must be non-empty")
	}

	if b.commandHandlers == nil {
		b.commandHandlers = new(commandRegistry)
	}

	if !strings.HasPrefix(cmd, "/") {
		cmd = "/" + cmd
	}

	if len(scopes) == 0 {
		b.commandHandlers.AddCommand(cmd, desc, handler)
	} else {
		keys := make([]scopeKey, 0, len(scopes))

		for _, scope := range scopes {
			switch s := scope.(type) {
			case botCommandScopeNoParams:
				keys = append(keys, scopeKey{Scope: s.Scope()})
			case BotCommandScopeChat:
				keys = append(keys, scopeKey{Scope: s.Scope(), ChatID: string(s)})
			case BotCommandScopeChatAdministrators:
				keys = append(keys, scopeKey{Scope: s.Scope(), ChatID: string(s)})
			case BotCommandScopeChatMember:
				keys = append(keys, scopeKey{Scope: s.Scope(), ChatID: s.ChatID, UserID: s.UserID})
			}
		}

		b.commandHandlers.AddCommand(cmd, desc, handler, keys...)
	}

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
	if b.initErr != nil {
		return fmt.Errorf("configuration error: %w", b.initErr)
	}

	// creating context, channels, settting defaults, etc etc...
	b.init()

	// checking if we can launch the bot
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

	// filtering updates that we're not handling
	for upd := range b.updateHandlers {
		b.Supplier.AllowUpdate(upd)
	}

	// adding the list of handled commands to the bot menu on the client side
	if err = b.setupCommands(); err != nil {
		return fmt.Errorf("setting up commands: %w", err)
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
		b.commandHandlers = new(commandRegistry)
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

func (b *Bot) setupCommands() error {
	scopes := b.commandHandlers.GetScopes()

	if scopes == nil {
		return nil // early exit, nothing to do
	}

	// FIXME: should rewrite scopes only if getCommands returns
	// a list of commands that differs from what we have

	for _, scope := range scopes {
		cmds := b.commandHandlers.GetCommands(scope)
		if cmds == nil {
			continue // it can't be an error, can it?
		}

		bc := make([]BotCommand, 0, len(cmds))
		for _, cmd := range cmds {
			bc = append(bc, BotCommand{
				Command:     cmd.Name,
				Description: cmd.Description,
			})
		}

		smc := SetMyCommands{
			Commands: bc,
		}

		switch scope.Scope {
		case "default":
			smc.Scope = BotCommandScopeDefault
		case "all_private_chats":
			smc.Scope = BotCommandScopeAllPrivateChats
		case "all_group_chats":
			smc.Scope = BotCommandScopeAllGroupChats
		case "all_chat_administrators":
			smc.Scope = BotCommandScopeAllChatAdministrators
		case "chat":
			smc.Scope = BotCommandScopeChat(scope.ChatID)
		case "chat_administrators":
			smc.Scope = BotCommandScopeChatAdministrators(scope.ChatID)
		case "chat_member":
			smc.Scope = BotCommandScopeChatMember{ChatID: scope.ChatID, UserID: scope.UserID}

		default:
			return fmt.Errorf("unknown bot command scope: %s", scope.Scope)
		}

		resp, err := b.Sender.Send(&smc)
		if err != nil {
			return fmt.Errorf("sending setMyCommands request: %w", err)
		}

		if err := resp.GetError(); err != nil {
			return fmt.Errorf("failed to set commands: %w", err)
		}
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

				handler, ok := b.commandHandlers.GetHandler(cmd)
				if ok {
					handler(ctx)
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
