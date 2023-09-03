package services

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
)

type MetadataService interface {
	SearchGame(name, platform, excludePlatforms string) (*GameMetadata, error)
}

type MetadataServiceRawg struct {
	ApiKey string
}

type GameMetadata struct {
	Name               string `json:"name"`
	BackgroundImageURL string `json:"background_image"`
	MetacriticScore    int    `json:"metacritic"`
	ReleaseDate        string `json:"released"`
}

type SearchGameResponse struct {
	Results []GameMetadata `json:"results"`
}

func (m *MetadataServiceRawg) SearchGame(name, platform, excludePlatforms string) (*GameMetadata, error) {
	requestURL := "https://api.rawg.io/api/games"

	// Remove all non-ASCII characters from name (i.e. trademark symbols, copyright, etc.)
	name = regexp.MustCompile(`[^\x00-\x7F]+`).ReplaceAllString(name, "")

	params := url.Values{
		"key":               {m.ApiKey},
		"search":            {name},
		"platforms":         {platform},
		"exclude_platforms": {excludePlatforms},
	}
	requestURL += "?" + params.Encode()

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var searchResponse SearchGameResponse
	err = json.NewDecoder(resp.Body).Decode(&searchResponse)
	if err != nil {
		return nil, err
	}

	if len(searchResponse.Results) == 0 {
		return nil, nil
	}

	return &searchResponse.Results[0], nil
}
