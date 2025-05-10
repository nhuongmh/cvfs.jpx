package ieservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/nhuongmh/cfvs.jpx/pkg/logger"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/ie"
	"github.com/pkg/errors"
)

func (ies *IEservice) ExtractVocab(ctx context.Context, id uint64) (*[]ie.ProposeWord, error) {
	cachedVocabs, ok := ies.vocabProposalCache[id]
	if ok {
		return cachedVocabs, nil
	}

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

	//cache the vocab list for this article
	ies.vocabProposalCache[id] = &words
	return &words, nil
}

func (ies *IEservice) GenVocabListFromProposal(ctx context.Context, articleId uint64, proposals *[]ie.ProposeWord) (*ie.IeVocabList, error) {
	//check if vocab list already exists
	article, err := ies.repo.FindByID(ctx, articleId)
	if err != nil {
		logger.Log.Warn().Msgf("Article ID %v not found", articleId)
	}
	processed, err := ies.processVocabProposalList(proposals)
	if err != nil {
		return nil, errors.Wrap(err, "failed to process vocab proposal list")
	}

	vocabList, err := ies.repo.GetAllVocabListByArticleId(ctx, article.ID)
	if err == nil {
		logger.Log.Warn().Msgf("Vocab list for article ID %v already exist, appending", article.ID)
		// vocabList.Vocabs = append(vocabList.Vocabs, *processed...)
		vocabList.Vocabs = *processed
		return ies.repo.UpdateVocabList(ctx, vocabList)
	}

	vocabList = &ie.IeVocabList{
		Name:         article.Title,
		RefArticleID: article.ID,
		Vocabs:       *processed,
	}
	// err = ies.generateAnkiDeck(vocabList)
	// if err != nil {
	// 	logger.Log.Warn().Err(err).Msg("failed to generate Anki deck")
	// }
	ies.vocabProposalCache[articleId] = nil //clear the cache for this article
	return ies.repo.SaveVocabList(ctx, vocabList)
}

func (ies *IEservice) GenAnkiDeckForVocabList(ctx context.Context, vocabListId uint64) error {
	vocabList, err := ies.repo.FindVocabListByID(ctx, vocabListId, true)
	if err != nil {
		return errors.Wrap(err, "failed to get vocab list")
	}

	return ies.generateAnkiDeck(vocabList)
}

func (ies *IEservice) GetVocabList(ctx context.Context, vocabListId uint64) (*ie.IeVocabList, error) {
	vocabList, err := ies.repo.FindVocabListByID(ctx, vocabListId, true)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get vocab list")
	}
	if vocabList == nil {
		return nil, errors.New("vocab list not found")
	}
	return vocabList, nil
}

func (ies *IEservice) GetVocabListByArticleId(ctx context.Context, articleId uint64) (*ie.IeVocabList, error) {
	vocabList, err := ies.repo.GetAllVocabListByArticleId(ctx, articleId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get vocab list")
	}
	if vocabList == nil {
		return nil, errors.New("vocab list not found")
	}
	return vocabList, nil
}

func (ies *IEservice) GetAllVocabList(ctx context.Context, limit, skip uint64) (*[]ie.IeVocabList, int, error) {
	vocabLists, total, err := ies.repo.GetAllVocabList(ctx, false, limit, skip)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to get all vocab lists")
	}
	return vocabLists, total, nil
}

func (ies *IEservice) DeleteVocabList(ctx context.Context, vocabListId uint64) error {
	vocabList, err := ies.repo.FindVocabListByID(ctx, vocabListId, false)
	if err != nil {
		return errors.Wrap(err, "failed to get vocab list")
	}
	if vocabList == nil {
		return errors.New("vocab list not found")
	}
	err = ies.repo.DeleteVocabList(ctx, vocabList.ID)
	if err != nil {
		return errors.Wrap(err, "failed to delete vocab list")
	}
	return nil
}

func (ies *IEservice) processVocabProposalList(proposals *[]ie.ProposeWord) (*[]ie.IeVocab, error) {
	listWord := make([]string, len(*proposals))
	for i, proposal := range *proposals {
		listWord[i] = proposal.Word
	}
	proposalJson, err := json.Marshal(listWord)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal proposal")
	}

	cvfspyUrl := "http://localhost:5000/api/dictionary/en/"
	resp, err := http.Post(cvfspyUrl, "application/json", bytes.NewBuffer(proposalJson))
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request to CVFSpy")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("failed to send request to CVFSpy, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	var wordMap map[string]ie.IeVocab
	err = json.Unmarshal(body, &wordMap)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response body")
	}

	// create search dict
	proposalMap := make(map[string]ie.ProposeWord)
	for _, proposal := range *proposals {
		proposalMap[proposal.Word] = proposal
	}

	processed := make([]ie.IeVocab, 0)
	for w, wordObj := range wordMap {
		wordObj.Word = w
		if proposal, ok := proposalMap[w]; ok {
			logger.Log.Debug().Msgf("found word %s in proposal: context=%v", w, proposal.Context)
			wordObj.Context = proposal.Context
			wordObj.WordFreq = proposal.WordFreq
		}
		processed = append(processed, wordObj)
	}
	return &processed, nil
}

func (ies *IEservice) generateAnkiDeck(vocabList *ie.IeVocabList) error {
	vocabJson, err := json.Marshal(vocabList)
	if err != nil {
		return errors.Wrap(err, "failed to marshal proposal")
	}

	cvfspyUrl := "http://localhost:5000/api/vocab/ankigen/local"
	resp, err := http.Post(cvfspyUrl, "application/json", bytes.NewBuffer(vocabJson))
	if err != nil {
		return errors.Wrap(err, "failed to send request to CVFSpy")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("failed to send request to CVFSpy, status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}
	logger.Log.Info().Msgf("Anki deck generated: %s", string(body))
	return nil
}
