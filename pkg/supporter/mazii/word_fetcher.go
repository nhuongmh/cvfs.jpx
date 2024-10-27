package mazii

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

func (m *MaziiFetcher) SearchMaziiWord(word string) (*MaziiSearchResult, error) {
	// Prepare the request data
	requestPayload := map[string]interface{}{
		"dict":  "javi",
		"type":  "word",
		"query": word,
		"page":  1,
	}

	jsonData, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request payload")
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", "https://mazii.net/api/search", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create HTTP request")
	}

	// Set headers
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept-lanuage", "en-US,en;q=0.9,ja;q=0.8,vi;q=0.7")
	req.Header.Set("authority", "mazii.net")

	// Make the request
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make HTTP request")
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP request failed with status %d", resp.StatusCode)
	}

	// Read and parse the response body
	var searchResult MaziiSearchResult
	err = json.NewDecoder(resp.Body).Decode(&searchResult)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode response body")
	}

	if len(searchResult.Results) == 0 {
		return nil, errors.New("no results found")
	}

	return &searchResult, nil
}
