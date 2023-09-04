package info

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const Url = "https://api.github.com/repos/lukecarr/trophies/releases/latest"

type ReleaseInfo struct {
	TagName     string `json:"tag_name"`
	PublishedAt string `json:"published_at"`
}

func GetLatestVersion() (string, string, error) {
	resp, err := http.Get(Url)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch GitHub releases: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response body: %w", err)
	}

	var release ReleaseInfo

	if err := json.Unmarshal(body, &release); err != nil {
		return "", "", fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return release.TagName, release.PublishedAt[:10], nil
}
