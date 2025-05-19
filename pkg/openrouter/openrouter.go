package openrouter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/rs/zerolog/log"
)

func OpenRouterCall[T any](ctx context.Context, apiKey string, req *http.Request, err error, allowedStatus ...int) Response[T] {
	if err != nil {
		return Response[T]{Err: fmt.Errorf("error creating request: %w", err)}
	}

	AddDefaultHeaders(apiKey, req)
	resp, err := http.DefaultClient.Do(req)
	return FromResponse[T](ctx, resp, err, allowedStatus...)
}

func NewRequest(ctx context.Context, method string, endpoint string, body any) (*http.Request, error) {
	requestJSON, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request for %s: %w", endpoint, err)
	}

	url := fmt.Sprintf("https://openrouter.ai/api/v1/%s", endpoint)

	log.Debug().Str("url", url).Str("method", method).Bytes("body", requestJSON).Msg("sending to openrouter")

	return http.NewRequestWithContext(ctx, method, url, strings.NewReader(string(requestJSON)))
}

func AddDefaultHeaders(APIKey string, req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+APIKey)
	req.Header.Set("HTTP-Referer", "https://github.com/kklipsch/billy-bot")
	req.Header.Set("X-Title", "Billy Bot")
}

type Response[T any] struct {
	Body   string
	Err    error
	Result T
}

func FromResponse[T any](ctx context.Context, resp *http.Response, err error, allowedStatus ...int) (oresp Response[T]) {
	oresp = Response[T]{}

	if err != nil {
		oresp.Err = fmt.Errorf("error sending request: %w", err)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		oresp.Err = fmt.Errorf("error reading response body: %w", err)
		return
	}
	defer resp.Body.Close()

	strbody := strings.TrimSpace(string(body))
	log.Debug().Int("status", resp.StatusCode).Str("body", strbody).Msg("response from openrouter")
	oresp.Body = strbody

	if !slices.Contains(allowedStatus, resp.StatusCode) {
		oresp.Err = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		return
	}

	if err = json.Unmarshal(body, &oresp.Result); err != nil {
		oresp.Err = fmt.Errorf("error unmarshaling response: %w", err)
	}

	return oresp
}
