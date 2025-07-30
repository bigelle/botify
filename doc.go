/*
Package botify provides a flexible Go library for creating Telegram bots with
customizable components and convenient update handling.

# Overview

Botify is designed with a modular architecture that allows developers to replace
components (RequestSender, UpdateReceiver, Logger) to fit specific needs. The
library features typed update handling with handlers for each update type,
inspired by http.ServeMux, and provides convenient command handling with
localization and scope support.

# Features

The botify package includes the following core features:

- Modular architecture with replaceable components
- Typed update handling for all Telegram update types
- Convenient command handling with localization support
- Automatic command menu creation and update filtering
- Support for both Long Polling and Webhooks
- Built-in logging using logr.Logger
- Complete set of Telegram Bot API objects and methods
- Context-based request handling

# Architecture

The package is structured around several core components:

- Bot: Main bot structure with configurable components
- UpdateReceiver: Interface for receiving updates (LongPolling/Webhook)
- RequestSender: Interface for sending requests to Telegram Bot API
- Context: Custom context providing direct access to RequestSender
- Logger: Logging interface using logr.Logger

# Requirements

- Go 1.24+

# Basic Usage

	import "github.com/bigelle/botify"

	func main() {
		bot := &botify.Bot{
			Token: "YOUR_BOT_TOKEN",
		}

		// Handle commands
		bot.HandleCommand("/start", "Start the bot", handleStart)
		bot.HandleCommand("help", "Get help", handleHelp)

		// Handle messages
		bot.Handle(botify.UpdateTypeMessage, handleMessage)

		// Start the bot
		if err := bot.Start(); err != nil {
			panic(err)
		}
	}

# Configuration

The Bot structure allows for flexible configuration:

	bot := &botify.Bot{
		Token:          "YOUR_BOT_TOKEN",        // Required
		UpdateReceiver: customUpdateReceiver,    // Optional
		RequestSender:  customRequestSender,     // Optional
		Logger:         customLogger,            // Optional
	}

Default components are used when not specified:
- UpdateReceiver: LongPolling
- RequestSender: TgBotAPIRequestSender
- Logger: logr.Discard()

# Update Handling

Register typed handlers for specific update types:

	// Handle messages
	bot.Handle(botify.UpdateTypeMessage, func(ctx *botify.Context) error {
		message := ctx.MustGetMessage()
		// Handle message logic
		return nil
	})

	// Handle callback queries
	bot.Handle(botify.UpdateTypeCallbackQuery, func(ctx *botify.Context) error {
		// Handle callback query logic
		return nil
	})

Alternative string-based registration is also supported:

	bot.Handle("message", messageHandler)
	bot.Handle("inline_query", inlineQueryHandler)

# Command Handling

Basic command registration:

	bot.HandleCommand("/start", "Start the bot", func(ctx *botify.Context) error {
		// Command handler logic
		return nil
	})

Commands with localization support:

	locales := map[string]string{
		"en": "Get help",
		"ru": "Получить помощь",
	}
	bot.HandleCommandWithLocales("help", locales, helpHandler, 
		botify.BotCommandScopeAllPrivateChats)

# Update Receivers

## Long Polling (Default)

	longPolling := &botify.LongPolling{
		// Configuration options
	}
	bot := &botify.Bot{
		Token:          "YOUR_BOT_TOKEN",
		UpdateReceiver: longPolling,
	}

## Webhook

	webhook := &botify.Webhook{
		// Configuration options
	}
	bot := &botify.Bot{
		Token:          "YOUR_BOT_TOKEN",
		UpdateReceiver: webhook,
	}

Note: Webhook requires HTTPS and valid SSL certificate.

# Request Sending

Custom RequestSender implementation:

	type CustomSender struct{}

	func (s *CustomSender) Send(obj APIMethod) (*APIResponse, error) {
		// Implementation
	}

	func (s *CustomSender) SendRaw(method string, obj any) (*APIResponse, error) {
		// Implementation
	}

	func (s *CustomSender) SendWithContext(ctx context.Context, obj APIMethod) (*APIResponse, error) {
		// Implementation
	}

	func (s *CustomSender) SendRawWithContext(ctx context.Context, method string, obj any) (*APIResponse, error) {
		// Implementation
	}

# Context Usage

The Context provides direct access to RequestSender methods:

	func myHandler(ctx *botify.Context) error {
		// Send API request
		sendMessage := &botify.SendMessage{
			ChatID: "@username",
			Text:   "Hello!",
		}
		
		response, err := ctx.SendRequest(sendMessage)
		if err != nil {
			return err
		}

		// Access update data
		message := ctx.MustGetMessage() // Non-nil for message updates
		// message := ctx.GetMessage()  // May be nil, use when uncertain

		return nil
	}

# API Methods

Complete set of Telegram Bot API objects and methods:

	sendMessage := botify.SendMessage{
		ChatID: "@username",
		Text:   "Hello World!",
	}

	resp, err := ctx.SendRequest(&sendMessage)
	if err != nil {
		// Handle error
	}

# Logging

Integration with logr.Logger:

	import "github.com/go-logr/logr"

	bot := &botify.Bot{
		Token:  "YOUR_BOT_TOKEN",
		Logger: myLogger, // logr.Logger implementation
	}

# Error Handling

All handlers should return an error. The bot will handle logging and
appropriate error responses:

	func myHandler(ctx *botify.Context) error {
		if someCondition {
			return fmt.Errorf("handling failed: %w", err)
		}
		return nil
	}

# Best Practices

- Always handle errors appropriately in your handlers

- Use MustGet* methods only when you're certain about the update type

- Implement custom components when default behavior doesn't fit your needs

- Use localization for commands in multi-language bots

- Consider using webhook for production deployments with high load

For more information and examples, see the package documentation and
Telegram Bot API reference at https://core.telegram.org/bots/api.
*/
package botify
