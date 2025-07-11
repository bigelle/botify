package botify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"
)

type UpdateSupplier interface {
	GetUpdates(context.Context, chan<- Update) error
	AllowUpdate(upds ...UpdateType) // maybe it should'nt be a part of the interface?
}

func allowUpdate(list []string, upds ...UpdateType) {
	if list == nil {
		list = []string{}
	}

	list = slices.Grow(list, len(upds))

	for _, upd := range upds {
		if upd == UpdateTypeAll {
			// TODO: add every possible update type
			return
		}

		if !slices.Contains(list, upd.String()) {
			list = append(list, upd.String())
		}
	}
}

type LongPollingSupplier struct {
	Sender RequestSender

	Offset  int
	Limit   int
	Timeout int
	// TODO: it should be filled by bot accoriding to the list of registered handlers
	AllowedUpdates *[]string
}

func (lps *LongPollingSupplier) AllowUpdate(upd ...UpdateType) {
	allowUpdate(*lps.AllowedUpdates, upd...)
}

func (e *LongPollingSupplier) GetUpdates(ctx context.Context, chUpdate chan<- Update) error {
	if e.Sender == nil {
		return fmt.Errorf("long polling bot requires request sender")
	}

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

type WebhookSupplier struct {
	// In format https://example.com
	Domain string
	// Webhook Path.
	Path string
	// Address for the server.
	// Telegram Bot API works only with port 443, 80, 88 and 8443
	Addr string
	// Optional.
	Certificate InputFile
	// Optional.
	IPAddress string
	// Optional.
	MaxConnections int
	// Optional.
	AllowedUpdates *[]string
	// Optional.
	DropPendingUpdates bool
	// Optional.
	SecretToken string
}

func (ws *WebhookSupplier) AllowUpdate(upds ...UpdateType) {
	allowUpdate(*ws.AllowedUpdates, upds...)
}

func (ws *WebhookSupplier) GetUpdates(ctx context.Context, chUpdate chan<- Update) error {
	mux := http.NewServeMux()

	mux.HandleFunc(ws.Path, ws.handlerFunc(chUpdate))

	server := &http.Server{
		Addr:    ws.Addr,
		Handler: mux,
	}

	serverErr := make(chan error, 1)

	go func() {
		log.Printf("Listening and serving on %s", ws.Addr)
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

		if ws.SecretToken != "" {
			t := r.Header.Get("X-Telegram-Bot-Api-Secret-Token")
			if ws.SecretToken != t {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
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

func (ws *WebhookSupplier) HandlePath(path string) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	ws.Path = path
}
