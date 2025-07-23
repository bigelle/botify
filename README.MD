# Botify

A flexible Go library for creating Telegram bots with customizable components and convenient update handling.

## Features

- 🔧 **Modular architecture** - replace components (RequestSender, UpdateReceiver, Logger) to fit your needs
- 🎯 **Typed update handling** - handlers for each update type, inspired by `http.ServeMux`
- 📝 **Convenient command handling** - with localization and scope support
- 🤖 **Automatic setup** - automatic command menu creation and update filtering
- 📡 **Long Polling and Webhooks** - choose your preferred method for receiving updates
- 📊 **Built-in logging** - using `logr.Logger`
- 🚀 **Ready-to-use API methods** - complete set of Telegram Bot API objects and methods

## Requirements

- Go 1.24+

## Installation

```bash
go get github.com/bigelle/botify
```

## Quick Start

```go
package main

import (
    "github.com/bigelle/botify"
)

func main() {
    bot := &botify.Bot{
        Token: "YOUR_BOT_TOKEN",
    }

    // Handle commands
    bot.HandleCommand("/start", "Start the bot", handleStart)
    // Equivalent for "/help":
    bot.HandleCommand("help", "Get help", handleHelp)

    // Handle messages
    bot.Handle(botify.UpdateTypeMessage, handleMessage)
    // Or:
    // bot.Handle("message", handleMessage)

    // Start the bot
    if err := bot.Start(); err != nil {
        panic(err)
    }
}

func handleStart(ctx *botify.Context) error {
    // TODO: Implementation of /start command handler
    return nil
}

func handleHelp(ctx *botify.Context) error {
    // TODO: Implementation of /help command handler
    return nil
}

func handleMessage(ctx *botify.Context) error {
    // TODO: Handle incoming messages
    return nil
}
```

## Components

### Bot

The main bot structure. Created without a constructor in `http.Client` style:

```go
bot := &botify.Bot{
    Token:          "YOUR_BOT_TOKEN",        // Required field
    UpdateReceiver: customUpdateReceiver,    // Optional
    RequestSender:  customRequestSender,     // Optional  
    Logger:         customLogger,            // Optional
}
```

> [!NOTE]
> If components are not specified, default values are used:
> - `UpdateReceiver`: `LongPollingBot`
> - `RequestSender`: `TgBotAPIRequestSender`
> - `Logger`: `logr.Discard()`

> [!WARNING]
> Using the bot or sender without a token will cause a panic.

### UpdateReceiver

Interface for receiving updates from Telegram:

#### LongPolling (default)

```go
longPolling := &botify.LongPolling{
    // TODO: Long Polling configuration
}

bot := &botify.Bot{
    Token:          "YOUR_BOT_TOKEN",
    UpdateReceiver: longPolling,
}
```

#### Webhook

```go
webhook := &botify.Webhook{
    // TODO: Webhook configuration
}

bot := &botify.Bot{
    Token:          "YOUR_BOT_TOKEN", 
    UpdateReceiver: webhook,
}
```

> [!WARNING]
> When using Webhook, ensure your server is accessible via HTTPS and has a valid SSL certificate.

### RequestSender

Interface for sending requests to Telegram Bot API:

```go
type MyCustomSender struct {
    // Whatever you need:
}

func (s *MyCustomSender) Send(obj APIMethod) (*APIResponse, error) {
    //TODO: 
    return nil, nil
}

func (s *MyCustomSender) SendRaw(method string, obj any) (*APIResponse, error) {
    //TODO: 
    return nil, nil
}

func (s *MyCustomSender) SendWithContext(ctx context.Context, obj APIMethod) (*APIResponse, error) {
    //TODO: 
    return nil, nil
}

func (s *MyCustomSender) SendRawWithContext(ctx context.Context, method string, obj any) (*APIResponse, error) {
    //TODO: 
    return nil, nil
}

customSender := &MyCustomRequestSender{}

bot := &botify.Bot{
    Token:         "YOUR_BOT_TOKEN",
    RequestSender: customSender,
}
```

### Logger

Uses `logr.Logger` for logging:

```go
import "github.com/go-logr/logr"

bot := &botify.Bot{
    Token:  "YOUR_BOT_TOKEN",
    Logger: myLogger, // logr.Logger
}
```

## Update Handling

### Typed Handlers

Register handlers for specific update types:

```go
// Messages
bot.Handle(botify.UpdateTypeMessage, func(ctx *botify.Context) error {
    // TODO: Handle messages
    return nil
})

// Inline queries
bot.Handle(botify.UpdateTypeInlineQuery, func(ctx *botify.Context) error {
    // TODO: Handle inline queries
    return nil
})

// Callback queries
bot.Handle(botify.UpdateTypeCallbackQuery, func(ctx *botify.Context) error {
    // TODO: Handle callback queries
    return nil
})
```

> [!NOTE]
> Instead of `botify.UpdateType...` you can use `"message"`, `"inline_query"`, etc.
> See [Update](https://core.telegram.org/bots/api#update) for a complete list of available update types.

### Command Handling

#### Basic Command Handling

```go
bot.HandleCommand("/start", "Start the bot", func(ctx *botify.Context) error {
    // TODO: Handle /start command
    return nil
})

// Equivalent for "/settings":
bot.HandleCommand("settings", "Settings", func(ctx *botify.Context) error {
    // TODO: Handle /settings command
    return nil
}, botify.BotCommandScopeDefault)
```

#### Commands with Localization

```go
locales := map[string]string{
    "en": "Get help",
    "ru": "Получить помощь",
}

bot.HandleCommandWithLocales("help", locales, func(ctx *botify.Context) error {
    // TODO: Handle /help command
    return nil
}, botify.BotCommandScopeAllPrivateChats)
```

> [!TIP]
> The bot will automatically set up the command menu and filter only the necessary update types when starting.

## Context

The custom context provides direct access to the RequestSender:

```go
func myHandler(ctx *botify.Context) error {
    // Direct access to RequestSender methods
    response, err := ctx.SendRequest(request) // equivalent to Sender.Send()
    if err != nil {
        return err
    }
    
    // Get update data
    message := ctx.MustGetMessage() // always non-nil, if this is a message update
    // message := ctx.GetMessage() // may be nil, use it if you're not sure
    
    return nil
}
```

> [!NOTE]  
> `Context` allows direct calls to the RequestSender used by the bot through `ctx.SendRequest()`, eliminating the need to store references to it.

## API Methods

The library provides a ready-to-use set of Telegram Bot API objects and methods:

```go
// Example: sending a message

sendMessage := botify.SendMessage{
    ChatID: "@username" // or fmt.Sprint(123456) for numeric IDs
    Text: "HIIII"
}

resp, err := ctx.SendRequest(&sendMessage) 
// Or using request sender:
// resp, err := sender.Send(&sendMessage)

// handle the response...
```

## License

MIT License. See [LICENSE](LICENSE) for details.

---

> [!TIP]
> For additional information, see the [documentation](https://pkg.go.dev/github.com/bigelle/botify).
> Also see the [Full API Reference for Developers](https://core.telegram.org/bots/api).
