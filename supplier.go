package botify

import (
	"context"
	"fmt"
	"time"
)

type UpdateSupplier interface {
	GetUpdates(context.Context, chan<- Update) error
}

type LongPollingSupplier struct {
	Sender RequestSender

	Offset         int
	Limit          int
	Timeout        int
	AllowedUpdates *[]string
}

func (e *LongPollingSupplier) GetUpdates(ctx context.Context, chUpdate chan<- Update) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			get := GetUpdates{
				Offset:         e.Offset,
				Limit:          e.Limit,
				Timeout:        e.Timeout,
				AllowedUpdates: e.AllowedUpdates,
			}

			resp, err := e.Sender.SendWithContext(ctx, &get)
			if err != nil {
				return fmt.Errorf("polling for updates: %w", err)
			}
			if !resp.Ok {
				return resp.GetError()
			}

			var upds []Update
			resp.BindResult(&upds)

			// to avoid any rate limits we are sleeping when there's no activity
			if len(upds) == 0 {
				time.Sleep(1 * time.Second)
				continue
			}

			for _, upd := range upds {
				chUpdate <- upd
				e.Offset = upd.UpdateID + 1
			}
		}
	}
}

func (e *LongPollingSupplier) Send(obj APIMethod) (*APIResponse, error) {
	if e.Sender == nil {
		return nil, fmt.Errorf("request sender is empty")
	}
	return e.Sender.Send(obj)
}

func (e *LongPollingSupplier) SendWithContext(ctx context.Context, obj APIMethod) (*APIResponse, error) {
	if e.Sender == nil {
		return nil, fmt.Errorf("request sender is empty")
	}
	return e.Sender.SendWithContext(ctx, obj)
}

func (e *LongPollingSupplier) SendRaw(method string, obj any) (*APIResponse, error) {
	if e.Sender == nil {
		return nil, fmt.Errorf("request sender is empty")
	}
	return e.Sender.SendRaw(method, obj)
}

func (e *LongPollingSupplier) SendRawWithContext(ctx context.Context, method string, obj any) (*APIResponse, error) {
	if e.Sender == nil {
		return nil, fmt.Errorf("request sender is empty")
	}
	return e.Sender.SendRawWithContext(ctx, method, obj)
}

type WebhookEngine struct {
	// TODO: setWebhook params
}

func (e *WebhookEngine) GetUpdates(ctx context.Context, chUpdate chan<- Update) error {
	// TODO:
	return nil
}
