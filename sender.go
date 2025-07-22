package botify

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bigelle/botify/internal/reused"
)

const TelegramBotAPIHost = "https://api.telegram.org"

var ErrNoResult = errors.New("the response has no result")

// ChatMigratedError is an error signalizing that the group has been migrated to the supergroup
// and holding the new group identifier as int
type ChatMigratedError int

func (e ChatMigratedError) Error() string {
	return fmt.Sprintf("the group has been migrated to the supergroup with the identiefier %d", e)
}

// TooManyRequestsError is an error signalizing that you have exceeded the flood control
// and holding the number of seconds left to wait before the request can be repeated
type TooManyRequestsError int

func (e TooManyRequestsError) Error() string {
	return fmt.Sprintf("too many requests; retry after %d seconds", e)
}

// RetryAfter returns the number of seconds left to wait before the request can be repeated
func (e TooManyRequestsError) RetryAfter() time.Duration {
	return time.Second * time.Duration(e)
}

// BadRequestError is an error signalizing that the request is failed
// and holding a human readable error description as string
type BadRequestError string

func (e BadRequestError) Error() string {
	return string(e)
}

// APIResponse is a response from Telegram Bot API
type APIResponse struct {
	// True if success, false otherwise
	Ok bool `json:"ok"`
	// Optional. If ok is true, it has a human-readable description of the result.
	// If ok is false, it explains the error
	Description string `json:"description"`
	// Optional. If ok is true, the result is available here through the BindResult method.
	Result json.RawMessage `json:"result"`
	// Optional. If ok is false, it has the error code
	ErrorCode int `json:"error_code"`
	// Optional. If ok is false, it may have extra information to automatically handle the error
	Parameters *ResponseParameters `json:"parameters"`
}

// BindResult is used to write response result to dest.
// It will return ErrNoResult if the response has no result,
// or any JSON decoding error if something went wrong
func (r *APIResponse) BindResult(dest any) error {
	if len(r.Result) == 0 {
		return ErrNoResult
	}

	buf := reused.Buf()
	defer reused.PutBuf(buf)
	buf.Write(r.Result)

	dec := json.NewDecoder(buf)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dest); err != nil {
		return fmt.Errorf("decoding result field: %w", err)
	}

	return nil
}

// IsSuccessful is a more readable wrapper around `if r.Ok`
func (r *APIResponse) IsSuccessful() bool {
	return r.Ok
}

// GetError returns nil if the request was successfull,
// ChatMigratedError if MigrateToChatId is not nil,
// TooManyRequestsError if RetryAfter is not nil,
// and BadRequestError with the error description if there's nothing more suitable
func (r *APIResponse) GetError() error {
	if r.IsSuccessful() {
		return nil
	}

	if r.Parameters != nil {
		params := r.Parameters
		if params.MigrateToChatID != nil {
			return fmt.Errorf("%d: %w", r.ErrorCode, ChatMigratedError(*params.MigrateToChatID))
		}
		if params.RetryAfter != nil {
			return fmt.Errorf("%d: %w", r.ErrorCode, TooManyRequestsError(*params.RetryAfter))
		}
	}

	return fmt.Errorf("%d: %s", r.ErrorCode, BadRequestError(r.Description))
}

// ResponseParameters helps to automatically handle the error
type ResponseParameters struct {
	MigrateToChatID *int `json:"migrate_to_chat_id"`
	RetryAfter      *int `json:"retry_after"`
}

// RequestSender is a unified way to send requests
type RequestSender interface {
	Send(obj APIMethod) (*APIResponse, error)
	SendWithContext(ctx context.Context, obj APIMethod) (*APIResponse, error)
	SendRaw(method string, obj any) (*APIResponse, error)
	SendRawWithContext(ctx context.Context, method string, obj any) (*APIResponse, error)
}

// TGBotAPIRequestSender is a default request sender.
// Every method returns [APIResponse] and the result of APIResponse.GetError(), no matter if it's successful or not.
// So there's no need to manually check for `if resp.GetError() != nil` after every request.
//
// If the request fails and if the response parameters contains a "retry_after" field,
// it will try to send the request again after n seconds, where n is the value of the "retry_after" field
type TGBotAPIRequestSender struct {
	Client   *http.Client
	APIToken string
	APIHost  string
}

// Send satisfies RequestSender interface
func (s *TGBotAPIRequestSender) Send(obj APIMethod) (apiResp *APIResponse, err error) {
	return s.SendWithContext(context.Background(), obj)
}

// SendWithContext satisfies RequestSender interface
func (s *TGBotAPIRequestSender) SendWithContext(ctx context.Context, obj APIMethod) (apiResp *APIResponse, err error) {
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

// SendRaw satisfies RequestSender interface
func (s *TGBotAPIRequestSender) SendRaw(method string, obj any) (apiResp *APIResponse, err error) {
	return s.SendRawWithContext(context.Background(), method, obj)
}

// SendRawWithContext satisfies RequestSender interface
func (s *TGBotAPIRequestSender) SendRawWithContext(ctx context.Context, method string, obj any) (apiResp *APIResponse, err error) {
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

func (s *TGBotAPIRequestSender) send(ctx context.Context, method string, payload io.Reader, contentType string) (apiResp *APIResponse, err error) {
	if s.Client == nil {
		s.Client = &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
				DisableKeepAlives:   false,
			},
		}
	}
	if s.APIHost == "" {
		s.APIHost = TelegramBotAPIHost
	}

	var req *http.Request
	var resp *http.Response

	reqURL := fmt.Sprintf("%s/bot%s/%s", s.APIHost, s.APIToken, method)
	forDebugURL := fmt.Sprintf("%sbot<API token with length = %d>/%s", s.APIHost, len(s.APIToken), method)

	req, err = http.NewRequestWithContext(ctx, "POST", reqURL, payload)
	if err != nil {
		return nil, fmt.Errorf("creating request with URL %s: %w", forDebugURL, err)
	}

	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	sendRequest := func(req *http.Request) (*APIResponse, error) {
		resp, err = s.Client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("sending request with URL %s: %w", forDebugURL, err)
		}
		defer resp.Body.Close()

		if err = json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			return nil, fmt.Errorf("reading API response: %w", err)
		}

		return apiResp, apiResp.GetError()
	}

	var errRateLimit TooManyRequestsError

	apiResp, err = sendRequest(req)
	if apiResp == nil && errors.As(err, &errRateLimit) {
		time.Sleep(errRateLimit.RetryAfter())

		apiResp, err = sendRequest(req)
	}
	return apiResp, err
}
