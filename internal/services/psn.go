package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type GetTitlesResponse struct {
	Titles []Title `json:"trophyTitles"`
}

type Title struct {
	ID          string `json:"npCommunicationId" db:"psnID"`
	ServiceName string `json:"npServiceName" db:"psnServiceName"`
	Name        string `json:"trophyTitleName" db:"name"`
	Description string `json:"trophyTitleDetail" db:"description"`
	IconURL     string `json:"trophyTitleIconUrl" db:"iconURL"`
	Platform    string `json:"trophyTitlePlatform" db:"platform"`
}

type GetTrophyGroupsResponse struct {
	TrophyGroups []TrophyGroup `json:"trophyGroups"`
}

type TrophyGroup struct {
	ID      string `json:"trophyGroupId" db:"psnID"`
	Name    string `json:"trophyGroupName" db:"name"`
	IconURL string `json:"trophyGroupIconUrl" db:"iconURL"`
}

type GetTrophiesResponse struct {
	Trophies []Trophy `json:"trophies"`
}

type Trophy struct {
	ID          uint   `json:"trophyID" db:"psnID"`
	Hidden      bool   `json:"trophyHidden" db:"hidden"`
	Rarity      string `json:"trophyType" db:"rarity"`
	Name        string `json:"trophyName" db:"name"`
	Description string `json:"trophyDetail" db:"description"`
	IconURL     string `json:"trophyIconUrl" db:"iconURL"`
	GroupID     string `json:"trophyGroupId"`
}

type PsnService interface {
	GetTitles() ([]Title, error)
	GetTrophyGroups(gameID, service string) ([]TrophyGroup, error)
	GetTrophies(gameID, service string) ([]Trophy, error)
}

type PsnClient struct {
	npssoToken   string
	accessToken  string
	refreshToken string
}

func fetchAccessCode(npssoToken string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	requestURL := "https://ca.account.sony.com/api/authz/v3/oauth/authorize"
	params := url.Values{
		"access_type":   {"offline"},
		"client_id":     {"09515159-7237-4370-9b40-3806e67c0891"},
		"redirect_uri":  {"com.scee.psxandroid.scecompcall://redirect"},
		"response_type": {"code"},
		"scope":         {"psn:mobile.v2.core psn:clientapp"},
	}
	requestURL += "?" + params.Encode()

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Cookie", "npsso="+npssoToken)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	location := resp.Header.Get("Location")
	if location == "" || !strings.Contains(location, "?code=") {
		return "", errors.New("failed to retrieve PSN access code from provided NPSSO token")
	}

	parsedLocation, err := url.Parse(location)
	if err != nil {
		return "", err
	}

	code := parsedLocation.Query().Get("code")
	if code == "" {
		return "", errors.New("failed to parse access code from Location header")
	}

	return code, nil
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func fetchAccessToken(accessCode string) (string, string, error) {
	client := &http.Client{}

	data := url.Values{
		"code":         {accessCode},
		"redirect_uri": {"com.scee.psxandroid.scecompcall://redirect"},
		"grant_type":   {"authorization_code"},
		"token_format": {"jwt"},
	}

	req, err := http.NewRequest("POST", "https://ca.account.sony.com/api/authz/v3/oauth/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic MDk1MTUxNTktNzIzNy00MzcwLTliNDAtMzgwNmU2N2MwODkxOnVjUGprYTV0bnRCMktxc1A=")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var res tokenResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", "", err
	}

	return res.AccessToken, res.RefreshToken, nil
}

func NewPsnClient(npssoToken string) *PsnClient {
	code, err := fetchAccessCode(npssoToken)
	if err != nil {
		log.Fatalln("Failed to fetch PSN access code:", err)
	}

	accessToken, refreshToken, err := fetchAccessToken(code)
	if err != nil {
		log.Fatalln("Failed to fetch PSN access token:", err)
	}

	return &PsnClient{
		npssoToken:   npssoToken,
		accessToken:  accessToken,
		refreshToken: refreshToken,
	}
}

func (c *PsnClient) GetTitles() ([]Title, error) {
	req, err := http.NewRequest("GET", "https://m.np.playstation.com/api/trophy/v1/users/me/trophyTitles?limit=800", nil)
	if err != nil {
		return nil, err
	}

	// Add the headers required for authentication.
	req.Header.Add("Authorization", "Bearer "+c.accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("Failed to get titles, status code: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var titlesResponse GetTitlesResponse
	if err := json.Unmarshal(body, &titlesResponse); err != nil {
		return nil, err
	}

	// If the title name ends with " trophies" (case-insensitive), remove it.
	for i, title := range titlesResponse.Titles {
		if strings.HasSuffix(strings.ToLower(strings.TrimSpace(title.Name)), " trophies") {
			titlesResponse.Titles[i].Name = strings.TrimSpace(title.Name[:len(title.Name)-len(" trophies")])
		}
	}

	return titlesResponse.Titles, nil
}

func (c *PsnClient) GetTrophyGroups(gameID, service string) ([]TrophyGroup, error) {
	requestURL := fmt.Sprintf("https://m.np.playstation.com/api/trophy/v1/npCommunicationIds/%v/trophyGroups", gameID)
	params := url.Values{
		"npServiceName": {service},
	}
	requestURL += "?" + params.Encode()

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	// Add the headers required for authentication.
	req.Header.Add("Authorization", "Bearer "+c.accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("Failed to get titles, status code: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var trophyGroupsResponse GetTrophyGroupsResponse
	if err := json.Unmarshal(body, &trophyGroupsResponse); err != nil {
		return nil, err
	}

	return trophyGroupsResponse.TrophyGroups, nil
}

func (c *PsnClient) GetTrophies(gameID, service string) ([]Trophy, error) {
	requestURL := fmt.Sprintf("https://m.np.playstation.com/api/trophy/v1/npCommunicationIds/%v/trophyGroups/all/trophies", gameID)
	params := url.Values{
		"npServiceName": {service},
	}
	requestURL += "?" + params.Encode()

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	// Add the headers required for authentication.
	req.Header.Add("Authorization", "Bearer "+c.accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("Failed to get titles, status code: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var trophiesResponse GetTrophiesResponse
	err = json.Unmarshal(body, &trophiesResponse)

	return trophiesResponse.Trophies, err
}
