package services

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type MetadataService interface {
	SearchGame(name string) (*GameMetadata, error)
}

type MetadataServiceRawg struct {
	ApiKey string
}

type GameMetadata struct {
	Name               string `json:"name"`
	BackgroundImageURL string `json:"background_image"`
}

type SearchGameResponse struct {
	Results []GameMetadata `json:"results"`
}

func (m *MetadataServiceRawg) SearchGame(name string) (*GameMetadata, error) {
	requestURL := "https://api.rawg.io/api/games"
	params := url.Values{
		"key":          {m.ApiKey},
		"search":       {name},
		"search_exact": {"true"},
		"platforms":    {"16,18,19,187"},
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
