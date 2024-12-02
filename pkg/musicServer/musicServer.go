package musicServer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type MusicInfoClient struct {
	BaseURL string
	Client  *http.Client
}

func NewMusicInfoClient(baseURL string) *MusicInfoClient {
	return &MusicInfoClient{
		BaseURL: baseURL,
		Client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (m *MusicInfoClient) GetSongDetails(group, song string) (*SongDetail, error) {
	url := fmt.Sprintf("%s/info?group=%s&song=%s", m.BaseURL, group, song)
	resp, err := m.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call music info API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	var songDetail SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &songDetail, nil
}
