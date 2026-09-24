package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/fr33dman/go-template/internal/clients/embeddings"
	"github.com/fr33dman/go-template/internal/core"
	apperrors "github.com/fr33dman/go-template/internal/errors"
)

type Client struct {
	baseURL    *url.URL
	apiKey     string
	httpClient *http.Client
}

type (
	embeddingsRequest struct {
		Model string   `json:"model"`
		Input []string `json:"input"`
	}
	embeddingsResponse struct {
		Data []embeddingData `json:"data"`
	}
	embeddingData struct {
		Index     int       `json:"index"`
		Embedding []float32 `json:"embedding"`
	}
	errorResponse struct {
		Error *apiError `json:"error"`
	}
	apiError struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	}
)

var _ embeddings.Provider = (*Client)(nil)

func NewClient(cfg embeddings.ProviderConfig) (*Client, error) {
	baseURL, err := url.Parse(strings.TrimRight(cfg.BaseURL, "/"))
	if err != nil {
		return nil, fmt.Errorf("parse embeddings provider %q base_url: %w", cfg.Name, err)
	}
	if baseURL.Scheme == "" || baseURL.Host == "" {
		return nil, fmt.Errorf("embeddings provider %q base_url must be absolute", cfg.Name)
	}

	return &Client{
		baseURL: baseURL,
		apiKey:  cfg.APIKey,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}, nil
}

func (p *Client) Generate(
	ctx context.Context,
	model core.Model,
	req embeddings.GenerateRequest,
) (embeddings.Generated, error) {
	input := make([]string, len(req.Input))
	for i, item := range req.Input {
		input[i] = item.Text
	}

	body, err := json.Marshal(embeddingsRequest{
		Model: model.UpstreamModel,
		Input: input,
	})
	if err != nil {
		return embeddings.Generated{}, fmt.Errorf("encode embeddings request: %w", err)
	}

	endpoint := p.baseURL.JoinPath("embeddings")
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint.String(),
		bytes.NewReader(body),
	)
	if err != nil {
		return embeddings.Generated{}, fmt.Errorf("build embeddings request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if p.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	}

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return embeddings.Generated{}, fmt.Errorf(
			"%w: send embeddings request: %w",
			apperrors.ErrUnavailable,
			err,
		)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return embeddings.Generated{}, fmt.Errorf(
			"%w: %w",
			apperrors.ErrUnavailable,
			decodeOpenAIError(resp),
		)
	}

	var decoded embeddingsResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return embeddings.Generated{}, fmt.Errorf("decode embeddings response: %w", err)
	}
	if len(decoded.Data) != len(req.Input) {
		return embeddings.Generated{}, fmt.Errorf(
			"embeddings response item count mismatch: expected %d, got %d",
			len(req.Input),
			len(decoded.Data),
		)
	}

	items := make([]core.Embedding, len(req.Input))
	for _, item := range decoded.Data {
		if item.Index < 0 || item.Index >= len(req.Input) {
			return embeddings.Generated{}, fmt.Errorf(
				"embeddings response index %d is out of range",
				item.Index,
			)
		}
		if len(item.Embedding) != model.Dimensions {
			return embeddings.Generated{}, fmt.Errorf(
				"embedding dimensions mismatch for model %q: expected %d, got %d",
				model.Name,
				model.Dimensions,
				len(item.Embedding),
			)
		}
		items[item.Index] = core.Embedding{
			ID:     req.Input[item.Index].ID,
			Vector: item.Embedding,
		}
	}

	return embeddings.Generated{
		Model:      model.Name,
		Dimensions: model.Dimensions,
		Embeddings: items,
	}, nil
}

func decodeOpenAIError(resp *http.Response) error {
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return fmt.Errorf(
			"embedding provider returned status %d and unreadable body: %w",
			resp.StatusCode,
			err,
		)
	}

	var decoded errorResponse
	if err := json.Unmarshal(
		body,
		&decoded,
	); err == nil && decoded.Error != nil &&
		decoded.Error.Message != "" {
		return fmt.Errorf(
			"embedding provider returned status %d: %s",
			resp.StatusCode,
			decoded.Error.Message,
		)
	}

	message := strings.TrimSpace(string(body))
	if message == "" {
		message = http.StatusText(resp.StatusCode)
	}
	return fmt.Errorf("embedding provider returned status %d: %s", resp.StatusCode, message)
}
