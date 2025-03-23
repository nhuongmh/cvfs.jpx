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
	articles, total, err := tc.Service.GetAllArticles(c, page, pageSize)
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
