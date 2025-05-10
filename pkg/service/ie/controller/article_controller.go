package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nhuongmh/cfvs.jpx/pkg/logger"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/ie"
	ieservice "github.com/nhuongmh/cfvs.jpx/pkg/service/ie"
)

type IeController struct {
	Service *ieservice.IEservice
}

type pageMeta struct {
	Page     uint64 `json:"page"`
	PageSize uint64 `json:"page_size"`
	Total    int    `json:"total"`
}

func (tc *IeController) GetAllArticle(c *gin.Context) {
	page, pageSize := tc.parsePagination(c, 1, 20)
	articles, total, err := tc.Service.GetAllArticles(c, pageSize, page)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to get all articles")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get all articles"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"articles": articles, "meta": pageMeta{Page: page, PageSize: pageSize, Total: total}})
}

func (tc *IeController) SaveArticle(c *gin.Context) {
	var article ie.Article
	err := c.ShouldBindJSON(&article)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to bind article")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind article"})
		return
	}
	savedArticle, err := tc.Service.SaveArticle(c, &article)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to save article")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save article"})
		return
	}
	c.JSON(http.StatusOK, savedArticle)
}

func (tc *IeController) GetArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to parse id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse id"})
		return
	}
	article, err := tc.Service.GetArticle(c, id)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to get article")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get article"})
		return
	}
	c.JSON(http.StatusOK, article)
}

func (tc *IeController) DeleteArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to parse id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse id"})
		return
	}
	err = tc.Service.DeleteArticle(c, id)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to delete article")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete article"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "article deleted"})
}

func (tc *IeController) FindArticleByTitle(c *gin.Context) {
	title := c.Query("title")
	if title == "" {
		logger.Log.Error().Msg("title is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is empty"})
		return
	}
	articles, err := tc.Service.FindArticleByTitle(c, title)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to find article by title")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find article by title"})
		return
	}
	c.JSON(http.StatusOK, articles)
}

func (tc *IeController) ParseArticleFromUrl(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		logger.Log.Error().Msg("url is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is empty"})
		return
	}
	article, err := tc.Service.FetchArticleUrl(c, url)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to fetch article from url")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch article from url"})
		return
	}

	c.JSON(http.StatusOK, article)
}

func (tc *IeController) GetArticleReading(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to parse id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse id"})
		return
	}
	articleReading, err := tc.Service.GetArticleReading(c, id)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to get question")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get question"})
		return
	}
	c.JSON(http.StatusOK, articleReading)
}

func (tc *IeController) ReGenArticleReading(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to parse id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse id"})
		return
	}
	err = tc.Service.ReGenArticleReading(c, id)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to regenerate question")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to regenerate question"})
		return
	}
	c.JSON(http.StatusOK, "OK")
}

func (tc *IeController) SubmitReadingTest(c *gin.Context) {
	readingId, err := strconv.ParseUint(c.Param("reading_id"), 10, 64)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to parse id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse id"})
		return
	}
	var answers map[int]string
	err = c.ShouldBindJSON(&answers)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to bind answers")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind answers"})
		return
	}

	testResult, err := tc.Service.GradeQuestionSubmit(c, readingId, answers)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to grade question submit")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to grade question submit"})
		return
	}
	c.JSON(http.StatusOK, testResult)
}

func (tc *IeController) GetTestSubmission(c *gin.Context) {
	submitId, err := strconv.ParseUint(c.Param("submit_id"), 10, 64)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to parse id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse id"})
		return
	}
	testResult, err := tc.Service.GetTestSubmission(c, submitId)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to get test submission")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get test submission"})
		return
	}
	c.JSON(http.StatusOK, testResult)
}

func (tc *IeController) GetTestSubmissionByReadingId(c *gin.Context) {
	readingId, err := strconv.ParseUint(c.Param("reading_id"), 10, 64)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to parse id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse id"})
		return
	}
	testResults, err := tc.Service.GetTestSubmissionByReadingId(c, readingId)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to get test submission by reading id")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get test submission by reading id"})
		return
	}
	c.JSON(http.StatusOK, testResults)
}

func (tc *IeController) DeleteTestSubmission(c *gin.Context) {
	submitId, err := strconv.ParseUint(c.Param("submit_id"), 10, 64)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to parse id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse id"})
		return
	}
	err = tc.Service.DeleteTestSubmission(c, submitId)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to delete test submission")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete test submission"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "test submission deleted"})
}

func (tc *IeController) ExtractProposedWordsForArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to parse id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse id"})
		return
	}
	words, err := tc.Service.ExtractVocab(c, id)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to extracting words")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to extracting words"})
		return
	}
	c.JSON(http.StatusOK, words)
}

func (tc *IeController) HandleVocabProposalSubmit(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to parse id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse id"})
		return
	}

	var proposals []ie.ProposeWord
	err = c.ShouldBindJSON(&proposals)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to bind answers")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind answers"})
		return
	}

	list, err := tc.Service.GenVocabListFromProposal(c, id, &proposals)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to gen vocab list from proposal")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to gen vocab list from proposal"})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (tc *IeController) GetVocabListByArticleId(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to parse id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse id"})
		return
	}

	list, err := tc.Service.GetVocabListByArticleId(c, id)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to get vocab list by article id")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get vocab list by article id"})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (tc *IeController) parsePagination(c *gin.Context, defaultPage, defaultSize uint64) (uint64, uint64) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	pageNum, err := strconv.ParseUint(page, 10, 64)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to parse page")
		pageNum = defaultPage
	}
	size, err := strconv.ParseUint(pageSize, 10, 64)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to parse page size")
		size = defaultSize
	}
	return pageNum, size
}
