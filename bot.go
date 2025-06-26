package botify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
)

type BotEngine interface {
	ProgressUpdate(chan<- Context) error
}

type LongPollingEngine struct {
	// TODO: getUpdates params
}

func (e *LongPollingEngine) ProgressUpdate(chCtx chan<- Context) error {
	// TODO:
	return nil
}

type WebhookEngine struct {
	// TODO: setWebhook params
}

func (e *WebhookEngine) ProgressUpdate(chCtx chan<- Context) error {
	// TODO
	return nil
}

type Bot struct {
	// configurable
	token      string
	handlers   map[UpdateType]HandlerFunc
	engine     BotEngine
	bufSize    int
	workerPool int

	// runtime
	chCtx      chan Context
	isOnline   bool
	hasWebhook bool
}

func NewBot(token string, engine BotEngine) *Bot {
	// NOTE: should i return err? what causes it?
	return &Bot{
		token:      token,
		handlers:   make(map[UpdateType]HandlerFunc), // FIXME: use a map full of default empty handlers
		engine:     engine,
		bufSize:    0,
		workerPool: runtime.NumCPU(),
	}
}

func (b *Bot) Serve() error {
	wh, err := b.getWebhookInfo()
	if err != nil {
		return fmt.Errorf("requesting info about webhook: %w", err)
	}
	if wh.URL != "" {
		b.hasWebhook = true
	}

	_, ok := b.engine.(*LongPollingEngine)
	if ok && b.hasWebhook {
		return fmt.Errorf("can't use long polling bot when webhook is set; call for deleteWebhook before using long polling bot")
	}

	return b.engine.ProgressUpdate(b.chCtx)
}

func (b *Bot) getWebhookInfo() (WebhookInfo, error) {
	reqURL := fmt.Sprintf("https://api.telegram.org/bot%s/getWebhookInfo", b.token)
	resp, err := http.Get(reqURL)
	if err != nil {
		return WebhookInfo{}, fmt.Errorf("requesting for webhook info: %w", err)
	}
	defer resp.Body.Close()

	var response APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return WebhookInfo{}, fmt.Errorf("reading API response: %w", err)
	}

	if !response.Ok {
		return WebhookInfo{}, fmt.Errorf("error from API: %s", response.Description)
	}

	return response.Result.(WebhookInfo), nil
}

type APIResponse struct {
	Ok          bool
	Description string
	Result      any
	// TODO:  response parameters
}

type WebhookInfo struct {
	URL                          string
	HasCustomCertificate         bool
	PendingUpdateCount           int
	IPAddress                    string
	LastErrorDate                int
	LastErrorMessage             string
	LastSynchronizationErrorDate int
	MaxConnections               int
	AllowedUpdates               []string
}
