package security

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const turnstileVerifyURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"

type TurnstileResponse struct {
	Success bool `json:"success"`
}

type TurnstileService struct {
	secretKey string
	client    *http.Client
}

func NewTurnstileService(secretKey string) *TurnstileService {
	return &TurnstileService{
		secretKey: secretKey,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *TurnstileService) Verify(ctx context.Context, token string) error {
	if token == "" || s.secretKey == "" {
		return errors.New("turnstile token or secret key is missing")
	}

	form := url.Values{}
	form.Set("secret", s.secretKey)
	form.Set("response", token)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		turnstileVerifyURL,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("turnstile request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("turnstile returned status: %d", resp.StatusCode)
	}

	var result TurnstileResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if !result.Success {
		return errors.New("turnstile verification failed")
	}

	return nil
}
