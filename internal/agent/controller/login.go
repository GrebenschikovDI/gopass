package controller

import (
	"GoPass/internal/agent/config"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func Login(ctx context.Context, cfg *config.ClientConfig, username, password string) ([]*http.Cookie, error) {
	return authRequest(ctx, cfg, username, password, "login")
}

func Register(ctx context.Context, cfg *config.ClientConfig, username, password string) ([]*http.Cookie, error) {
	return authRequest(ctx, cfg, username, password, "register")
}

func authRequest(ctx context.Context, cfg *config.ClientConfig, username, password, request string) ([]*http.Cookie, error) {
	server := cfg.Server
	u := fmt.Sprintf("%s/api/user/%s", server, request)

	bodyData := map[string]string{
		"login":    username,
		"password": password,
	}

	bodyBytes, err := json.Marshal(bodyData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON body: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
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
