package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nhuongmh/cfvs.jpx/api/controller"
	"github.com/nhuongmh/cfvs.jpx/bootstrap"
	"github.com/nhuongmh/cfvs.jpx/pkg/model/langfi"
	"github.com/nhuongmh/cfvs.jpx/pkg/service/jpxgen"
)

func NewJpxServiceRouter(app *bootstrap.Application, repo langfi.PracticeRepo, timeout time.Duration, publicRouter, privateRouter *gin.RouterGroup) {

	ts := jpxgen.NewJpxService(repo, timeout, app.Env)
	tc := &controller.JpxController{JpxService: ts}

	privateRouter.PUT(DEFAULT_API_PREFIX+"/core/initdb", tc.InitData)
	privateRouter.PUT(DEFAULT_API_PREFIX+"/core/buildcards", tc.GenerateProposalCards)
	publicRouter.GET(DEFAULT_API_PREFIX+"/core/langs", tc.GetAvailableLang)
	publicRouter.GET(DEFAULT_API_PREFIX+"/process/:lang-id/fetch", tc.FetchProposal)
	publicRouter.POST(DEFAULT_API_PREFIX+"/process/:lang-id/submit", tc.SubmitProposal)
	// publicRouter.GET(DEFAULT_API_PREFIX+"/process/groups", tc.GetProcessGroups)
	// publicRouter.GET(DEFAULT_API_PREFIX+"/process/:lang-id/:group-id", tc.GetProposalGroups)

}
