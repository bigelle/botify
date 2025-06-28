package botify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
)

type BotEngine interface {
	UpdateSupplier
	RequestSender
}

type UpdateSupplier interface {
	GetUpdates(context.Context, chan<- Update) error
}

type RequestSender interface {
	Send(obj APIMethod) (*APIResponse, error)
	SendRaw(method string, obj any) (*APIResponse, error)
}

type LongPollingEngine struct {
	Sender RequestSender

	PollingParams GetUpdates
}

func (e *LongPollingEngine) GetUpdates(ctx context.Context, chUpdate chan<- Update) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			get := GetUpdates(e.PollingParams)

			resp, err := e.Sender.Send(&get)
			if err != nil {
				return fmt.Errorf("polling for updates: %w", err)
			}
			if !resp.Ok {
				return fmt.Errorf("error from API: %s", resp.Description)
			}

			var upds []Update
			resp.Bind(&upds)

			for _, upd := range upds {
				chUpdate <- upd
				e.PollingParams.Offset = upd.UpdateID + 1
			}
		}
	}
}

func (e *LongPollingEngine) Send(obj APIMethod) (*APIResponse, error) {
	return e.Sender.Send(obj)
}

func (e *LongPollingEngine) SendRaw(method string, obj any) (*APIResponse, error) {
	return e.Sender.SendRaw(method, obj)
}

type WebhookEngine struct {
	// TODO: setWebhook params
}

func (e *WebhookEngine) GetUpdates(ctx context.Context, chUpdate chan<- Update) error {
	// TODO:
	return nil
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

	payload, err := obj.Payload()
	if err != nil {
		return nil, fmt.Errorf("forming request payload: %w", err)
	}

	req, err = http.NewRequest(m, reqURL, payload)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", obj.ContentType())

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
	handlers   map[UpdateType][]HandlerFunc
	engine     BotEngine
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
			ctx := Context{bot: b}
			if upd.Message != nil {
				ctx.updType = UpdateTypeMessage
				ctx.updObj = upd.Message

				handlers := b.handlers[ctx.updType]
				for _, h := range handlers {
					h(ctx)
				}
			}

		}
	}
}

func (b *Bot) Handle(t UpdateType, handler ...HandlerFunc) {
	b.handlers[t] = append(b.handlers[t], handler...)
}

func (b *Bot) Serve() error {
	ctx, cancel := context.WithCancel(context.Background())
	b.ctx = ctx
	b.cancel = cancel

	b.chUpdate = make(chan Update, b.bufSize)

	go b.work()
	return b.engine.GetUpdates(b.ctx, b.chUpdate)
}

type LongPollingBuilder struct {
	token  string
	engine LongPollingEngine
}

func NewLongPollingBuilder(token string) *LongPollingBuilder {
	return &LongPollingBuilder{
		token: token,
		engine: LongPollingEngine{
			Sender: &DefaultRequestSender{
				Client:   http.DefaultClient,
				APIToken: token,
				APIHost:  "https://api.telegram.org/",
			},
			PollingParams: struct {
				Offset         int       "json:\"offset\""
				Limit          int       "json:\"limit\""
				Timeout        int       "json:\"timeout\""
				AllowedUpdates *[]string "json:\"allowed_updates\""
			}{
				Offset:         0,
				Limit:          100,
				Timeout:        30,
				AllowedUpdates: &[]string{},
			},
		},
	}
}

func (b *LongPollingBuilder) Build() *Bot {
	return &Bot{
		handlers:   make(map[UpdateType][]HandlerFunc),
		engine:     &b.engine,
		bufSize:    0,
		workerPool: runtime.NumCPU(),
	}
}

type APIResponse struct {
	Ok          bool
	Description string
	Result      json.RawMessage
	// TODO:  response parameters
}

func (r *APIResponse) Bind(dest any) error {
	return json.NewDecoder(bytes.NewReader(r.Result)).Decode(dest)
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
