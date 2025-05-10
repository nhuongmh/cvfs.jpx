package ieservice

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/nhuongmh/cfvs.jpx/pkg/logger"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/ie"
	"github.com/pkg/errors"
)

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
		vocabList.Vocabs = append(vocabList.Vocabs, *processed...)
		return ies.repo.UpdateVocabList(ctx, vocabList)
	}

	vocabList = &ie.IeVocabList{
		Name:         article.Title,
		RefArticleID: article.ID,
		Vocabs:       *processed,
	}

	return ies.repo.SaveVocabList(ctx, vocabList)
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

	processed := make([]ie.IeVocab, 0)
	for w, wordObj := range wordMap {
		wordObj.Word = w
		processed = append(processed, wordObj)
	}
	return &processed, nil
}
