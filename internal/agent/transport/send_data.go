package transport

import (
	"GoPass/internal/agent/config"
	"GoPass/internal/agent/records"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func CreateRecord(ctx context.Context, cfg config.ClientConfig, record records.Record, cookies []*http.Cookie) error {
	server := cfg.Server
	url := fmt.Sprintf("%s/api/records", server)

	// Преобразование записи в JSON
	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal record to JSON: %v", err)
	}

	// Создание запроса с телом JSON
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(recordJSON))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Добавление кук в запрос
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	// Установка заголовка Content-Type
	req.Header.Set("Content-Type", "application/json")

	// Отправка запроса
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Проверка статуса ответа
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("server returned non-201 status code: %d", resp.StatusCode)
	}

	// Все прошло успешно
	return nil
}
