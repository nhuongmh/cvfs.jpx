package router

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nhuongmh/cfvs.jpx/bootstrap"
	ieservice "github.com/nhuongmh/cfvs.jpx/pkg/service/ie"
	"github.com/nhuongmh/cfvs.jpx/pkg/service/ie/controller"
)

func NewIeAiRouter(app *bootstrap.Application, timeout time.Duration, publicRouter, privateRouter *gin.RouterGroup) {

	ts := ieservice.NewIEservice(timeout, app.Env, app.DB)
	tc := &controller.IeController{Service: ts}

	privateRouter.GET(DEFAULT_API_PREFIX+"/ie/article", tc.GetAllArticle)
	privateRouter.GET(DEFAULT_API_PREFIX+"/ie/article/:id", tc.GetArticle)
	privateRouter.POST(DEFAULT_API_PREFIX+"/ie/article", tc.SaveArticle)
	privateRouter.DELETE(DEFAULT_API_PREFIX+"/ie/article/:id", tc.DeleteArticle)
	privateRouter.GET(DEFAULT_API_PREFIX+"/ie/article/url", tc.ParseArticleFromUrl)
	privateRouter.GET(DEFAULT_API_PREFIX+"/ie/article/:id/reading", tc.GetArticleReading)
	privateRouter.PUT(DEFAULT_API_PREFIX+"/ie/article/:id/reading", tc.ReGenArticleReading)
	privateRouter.GET(DEFAULT_API_PREFIX+"/ie/article/:id/proposed_vocab", tc.ExtractProposedWordsForArticle)

}
