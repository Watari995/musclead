package paymentinfra

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	paymentdomain "github.com/Watari995/musclead/internal/payment/internal/domain"
)

const resendEndpoint = "https://api.resend.com/emails"

type resendMailer struct {
	apiKey     string
	from       string
	httpClient *http.Client
}

// NewResendMailer は Resend (https://resend.com) でメール送信する Mailer を返す。
// AWS SDK 非依存 (標準の net/http のみ)。
func NewResendMailer(apiKey, from string) paymentdomain.Mailer {
	return &resendMailer{
		apiKey:     apiKey,
		from:       from,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (m *resendMailer) Send(ctx context.Context, to, subject, body string) error {
	payload, err := json.Marshal(map[string]any{
		"from":    m.from,
		"to":      []string{to},
		"subject": subject,
		"text":    body,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, resendEndpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := m.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= http.StatusMultipleChoices {
		b, _ := io.ReadAll(res.Body)
		return fmt.Errorf("resend: status %d: %s", res.StatusCode, string(b))
	}
	return nil
}
