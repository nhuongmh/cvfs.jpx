package ieservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/nhuongmh/cfvs.jpx/pkg/logger"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/ie"
	"github.com/nhuongmh/cfvs.jpx/pkg/utils"
	"github.com/pkg/errors"
)

func (ies *IEservice) FetchArticleUrl(ctx context.Context, link string) (*ie.Article, error) {
	cvfspyUrl := fmt.Sprintf("http://localhost:5000/api/article?url=%s", link)
	resp, err := http.Get(cvfspyUrl)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch article from link")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("failed to fetch article from link, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	var article ie.Article
	err = json.Unmarshal(body, &article)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response body")
	}

	return &article, nil
}

// parse and propose article for user to approve
func (ies *IEservice) FetchArticleUrlWithLocalScript(link string) (*ie.Article, error) {
	output, err := utils.ExecuteHostCmd(fmt.Sprintf("python3 scripts/article_fetcher.py %v", link))
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch article from link")
		// Find the file path in the output with the pattern "ARTICLE_JSON: <file_path>" using regex
	}
	var filePath string
	re := regexp.MustCompile(`ARTICLE_JSON:\s*(.+)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) > 1 {
		filePath = matches[1]
	}

	if filePath == "" {
		return nil, errors.New("failed to find ARTICLE_JSON file path in output")
	}

	logger.Log.Info().Msgf("parsed article: %s", filePath)

	// Load the article from the file path
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch open file %s", filePath)
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file %s", filePath)
	}
	var article ie.Article
	err = json.Unmarshal(byteValue, &article)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal json file %s", filePath)
	}

	return &article, nil
}

func (ies *IEservice) ExtractVocab(ctx context.Context, id uint64) (*[]ie.ProposeWord, error) {
	article, err := ies.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get article")
	}
	//serialize article to pass to http.Post command
	articleJson, err := json.Marshal(article)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal article")
	}

	cvfspyUrl := fmt.Sprintf("http://localhost:5000/api/vocab_extractor")
	resp, err := http.Post(cvfspyUrl, "application/json", bytes.NewBuffer(articleJson))
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch article from link")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("failed to fetch article from link, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	var words []ie.ProposeWord
	err = json.Unmarshal(body, &words)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response body")
	}

	return &words, nil
}
