package transport

import (
	"GoPass/internal/agent/config"
	"context"
	"fmt"
	"io"
	"net/http"
)

func GetList(ctx context.Context, cfg config.ClientConfig, cookies []*http.Cookie) ([]byte, error) {
	server := cfg.Server
	url := fmt.Sprintf("%s/api/records", server)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	client := http.Client{}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned non-200 status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return body, nil
}

func DeleteRecord(ctx context.Context, cfg *config.ClientConfig, cookies []*http.Cookie, recordID int) error {
	server := cfg.Server
	url := fmt.Sprintf("%s/api/records/%d", server, recordID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	client := http.Client{}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("server returned non-204 status code: %d", resp.StatusCode)
	}

	return nil
}
