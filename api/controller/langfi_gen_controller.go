package controller

import (
	"net/http"
	"strconv"

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

func (jctl *JpxController) GetAvailableLang(gc *gin.Context) {
	langs := []string{"jp"}

	gc.JSON(http.StatusOK, langs)
}

func (jctl *JpxController) FetchProposal(gc *gin.Context) {
	proposal, err := jctl.JpxService.FetchProposal(gc)
	if err != nil {
		logger.Log.Error().Err(err).Msg("request process failed")
		gc.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	gc.JSON(http.StatusOK, *proposal)
}

func (jctl *JpxController) SubmitProposal(gc *gin.Context) {
	//get card proposal from gin context

	cardIDStr := gc.DefaultQuery("cardID", "")
	status := gc.DefaultQuery("status", "")
	if cardIDStr == "" || status == "" {
		gc.JSON(http.StatusBadRequest, ErrorResponse{Message: "cardID and status are required"})
		return
	}

	cardId, err := strconv.ParseUint(cardIDStr, 10, 64)
	if err != nil {
		gc.JSON(http.StatusBadRequest, ErrorResponse{Message: "cardID must be a number"})
		return
	}

	err = jctl.JpxService.SubmitProposal(gc, cardId, status)
	if err != nil {
		logger.Log.Error().Err(err).Msg("request process failed")
		gc.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	gc.JSON(http.StatusOK, "Success")
}
