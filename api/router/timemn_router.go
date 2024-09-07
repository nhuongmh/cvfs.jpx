package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nhuongmh/cfvs.jpx/api/controller"
	"github.com/nhuongmh/cfvs.jpx/bootstrap"
	"github.com/nhuongmh/cfvs.jpx/pkg/service/jpxservice"
	"github.com/nhuongmh/cfvs.jpx/pkg/service/jpxservice/repo"
)

func NewJpxServiceRouter(app *bootstrap.Application, timeout time.Duration, publicRouter, privateRouter *gin.RouterGroup) {

	tr := repo.NewJpxRepo(app.DB)
	ts := jpxservice.NewJpxService(tr, timeout, app.Env)
	tc := &controller.JpxController{JpxService: ts}

	privateRouter.PUT(DEFAULT_API_PREFIX+"/jpx/initdb", tc.InitData)
}
