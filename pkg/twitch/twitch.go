// pkg/twitchlive/twitchlive.go
package twitch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Config struct {
	ClientID    string
	AccessToken string
	APIBaseURL  string // optional; defaults to https://api.twitch.tv/helix
	HTTPTimeout time.Duration
}

// Response format for /helix/streams
type helixStreamsResponse struct {
	Data []struct {
		ID           string `json:"id"`
		UserID       string `json:"user_id"`
		UserLogin    string `json:"user_login"`
		UserName     string `json:"user_name"`
		GameID       string `json:"game_id"`
		Type         string `json:"type"` // "live" or ""
		Title        string `json:"title"`
		ViewerCount  int    `json:"viewer_count"`
		StartedAt    string `json:"started_at"`
		Language     string `json:"language"`
		ThumbnailURL string `json:"thumbnail_url"`
	} `json:"data"`
}

// StreamStatus is what we return to callers.
type StreamStatus struct {
	Live        bool
	Title       string
	ViewerCount int
	StartedAt   time.Time
}

// CheckChannelLive checks if a given Twitch channel is currently live.
// channelName should be the login name, e.g. "ninja", not the display name.
func CheckChannelLive(cfg Config, channelName string) (*StreamStatus, error) {
	if cfg.ClientID == "" || cfg.AccessToken == "" {
		return nil, fmt.Errorf("twitchlive: ClientID and AccessToken are required")
	}

	baseURL := cfg.APIBaseURL
	if baseURL == "" {
		baseURL = "https://api.twitch.tv/helix"
	}

	u, err := url.Parse(baseURL + "/streams")
	if err != nil {
		return nil, fmt.Errorf("twitchlive: parse base url: %w", err)
	}

	q := u.Query()
	q.Set("user_login", channelName)
	u.RawQuery = q.Encode()

	client := &http.Client{
		Timeout: cfg.HTTPTimeout,
	}
	if client.Timeout == 0 {
		client.Timeout = 5 * time.Second
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("twitchlive: new request: %w", err)
	}

	req.Header.Set("Client-Id", cfg.ClientID)
	req.Header.Set("Authorization", "Bearer "+cfg.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("twitchlive: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("twitchlive: non-200 status: %s", resp.Status)
	}

	var data helixStreamsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("twitchlive: decode json: %w", err)
	}

	if len(data.Data) == 0 {
		// not live
		return &StreamStatus{Live: false}, nil
	}

	s := data.Data[0]
	startedAt, _ := time.Parse(time.RFC3339, s.StartedAt)

	return &StreamStatus{
		Live:        s.Type == "live",
		Title:       s.Title,
		ViewerCount: s.ViewerCount,
		StartedAt:   startedAt,
	}, nil
}
