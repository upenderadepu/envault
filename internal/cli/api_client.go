package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type cliClient struct {
	baseURL string
	token   string
	http    *http.Client
}

func newAuthClient() *cliClient {
	return &cliClient{
		baseURL: getAPIURL(),
		token:   getAuthToken(),
		http:    &http.Client{},
	}
}

func (c *cliClient) request(method, path string, body interface{}) ([]byte, int, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp struct {
			Error string `json:"error"`
		}
		if json.Unmarshal(respBody, &errResp) == nil && errResp.Error != "" {
			return nil, resp.StatusCode, fmt.Errorf("%s", errResp.Error)
		}
		return nil, resp.StatusCode, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	return respBody, resp.StatusCode, nil
}
