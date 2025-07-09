package botify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type UpdateSupplier interface {
	GetUpdates(context.Context, chan<- Update) error
}

type LongPollingSupplier struct {
	Sender RequestSender

	Offset  int
	Limit   int
	Timeout int
	// TODO: it should be filled by bot accoriding to the list of registered handlers
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

type WebhookSupplier struct {
	url                string
	certificate        InputFile
	ipAddress          string
	maxConnections     int
	allowedUpdates     *[]string
	dropPendingUpdates bool
	secretToken        string

	addr string
}

func NewWebhookSupplier(url string) *WebhookSupplier {
	return &WebhookSupplier{
		url: url,

		addr: ":8080",
	}
}

func (ws *WebhookSupplier) GetUpdates(ctx context.Context, chUpdate chan<- Update) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/webhook", ws.handlerFunc(chUpdate))

	server := &http.Server{
		Addr:    ws.addr,
		Handler: mux,
	}

	serverErr := make(chan error, 1)

	go func() {
		log.Printf("Listening and serving on %s", ws.addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Printf("Stopping server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("error while stopping server: %v", err)
			return err
		}

		log.Print("Server is stopped")
		return ctx.Err()

	case err := <-serverErr:
		log.Printf("Server error: %v", err)
		return err
	}
}

func (ws *WebhookSupplier) handlerFunc(chUpdate chan<- Update) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		b, _ := io.ReadAll(r.Body)

		dec := json.NewDecoder(bytes.NewReader(b))
		// dec.DisallowUnknownFields() //FIXME:

		var upd Update
		if err := dec.Decode(&upd); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("parsing body: %s", err.Error())
			fmt.Println(string(b))
			return
		}

		chUpdate <- upd

		w.WriteHeader(http.StatusOK)
	}
}
