package botify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
)

type BotEngine interface {
	GetUpdates(chan<- Update) error
}

type LongPollingEngine struct {
	// TODO: getUpdates params
}

func (e *LongPollingEngine) GetUpdates(chCtx chan<- Update) error {
	// TODO:
	return nil
}

type WebhookEngine struct {
	// TODO: setWebhook params
}

func (e *WebhookEngine) ProgressUpdate(chCtx chan<- Update) error {
	// TODO
	return nil
}

type RequestSender interface {
	Send(obj APIMethod) (*APIResponse, error)
	SendRaw(method string, obj any) (*APIResponse, error)
}

type DefaultRequestSender struct {
	Client   *http.Client
	APIToken string
	APIHost  string
	UsePOST  bool
}

func (s *DefaultRequestSender) Send(obj APIMethod) (apiResp *APIResponse, err error) {
	var req *http.Request
	var resp *http.Response

	m := "GET"
	if s.UsePOST {
		m = "POST"
	}

	reqURL := fmt.Sprintf("%sbot%s/%s", s.APIHost, s.APIToken, obj.Method())

	req, err = http.NewRequest(m, reqURL, obj.Payload())
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err = s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("reading API response: %w", err)
	}

	return apiResp, nil
}

func (s *DefaultRequestSender) SendRaw(method string, obj any) (apiResp *APIResponse, err error) {
	var req *http.Request
	var resp *http.Response
	payload := &bytes.Buffer{}

	if obj != nil {
		if err = json.NewEncoder(payload).Encode(obj); err != nil {
			return nil, fmt.Errorf("encoding request payload: %w", err)
		}
	}

	m := "GET"
	if s.UsePOST {
		m = "POST"
	}

	reqURL := fmt.Sprintf("%sbot%s/%s", s.APIHost, s.APIToken, method)

	req, err = http.NewRequest(m, reqURL, payload)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err = s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("reading API response: %w", err)
	}

	return apiResp, nil
}

type Bot struct {
	// configurable
	handlers      map[UpdateType][]HandlerFunc
	engine        BotEngine
	requestSender RequestSender
	bufSize       int
	workerPool    int

	// runtime
	chUpdate      chan Update
	isOnline   bool
	hasWebhook bool
}

func NewBot(token string, engine BotEngine) *Bot {
	// NOTE: should i return err? what causes it?
	return &Bot{
		handlers:      make(map[UpdateType][]HandlerFunc), // FIXME: use a map full of default empty handlers
		engine:        engine,
		requestSender: &DefaultRequestSender{APIToken: token},
		bufSize:       0,
		workerPool:    runtime.NumCPU(),
	}
}

func (b *Bot) Handle(t UpdateType, handler ...HandlerFunc) {
	b.handlers[t] = append(b.handlers[t], handler...)
}

func (b *Bot) Serve() error {
	_, ok := b.engine.(*LongPollingEngine)
	if ok && b.hasWebhook {
		return fmt.Errorf("can't use long polling bot when webhook is set; call for deleteWebhook before using long polling bot")
	}

	return b.engine.GetUpdates(b.chUpdate)
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
