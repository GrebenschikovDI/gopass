package controller

import (
	"GoPass/internal/agent/config"
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func Login(ctx context.Context, cfg *config.ClientConfig, username, password string) ([]*http.Cookie, error) {
	server := cfg.Server
	u := fmt.Sprintf("%s/api/user/login", server)

	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.PostForm = data

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned non-200 status code: %d", resp.StatusCode)
	}

	return resp.Cookies(), nil
}
