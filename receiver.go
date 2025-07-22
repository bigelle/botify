package botify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/bigelle/botify/internal/reused"
)

// UpdateReceiver is an interface used to receive updates and send them into chUpdate.
// Since sometimes to receive updates you need to send requests,
// to make it less awkward it pairs with the bot to get access to it's sender.
type UpdateReceiver interface {
	ReceiveUpdates(ctx context.Context, chUpdate chan<- Update) error
	PairBot(*Bot)
}

// LongPolling is a long-polling implementation of UpdateReceiver.
// It sends /getUpdates requests, deserializes the response
// and sends the updates to the channel
type LongPolling struct {
	Offset  int
	Limit   int
	Timeout int

	bot *Bot
}

// PairBot satisfies the UpdateReceiver interface
func (lp *LongPolling) PairBot(b *Bot) {
	lp.bot = b
	b.Receiver = lp
}

// ReceiveUpdates sends /getUpdates requests, deserializes the response
// and sends the updates to the chUpdate
func (lp *LongPolling) ReceiveUpdates(ctx context.Context, chUpdate chan<- Update) error {
	if lp.bot.Sender == nil {
		return fmt.Errorf("long polling bot requires request sender")
	}

	allowedUpdates := make([]string, 0, len(lp.bot.updateHandlers))
	for upd := range lp.bot.updateHandlers {
		allowedUpdates = append(allowedUpdates, upd)
	}
	if len(lp.bot.commandHandlers.byCommand) > 0 && !slices.Contains(allowedUpdates, UpdateTypeMessage) {
		allowedUpdates = append(allowedUpdates, UpdateTypeMessage)
	}

	var (
		get  GetUpdates
		resp *APIResponse
		err  error
	)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			get = GetUpdates{
				Offset:         lp.Offset,
				Limit:          lp.Limit,
				Timeout:        lp.Timeout,
				AllowedUpdates: &allowedUpdates,
			}

			resp, err = lp.bot.Sender.SendWithContext(ctx, &get)
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
				lp.Offset = upd.UpdateID + 1
			}
		}
	}
}

// Webhook is a implementation for UpdateReceiver,
// which creates a server on a given ListenAddr
// and sends /setWebhook request with URL "Host:ExposedPort/Path"
type Webhook struct {
	// In format https://example.com
	Host string
	// Webhook Path.
	Path string
	// Will be send to the Telegram Bot API server.
	// Defaults to 443
	ExposedPort string
	// Will be used to run the webhook server
	ListenAddr string
	// Optional.
	Certificate InputFile
	// Optional.
	IPAddress string
	// Optional.
	MaxConnections int
	// Optional.
	DropPendingUpdates bool
	// Optional.
	SecretToken string

	bot *Bot
}

// ReceiveUpdates creates a webhook server which will send every incoming update into chUpdate
func (ws *Webhook) ReceiveUpdates(ctx context.Context, chUpdate chan<- Update) (err error) {
	if ws.bot.Sender == nil {
		return fmt.Errorf("can't set webhook: no request sender")
	}

	allowedUpdates := make([]string, 0, len(ws.bot.updateHandlers))
	for upd := range ws.bot.updateHandlers {
		allowedUpdates = append(allowedUpdates, upd)
	}
	if len(ws.bot.commandHandlers.byCommand) > 0 && !slices.Contains(allowedUpdates, UpdateTypeMessage) {
		allowedUpdates = append(allowedUpdates, UpdateTypeMessage)
	}

	mux := http.NewServeMux()
	mux.HandleFunc(ws.Path, ws.handlerFunc(chUpdate))

	if ws.ListenAddr == "" {
		ws.ListenAddr = ":443"
	}
	server := &http.Server{
		Addr:    ws.ListenAddr,
		Handler: mux,
	}
	serverErr := make(chan error, 1)

	go func() {
		if err = ws.SetWebhook(ctx, allowedUpdates); err != nil {
			serverErr <- fmt.Errorf("setting webhook: %w", err)
		}
	}()

	go func() {
		ws.bot.Logger.Info("webhook server is listening and serving on", "address", ws.ListenAddr, "exposed port", ws.ExposedPort, "path", ws.Path)
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		ws.bot.Logger.Info("stopping webhook server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err = server.Shutdown(shutdownCtx); err != nil {
			ws.bot.Logger.Error(err, "error shutting down webhook server")
			return err
		}

		ws.bot.Logger.Info("webhook server is stopped")
		return ctx.Err()

	case err = <-serverErr:
		ws.bot.Logger.Error(err, "server error")
		return err
	}
}

// PairBot satisfies the UpdateReceiver interface
func (wh *Webhook) PairBot(b *Bot) {
	wh.bot = b
	b.Receiver = wh
}

// WebhookURL returns tthe URL to which Telegram Bot API will send updates.
// If ExposedPort is 443, it will be omitted since it's a default HTTPS port
func (ws *Webhook) WebhookURL() string {
	port := ""
	if ws.ExposedPort != "" && ws.ExposedPort != "443" && ws.ExposedPort != ":443" {
		port = ws.ExposedPort
		if !strings.HasPrefix(port, ":") {
			port = ":" + port
		}
	}
	if !strings.HasPrefix(ws.Path, "/") {
		ws.Path = "/" + ws.Path
	}
	return fmt.Sprintf("%s%s%s", ws.Host, port, ws.Path)
}

// SetWebhook sends /setWebhook request
func (ws *Webhook) SetWebhook(ctx context.Context, allowedUpdates []string) error {
	if ws.bot == nil || ws.bot.Sender == nil {
		return fmt.Errorf("webhook is not paired with bot")
	}
	swh := SetWebhook{
		URL:                ws.WebhookURL(),
		Certificate:        ws.Certificate,
		IPAddress:          ws.IPAddress,
		MaxConnections:     ws.MaxConnections,
		AllowedUpdates:     &allowedUpdates,
		DropPendingUpdates: ws.DropPendingUpdates,
		SecretToken:        ws.SecretToken,
	}
	resp, err := ws.bot.Sender.SendWithContext(ctx, &swh)
	if err != nil {
		return fmt.Errorf("sending setWebhook request: %w", err)
	}
	if err = resp.GetError(); err != nil {
		return fmt.Errorf("setting webhook: %w", err)
	}
	return nil
}

func (ws *Webhook) handlerFunc(chUpdate chan<- Update) http.HandlerFunc {
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

		var err error
		buf := reused.Buf()
		defer reused.PutBuf(buf)

		_, err = io.Copy(buf, r.Body)
		if err != nil {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			ws.bot.Logger.Error(err, "reading request body")
			return
		}

		dec := json.NewDecoder(buf)
		// dec.DisallowUnknownFields() //FIXME:
		var upd Update
		if err = dec.Decode(&upd); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			ws.bot.Logger.Error(err, "parsing request body")
			return
		}

		chUpdate <- upd
		w.WriteHeader(http.StatusOK)
	}
}
