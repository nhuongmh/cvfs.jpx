package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nhuongmh/cfvs.jpx/pkg/logger"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/jp"
)

type JpxController struct {
	JpxService jp.JpxGeneratorService
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func (jctl *JpxController) InitData(gc *gin.Context) {
	err := jctl.JpxService.InitData(gc)
	if err != nil {
		logger.Log.Error().Err(err).Msg("request process failed")
		gc.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	gc.JSON(http.StatusOK, "Success")
}

func (jctl *JpxController) GenerateProposalCards(gc *gin.Context) {
	cards, err := jctl.JpxService.BuildCards(gc)
	if err != nil {
		logger.Log.Error().Err(err).Msg("request process failed")
		gc.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	gc.JSON(http.StatusOK, *cards)
}

func (jctl *JpxController) GetWordList(gc *gin.Context) {
	words := jctl.JpxService.GetWordList(gc)

	gc.JSON(http.StatusOK, *words)
}
