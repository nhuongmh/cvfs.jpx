package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nhuongmh/cfvs.jpx/api/controller"
	"github.com/nhuongmh/cfvs.jpx/bootstrap"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/langfi"
	"github.com/nhuongmh/cfvs.jpx/pkg/service/jpxpractice"
)

func NewJpxPraServiceRouter(app *bootstrap.Application, repo langfi.PracticeRepo, timeout time.Duration, publicRouter, privateRouter *gin.RouterGroup) {

	ts := jpxpractice.NewJpxPracService(timeout, repo, app.Env)
	tc := &controller.PracticeController{PracticeSrv: ts}

	publicRouter.GET(DEFAULT_API_PREFIX+"/practice/:lang-id/groups", tc.GetPracticeGroups)
	publicRouter.GET(DEFAULT_API_PREFIX+"/practice/:lang-id/:group-id/fetch", tc.FetchPracticeCard)
	publicRouter.POST(DEFAULT_API_PREFIX+"/practice/:lang-id/:group-id/submit", tc.SubmitPracticeCard)
	// publicRouter.POST(DEFAULT_API_PREFIX+"/practice/:card-id", tc.GetCard)

}
