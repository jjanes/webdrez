package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
)



func IsKickStreamLive(username string) (bool, error) {
	url := fmt.Sprintf("https://kick.com/api/v1/channels/%s", username)

	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, fmt.Errorf("non-200 response: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	// Kick's API returns a field called "livestream"
	var result struct {
		Livestream interface{} `json:"livestream"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return false, err
	}

	// If livestream is not nil, stream is live
	return result.Livestream != nil, nil
}
