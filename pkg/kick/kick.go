package kick

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type kickChannel struct {
	IsLive     bool `json:"is_live"`
	Livestream *struct {
		IsLive bool `json:"is_live"`
	} `json:"livestream"`
}

func IsKickLive(slug string) (bool, error) {
	url := "https://kick.com/api/v2/channels/" + slug // check version in DevTools

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", "https://kick.com/"+slug)
	req.Header.Set("Origin", "https://kick.com")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("bad status: %s body: %s", resp.Status, string(body))
	}

	var ch kickChannel
	if err := json.Unmarshal(body, &ch); err != nil {
		return false, err
	}

	if ch.Livestream != nil {
		return ch.Livestream.IsLive, nil
	}
	return ch.IsLive, nil
}
