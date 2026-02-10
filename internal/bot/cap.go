package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type CapVerifier struct {
	SiteVerifyURL string
	Secret        string
	Client        *http.Client
}

func (v *CapVerifier) Enabled() bool {
	return v != nil && v.SiteVerifyURL != "" && v.Secret != ""
}

type capVerifyRequest struct {
	Secret   string `json:"secret"`
	Response string `json:"response"`
}

type capVerifyResponse struct {
	Success bool `json:"success"`
}

func (v *CapVerifier) Verify(ctx context.Context, token string) error {
	if !v.Enabled() {
		return nil
	}
	if token == "" {
		return errors.New("missing cap token")
	}

	client := v.Client
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}

	payload, err := json.Marshal(capVerifyRequest{Secret: v.Secret, Response: token})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, v.SiteVerifyURL, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("cap siteverify returned status %d", resp.StatusCode)
	}

	var result capVerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	if !result.Success {
		return errors.New("cap verification failed")
	}

	return nil
}
