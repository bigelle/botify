package botify

import (
	"context"
	"fmt"
	"runtime"
	"slices"
	"strings"
)

func DefaultBot(token string) *Bot {
	sender := DefaultRequestSender(token)

	bot := Bot{
		Token:  token,
		Sender: sender,
		Receiver: &LongPolling{
			Offset:  0,
			Timeout: 30,
			Limit:   100,
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
	Receiver UpdateReceiver

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

// LocaleMap is used to provide a localization for the command.
// The key must be a two-letter ISO 639-1 language code.
// The value must be 1-256 characters long.
type LocaleMap map[string]string

// NOTE: if, for example, bot has commands in both private chat and default scope, and you're opening the bot menu in a private chat,
// you would see only private chat commands.
//
// See [Determining list of commands] for details.
//
// [Determining list of commands]: https://core.telegram.org/bots/api#determining-list-of-commands
func (b *Bot) HandleCommandWithLocales(cmd string, locales LocaleMap, handler HandlerFunc, scopes ...BotCommandScope) *Bot {
	if cmd == "" {
		b.initErr = fmt.Errorf("cmd must be non-empty")
		return b
	}

	if b.commandHandlers == nil {
		b.commandHandlers = new(commandRegistry)
	}

	if !strings.HasPrefix(cmd, "/") {
		cmd = "/" + cmd
	}

	var (
		scope BotCommandScope
		keys  = make([]scopeKey, 0, len(scopes))
	)

	for code, desc := range locales {
		if len(code) != 2 {
			b.initErr = fmt.Errorf("language code must be a two-letter ISO 639-1: %s", code)
			return b
		}

		if len(desc) < 1 || 256 < len(desc) {
			b.initErr = fmt.Errorf("command description must be 1-256 characters")
			return b
		}

		if len(scopes) == 0 {
			b.commandHandlers.AddCommand(cmd, desc, handler)
		} else {

			if code == "en" {
				code = "" // to make sure that english description is applied by default
			}

			for _, scope = range scopes {
				switch s := scope.(type) {
				case botCommandScopeNoParams:
					keys = append(keys, scopeKey{Scope: s.Scope(), LanguageCode: code})
				case BotCommandScopeChat:
					keys = append(keys, scopeKey{Scope: s.Scope(), LanguageCode: code, ChatID: string(s)})
				case BotCommandScopeChatAdministrators:
					keys = append(keys, scopeKey{Scope: s.Scope(), LanguageCode: code, ChatID: string(s)})
				case BotCommandScopeChatMember:
					keys = append(keys, scopeKey{Scope: s.Scope(), LanguageCode: code, ChatID: s.ChatID, UserID: s.UserID})
				}
			}

			b.commandHandlers.AddCommand(cmd, desc, handler, keys...)

			keys = keys[:0]
		}
	}

	return b
}

// NOTE: if , for example, bot has commands in both private chat and default scope, and you're opening the bot menu in a private chat,
// you would see only private chat commands.
//
// See [Determining list of commands] for details.
//
// [Determining list of commands]: https://core.telegram.org/bots/api#determining-list-of-commands
func (b *Bot) HandleCommand(cmd, desc string, handler HandlerFunc, scopes ...BotCommandScope) *Bot {
	return b.HandleCommandWithLocales(cmd, LocaleMap{"en": desc}, handler, scopes...)
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
	if _, ok := b.Receiver.(*LongPolling); ok && wh.URL != "" {
		return fmt.Errorf("can't use long polling when webhook is set; use deleteWebhook before running long polling bot")
	}

	// adding the list of handled commands to the bot menu on the client side
	if err = b.setupCommands(); err != nil {
		return fmt.Errorf("setting up commands: %w", err)
	}

	defer b.Shutdown()

	for range b.workerPool {
		go b.work()
	}

	return b.Receiver.ReceiveUpdates(b.ctx, b.chUpdate)
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

	if b.Receiver == nil {
		b.Receiver = &LongPolling{
			Timeout: 30,
			Offset:  0,
			Limit:   100,
		}
	}
	b.Receiver.PairBot(b)

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

	b.ctx, b.cancel = context.WithCancel(context.Background())

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

var scopeMap = map[string]func(scopeKey) BotCommandScope{
	"default":                 func(scopeKey) BotCommandScope { return BotCommandScopeDefault },
	"all_private_chats":       func(scopeKey) BotCommandScope { return BotCommandScopeAllPrivateChats },
	"all_group_chats":         func(scopeKey) BotCommandScope { return BotCommandScopeAllGroupChats },
	"all_chat_administrators": func(scopeKey) BotCommandScope { return BotCommandScopeAllChatAdministrators },
	"chat":                    func(key scopeKey) BotCommandScope { return BotCommandScopeChat(key.ChatID) },
	"chat_administrators":     func(key scopeKey) BotCommandScope { return BotCommandScopeChatAdministrators(key.ChatID) },
	"chat_member": func(key scopeKey) BotCommandScope {
		return BotCommandScopeChatMember{ChatID: key.ChatID, UserID: key.UserID}
	},
}

func (b *Bot) setupCommands() error {
	scopes := b.commandHandlers.GetScopes()
	if len(scopes) == 0 {
		return nil
	}

	for _, scope := range scopes {
		if err := b.syncCommandsByScope(scope); err != nil {
			return fmt.Errorf("syncing commands for scope %s: %w", scope.Scope, err)
		}
	}
	return nil
}

func (b *Bot) syncCommandsByScope(key scopeKey) error {
	scopeFunc, exists := scopeMap[key.Scope]
	if !exists {
		return fmt.Errorf("unknown bot command scope: %s", key.Scope)
	}

	scope := scopeFunc(key)

	currentCommands, err := b.getCurrentCommands(scope)
	if err != nil {
		return fmt.Errorf("getting current commands: %w", err)
	}

	myCommands := b.commandHandlers.GetCommands(key)

	if !isEqualCommands(myCommands, currentCommands) {
		if err = b.setCommands(scope, myCommands); err != nil {
			return fmt.Errorf("setting commands: %w", err)
		}
	}

	return nil
}

func (b *Bot) getCurrentCommands(scope BotCommandScope) ([]BotCommand, error) {
	gmc := GetMyCommands{Scope: scope}

	resp, err := b.Sender.Send(&gmc)
	if err != nil {
		return nil, fmt.Errorf("sending getMyCommands request: %w", err)
	}

	if err = resp.GetError(); err != nil {
		return nil, fmt.Errorf("getting current commands: %w", err)
	}

	var commands []BotCommand
	if err = resp.BindResult(&commands); err != nil {
		return nil, fmt.Errorf("binding getMyCommands result: %w", err)
	}

	return commands, nil
}

func (b *Bot) setCommands(scope BotCommandScope, commands []command) error {
	botCommands := make([]BotCommand, 0, len(commands))

	for _, cmd := range commands {
		botCommands = append(botCommands, BotCommand{
			Command:     cmd.Name,
			Description: cmd.Description,
		})
	}

	smc := SetMyCommands{
		Scope:    scope,
		Commands: botCommands,
	}

	resp, err := b.Sender.Send(&smc)
	if err != nil {
		return fmt.Errorf("sending setMyCommands request: %w", err)
	}

	if err = resp.GetError(); err != nil {
		return fmt.Errorf("setting bot commands: %w", err)
	}

	return nil
}

func isEqualCommands(myCommands []command, telegramCommands []BotCommand) bool {
	if len(myCommands) != len(telegramCommands) {
		return false
	}

	if len(myCommands) == 0 {
		return true
	}

	mySlice := make([]BotCommand, len(myCommands))
	for i, cmd := range myCommands {
		mySlice[i] = BotCommand{
			Command:     cmd.Name,
			Description: cmd.Description,
		}
	}

	telegramSlice := make([]BotCommand, len(telegramCommands))
	copy(telegramSlice, telegramCommands)

	compareFunc := func(a, b BotCommand) int {
		if cmp := strings.Compare(a.Command, b.Command); cmp != 0 {
			return cmp
		}
		return strings.Compare(a.Description, b.Description)
	}

	slices.SortFunc(mySlice, compareFunc)
	slices.SortFunc(telegramSlice, compareFunc)

	return slices.Equal(mySlice, telegramSlice)
}

func (b *Bot) work() {
	var (
		ctx     Context
		cmd     string
		handler HandlerFunc
		ok      bool
	)

	for {
		select {
		case <-b.ctx.Done():
			return

		case upd := <-b.chUpdate:
			ctx = Context{
				bot:     b,
				updType: upd.UpdateType(),
				upd:     &upd,
				ctx:     b.ctx,
			}

			if upd.Message != nil && upd.Message.IsCommand() {
				cmd, _ = upd.Message.GetCommand()

				handler, ok = b.commandHandlers.GetHandler(cmd)
				if ok {
					handler(ctx)
				}
			} else {
				handler, ok = b.updateHandlers[ctx.updType]
				if ok {
					handler(ctx)
				}
			}

		}
	}
}
