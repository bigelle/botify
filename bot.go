package botify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
)

type UpdateSupplier interface {
	GetUpdates(context.Context, chan<- Update) error
}

type RequestSender interface {
	Send(obj APIMethod) (*APIResponse, error)
	SendWithContext(ctx context.Context, obj APIMethod) (*APIResponse, error)
	SendRaw(method string, obj any) (*APIResponse, error)
	SendRawWithContext(ctx context.Context, method string, obj any) (*APIResponse, error)
}

type Bot struct {
	// configurable
	handlers   map[UpdateType]HandlerFunc
	sender     RequestSender
	supplier   UpdateSupplier
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
			ctx := Context{
				bot:     b,
				updType: upd.UpdateType(),
				upd:     &upd,
			}

			all, ok := b.handlers[UpdateTypeAll]
			if ok {
				all(ctx)
			}

			exact, ok := b.handlers[ctx.updType]
			if ok {
				exact(ctx)
			}
		}
	}
}

func (b *Bot) Handle(t UpdateType, handler HandlerFunc) {
	b.handlers[t] = handler
}

func (b *Bot) RequestSender() RequestSender {
	return b.sender
}

func (b *Bot) WithRequestSender(sender RequestSender) {
	b.sender = sender
	lps, ok := b.supplier.(*LongPollingSupplier)
	if ok {
		lps.Sender = sender
		b.supplier = lps
	}
}

func (b *Bot) UpdateSupplier() UpdateSupplier {
	return b.supplier
}

func (b *Bot) WithUpdateSupplier(supp UpdateSupplier) {
	b.supplier = supp
}

func (b *Bot) Serve() error {
	ctx, cancel := context.WithCancel(context.Background())
	b.ctx = ctx
	b.cancel = cancel

	b.chUpdate = make(chan Update, b.bufSize)

	go b.work()
	return b.supplier.GetUpdates(b.ctx, b.chUpdate)
}

func DefaultLongPollingBot(token string) *Bot {
	sender := DefaultRequestSender{
		Client:   http.DefaultClient,
		APIToken: token,
		APIHost:  "https://api.telegram.org/",
		UsePOST:  false,
	}
	bot := Bot{
		handlers: make(map[UpdateType]HandlerFunc),
		sender:   &sender,
		supplier: &LongPollingSupplier{
			Sender: &sender,
			PollingParams: GetUpdates{
				Offset:         0,
				Timeout:        30,
				Limit:          100,
				AllowedUpdates: &[]string{},
			},
		},
		bufSize:    0,
		workerPool: runtime.NumCPU(),
	}
	return &bot
}

type LongPollingSupplier struct {
	Sender RequestSender

	PollingParams GetUpdates
}

func (e *LongPollingSupplier) GetUpdates(ctx context.Context, chUpdate chan<- Update) error {
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

func (e *LongPollingSupplier) Send(obj APIMethod) (*APIResponse, error) {
	return e.Sender.Send(obj)
}

func (e *LongPollingSupplier) SendWithContext(ctx context.Context, obj APIMethod) (*APIResponse, error) {
	return e.Sender.SendWithContext(ctx, obj)
}

func (e *LongPollingSupplier) SendRaw(method string, obj any) (*APIResponse, error) {
	return e.Sender.SendRaw(method, obj)
}

func (e *LongPollingSupplier) SendRawWithContext(ctx context.Context, method string, obj any) (*APIResponse, error) {
	return e.Sender.SendRawWithContext(ctx, method, obj)
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
	return s.SendWithContext(context.Background(), obj)
}

func (s *DefaultRequestSender) SendWithContext(ctx context.Context, obj APIMethod) (apiResp *APIResponse, err error) {
	if obj == nil {
		return nil, fmt.Errorf("obj can't be empty")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	var payload io.Reader
	payload, err = obj.Payload()
	if err != nil {
		return nil, fmt.Errorf("forming request payload: %w", err)
	}

	return s.send(ctx, obj.Method(), payload)
}

func (s *DefaultRequestSender) SendRaw(method string, obj any) (apiResp *APIResponse, err error) {
	return s.SendRawWithContext(context.Background(), method, obj)
}

func (s *DefaultRequestSender) SendRawWithContext(ctx context.Context, method string, obj any) (apiResp *APIResponse, err error) {
	if method == "" {
		return nil, fmt.Errorf("method can't be empty")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	var payload io.ReadWriter
	if obj != nil {
		payload = &bytes.Buffer{}

		if err = json.NewEncoder(payload).Encode(obj); err != nil {
			return nil, fmt.Errorf("encoding request payload: %w", err)
		}
	}

	return s.send(ctx, method, payload)
}

func (s *DefaultRequestSender) send(ctx context.Context, method string, payload io.Reader) (apiResp *APIResponse, err error) {
	var req *http.Request
	var resp *http.Response

	m := "GET"
	if s.UsePOST {
		m = "POST"
	}

	reqURL := fmt.Sprintf("%sbot%s/%s", s.APIHost, s.APIToken, method)

	req, err = http.NewRequestWithContext(ctx, m, reqURL, payload)
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

type ChatMigratedError int

func (e ChatMigratedError) Error() string {
	return fmt.Sprintf("the group has been migrated to the supergroup with the identiefier %d", e)
}

type TooManyRequestsError int

func (e TooManyRequestsError) Error() string {
	return fmt.Sprintf("too many requests; retry after %d seconds", e)
}

type BadRequestError string

func (e BadRequestError) Error() string {
	return string(e)
}

type APIResponse struct {
	Ok          bool                `json:"ok"`
	Description string              `json:"description"`
	Result      json.RawMessage     `json:"result"`
	Parameters  *ResponseParameters `json:"parameters"`
}

func (r *APIResponse) Bind(dest any) error {
	return json.NewDecoder(bytes.NewReader(r.Result)).Decode(dest)
}

func (r *APIResponse) IsSuccessful() bool {
	return r.Ok
}

func (r *APIResponse) GetError() error {
	if r.IsSuccessful() {
		return nil
	}

	if r.Parameters != nil {
		params := r.Parameters
		if params.MigrateToChatID != nil {
			return ChatMigratedError(*params.MigrateToChatID)
		}
		if params.RetryAfter != nil {
			return TooManyRequestsError(*params.RetryAfter)
		}
	}

	return BadRequestError(r.Description)
}

type ResponseParameters struct {
	MigrateToChatID *int `json:"migrate_to_chat_id"`
	RetryAfter      *int `json:"retry_after"`
}

type WebhookInfo struct {
	URL                          string    `json:"url"`
	HasCustomCertificate         *bool     `json:"has_custom_certificate"`
	PendingUpdateCount           *int      `json:"pending_update_count"`
	IPAddress                    *string   `json:"ip_address"`
	LastErrorDate                *int      `json:"last_error_date"`
	LastErrorMessage             *string   `json:"last_error_message"`
	LastSynchronizationErrorDate *int      `json:"last_synchronization_error_date"`
	MaxConnections               *int      `json:"max_connections"`
	AllowedUpdates               *[]string `json:"allowed_updates"`
}
