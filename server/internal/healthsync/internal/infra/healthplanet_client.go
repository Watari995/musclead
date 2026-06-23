package healthsyncinfra

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	healthsyncdomain "github.com/Watari995/musclead/internal/healthsync/internal/domain"
)

const healthPlanetBaseURL = "https://www.healthplanet.jp"

// jst は HealthPlanet API が JST 基準で日時を解釈するため、パッケージ全体で共有する。
// "Asia/Tokyo" は Go 標準 tzdata に含まれており失敗しない。
var jst, _ = time.LoadLocation("Asia/Tokyo")

type HealthPlanetClient struct {
	httpClient   *http.Client
	clientID     string
	clientSecret string
}

func NewHealthPlanetClient(httpClient *http.Client, clientID, clientSecret string) *HealthPlanetClient {
	return &HealthPlanetClient{
		httpClient:   httpClient,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// noRedirectClient はリダイレクトを追従しない HTTP クライアント。
// HealthPlanet のトークンエンドポイントはエラー時に callback へ 302 リダイレクトするため、
// 追従すると最終的にフロントエンド HTML が返り JSON デコード失敗で 500 になる。
var noRedirectClient = &http.Client{
	Timeout: 10 * time.Second,
	CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func (c *HealthPlanetClient) ExchangeCode(ctx context.Context, code string) (accessToken, refreshToken string, expiresAt time.Time, err error) {
	values := url.Values{}
	values.Set("grant_type", "authorization_code")
	values.Set("client_id", c.clientID)
	values.Set("client_secret", c.clientSecret)
	values.Set("redirect_uri", "https://api.musclead.com/integrations/healthplanet/callback")
	values.Set("code", code)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		healthPlanetBaseURL+"/oauth/token",
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return "", "", time.Time{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// RFC 6749 §2.3.1: Basic Auth でも client_id/secret を送る
	req.SetBasicAuth(c.clientID, c.clientSecret)

	resp, err := noRedirectClient.Do(req)
	if err != nil {
		return "", "", time.Time{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("healthplanet code exchange failed",
			"status", resp.StatusCode,
			"location", resp.Header.Get("Location"),
			"body", string(body),
		)
		return "", "", time.Time{}, fmt.Errorf("healthplanet code exchange: status %d location=%s body=%s",
			resp.StatusCode, resp.Header.Get("Location"), string(body))
	}

	var result tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", time.Time{}, err
	}

	expiresAt = time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)
	return result.AccessToken, result.RefreshToken, expiresAt, nil
}

func (c *HealthPlanetClient) RefreshToken(ctx context.Context, refreshToken string) (accessToken, newRefreshToken string, expiresAt time.Time, err error) {
	values := url.Values{}
	values.Set("grant_type", "refresh_token")
	values.Set("client_id", c.clientID)
	values.Set("client_secret", c.clientSecret)
	values.Set("refresh_token", refreshToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		healthPlanetBaseURL+"/oauth/token",
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return "", "", time.Time{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", "", time.Time{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", time.Time{}, fmt.Errorf("healthplanet token refresh: unexpected status %d", resp.StatusCode)
	}

	var result tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", time.Time{}, err
	}

	expiresAt = time.Now().Add(time.Duration(result.ExpiresIn) * time.Second)
	return result.AccessToken, result.RefreshToken, expiresAt, nil
}

type innerscanResponse struct {
	Data []innerscanEntry `json:"data"`
}

type innerscanEntry struct {
	Date    string `json:"date"`    // "202601151030" (YYYYMMDDHHMM, JST)
	Tag     string `json:"tag"`     // "6021"=体重, "6022"=体脂肪率, "6023"=骨格筋量
	Keydata string `json:"keydata"` // "66.7"
}

func (c *HealthPlanetClient) FetchMetrics(ctx context.Context, accessToken string, from, to time.Time) ([]healthsyncdomain.BodyMetrics, error) {
	values := url.Values{}
	values.Set("access_token", accessToken)
	values.Set("date", "1")
	values.Set("from", from.In(jst).Format("20060102150405"))
	values.Set("to", to.In(jst).Format("20060102150405"))
	values.Set("tag", "6021,6022,6023")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		healthPlanetBaseURL+"/status/innerscan.json",
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("healthplanet innerscan: unexpected status %d", resp.StatusCode)
	}

	var result innerscanResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return toBodyMetrics(result.Data)
}

func toBodyMetrics(entries []innerscanEntry) ([]healthsyncdomain.BodyMetrics, error) {
	metricsMap := map[string]*healthsyncdomain.BodyMetrics{}
	order := []string{}

	for _, entry := range entries {
		if _, ok := metricsMap[entry.Date]; !ok {
			measuredAt, err := time.ParseInLocation("200601021504", entry.Date, jst)
			if err != nil {
				continue
			}
			metricsMap[entry.Date] = &healthsyncdomain.BodyMetrics{MeasuredAt: measuredAt}
			order = append(order, entry.Date)
		}

		val, err := strconv.ParseFloat(entry.Keydata, 64)
		if err != nil {
			continue
		}

		m := metricsMap[entry.Date]
		switch entry.Tag {
		case "6021":
			m.Weight = val
		case "6022":
			m.BodyFatPercent = &val
		case "6023":
			m.SkeletalMuscleKg = &val
		}
	}

	metrics := make([]healthsyncdomain.BodyMetrics, 0, len(order))
	for _, date := range order {
		metrics = append(metrics, *metricsMap[date])
	}
	return metrics, nil
}
