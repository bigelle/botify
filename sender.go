package botify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bigelle/botify/internal/reused"
)

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

func (r *APIResponse) BindResult(dest any) error {
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
type RequestSender interface {
	Send(obj APIMethod) (*APIResponse, error)
	SendWithContext(ctx context.Context, obj APIMethod) (*APIResponse, error)
	SendRaw(method string, obj any) (*APIResponse, error)
	SendRawWithContext(ctx context.Context, method string, obj any) (*APIResponse, error)
}

type DefaultRequestSender struct {
	Client   *http.Client
	APIToken string
	APIHost  string
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

	return s.send(ctx, obj.Method(), payload, obj.ContentType())
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

	var payload *bytes.Buffer
	if obj != nil {
		payload = reused.Buf()
		defer reused.PutBuf(payload)

		if err = json.NewEncoder(payload).Encode(obj); err != nil {
			return nil, fmt.Errorf("encoding request payload: %w", err)
		}
	}

	return s.send(ctx, method, payload, "application/json")
}

func (s *DefaultRequestSender) send(ctx context.Context, method string, payload io.Reader, contentType string) (apiResp *APIResponse, err error) {
	if s.Client == nil {
		s.Client = http.DefaultClient
	}

	var req *http.Request
	var resp *http.Response

	reqURL := fmt.Sprintf("%sbot%s/%s", s.APIHost, s.APIToken, method)
	forDebugURL := fmt.Sprintf("%sbot<API token with length = %d>/%s", s.APIHost, len(s.APIToken), method)

	req, err = http.NewRequestWithContext(ctx, "POST", reqURL, payload)
	if err != nil {
		return nil, fmt.Errorf("creating request with URL %s: %w", forDebugURL, err)
	}

	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	resp, err = s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request with URL %s: %w", forDebugURL, err)
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("reading API response: %w", err)
	}

	return apiResp, nil
}
