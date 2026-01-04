package atylabapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ApiPostInterface interface {
	Post(
		path string,
		body map[string][]string,
		ctx *ApiCtxSvc,
	) ([]byte, error)
}

type ApiPostStruct struct {
	client  *http.Client
	apiKey  string
	baseUrl string
}

func NewApiPostStruct(
	apiKey string,
	baseUrl string,
) *ApiPostStruct {
	return &ApiPostStruct{
		apiKey:  apiKey,
		baseUrl: baseUrl,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *ApiPostStruct) Post(
	path string,
	body map[string][]string,
	ctx *ApiCtxSvc,
) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx.Ctx,
		http.MethodPost,
		s.baseUrl+path,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Common-Key", s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
