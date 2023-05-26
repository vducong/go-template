package httpclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"promotion/pkg/failure"

	gcppropagator "github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"go.opentelemetry.io/otel/propagation"
)

type CreateHTTPRequestDTO struct {
	Method     string
	URL        string
	Body       io.Reader
	AuthHeader *AuthHeader
}

type HTTPResponse[DataType interface{}] struct {
	Status int      `json:"status"`
	Data   DataType `json:"data"`
}

type HTTPClient struct {
	client *http.Client
}

func New() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{},
	}
}

func CreateRequest(ctx context.Context, dto *CreateHTTPRequestDTO) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, dto.Method, dto.URL, dto.Body)
	if err != nil {
		return nil, failure.ErrWithTrace(fmt.Errorf("Failed to create request=%+v: %w", dto, err))
	}
	req.Header.Set("Content-Type", "application/json")
	setAuthHeader(req, dto.AuthHeader)

	propagator := gcppropagator.CloudTraceFormatPropagator{}
	propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))

	return req, nil
}

func setAuthHeader(req *http.Request, authHeader *AuthHeader) {
	switch authHeader.Method {
	case AuthMethodJWT, AuthMethodBearerAPIKey:
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authHeader.Token))
	case AuthMethodAPIKey:
		req.Header.Set("x-api-key", authHeader.Token)
	}
}

func (c *HTTPClient) Do(req *http.Request) ([]byte, error) {
	res, err := c.client.Do(req)
	if err != nil {
		return nil, failure.ErrWithTrace(fmt.Errorf("Failed to send HTTP request to %s: %w", req.URL, err))
	}

	data, err := readResponseBody(res)
	if err != nil {
		return nil, failure.ErrWithTrace(fmt.Errorf("Failed to read response's body from %s: %w", req.URL, err))
	}

	if res.StatusCode != http.StatusOK {
		return nil, failure.ErrWithTrace(fmt.Errorf("Request to %s failed: %+v", req.URL, string(data)))
	}
	return data, nil
}

func readResponseBody(res *http.Response) ([]byte, error) {
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, failure.ErrWithTrace(fmt.Errorf("Failed to read body: %w", err))
	}
	defer res.Body.Close()
	return bodyBytes, nil
}

func ParseResponseBody[BodyDataType any](body []byte) (*BodyDataType, error) {
	parsed := struct {
		Status int          `json:"status"`
		Data   BodyDataType `json:"data"`
	}{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, failure.ErrWithTrace(fmt.Errorf("Failed to unmarshal body: %w", err))
	}
	if parsed.Status != http.StatusOK {
		return nil, failure.ErrWithTrace(fmt.Errorf(
			"Got status=%d. Failed to send request: %+v", parsed.Status, parsed.Data,
		))
	}
	if parsed.Status == http.StatusUnauthorized {
		return nil, failure.ErrWithTrace(errors.New("Unauthorized request"))
	}
	return &parsed.Data, nil
}
