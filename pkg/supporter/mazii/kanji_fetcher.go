package mazii

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

func (m *MaziiFetcher) FetchMaziiKanji(kj string) (*KanjiResultResp, error) {
	// Prepare the request data
	requestPayload := map[string]interface{}{
		"dict":  "javi",
		"type":  "kanji",
		"query": kj,
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
	var kanjiReps KanjiResultResp
	err = json.NewDecoder(resp.Body).Decode(&kanjiReps)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode response body")
	}

	if len(kanjiReps.Results) == 0 {
		return nil, errors.New("no results found")
	}

	return &kanjiReps, nil
}

func (m *MaziiFetcher) FetchBestComment(wordID int) (*MaziiCommentEntry, error) {
	// Prepare the request data
	searchData := map[string]interface{}{
		"dict":   "javi",
		"type":   "kanji",
		"wordId": wordID,
	}

	jsonData, err := json.Marshal(searchData)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal search data for wordID %d", wordID)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", "https://api.mazii.net/api/get-mean", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create HTTP request for wordID %d", wordID)
	}

	// Set headers
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept-lanuage", "en-US,en;q=0.9,ja;q=0.8,vi;q=0.7")
	req.Header.Set("authority", "mazii.net")

	// Make the request
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to make HTTP request for wordID %d", wordID)
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("request failed with status %d for wordID %d", resp.StatusCode, wordID)
	}

	// Read and parse the response body
	var commentResp MaziiComments
	err = json.NewDecoder(resp.Body).Decode(&commentResp)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode response body for wordID %d", wordID)
	}

	if len(commentResp.Comments) == 0 {
		return nil, errors.New(fmt.Sprintf("no comments found for wordID %d", wordID))
	}

	return &commentResp.Comments[0], nil
}

func extractExamples(examples []map[string]string) string {
	var exampleBuilder string
	for i, example := range examples {
		if i > 4 {
			break
		}
		exampleBuilder += fmt.Sprintf("%s - %s - %s - %s\n", example["w"], example["h"], example["p"], example["m"])
	}
	return exampleBuilder
}
