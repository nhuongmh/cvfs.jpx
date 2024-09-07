package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nhuongmh/cfvs.jpx/pkg/logger"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/jp"
)

type JpxController struct {
	JpxService jp.JpxService
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
