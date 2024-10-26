package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/langfi"
)

type PracticeController struct {
	PracticeSrv langfi.PracticeService
}

func (pctl *PracticeController) GetPracticeGroups(gc *gin.Context) {
	// langID := gc.Param("lang-id")
	groups := pctl.PracticeSrv.GetGroups(gc)
	gc.JSON(http.StatusOK, groups)
}

func (pctl *PracticeController) FetchPracticeCard(gc *gin.Context) {
	// langID := gc.Param("lang-id")
	groupID := gc.Param("group-id")
	card, err := pctl.PracticeSrv.FetchCard(gc, groupID)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	gc.JSON(http.StatusOK, card)
}

func (pctl *PracticeController) SubmitPracticeCard(gc *gin.Context) {
	cardIDStr := gc.DefaultQuery("cardID", "")
	rating := gc.DefaultQuery("rating", "")
	if cardIDStr == "" || rating == "" {
		gc.JSON(http.StatusBadRequest, ErrorResponse{Message: "cardID and rating are required"})
		return
	}

	cardId, err := strconv.ParseUint(cardIDStr, 10, 64)
	if err != nil {
		gc.JSON(http.StatusBadRequest, ErrorResponse{Message: "cardID must be a number"})
		return
	}

	rateID, err := strconv.ParseUint(rating, 10, 64)
	if err != nil {
		gc.JSON(http.StatusBadRequest, ErrorResponse{Message: "rating must be a number"})
		return
	}

	err = pctl.PracticeSrv.SubmitCard(gc, cardId, rateID)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	gc.JSON(http.StatusOK, "Success")
}

func (pctl *PracticeController) GetCard(gc *gin.Context) {
	cardIDStr := gc.Param("card-id")
	if cardIDStr == "" {
		gc.JSON(http.StatusBadRequest, ErrorResponse{Message: "cardID is required"})
		return
	}

	cardId, err := strconv.ParseUint(cardIDStr, 10, 64)
	if err != nil {
		gc.JSON(http.StatusBadRequest, ErrorResponse{Message: "cardID must be a number"})
		return
	}

	card, err := pctl.PracticeSrv.GetCard(gc, cardId)
	if err != nil {
		gc.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	gc.JSON(http.StatusOK, card)
}
