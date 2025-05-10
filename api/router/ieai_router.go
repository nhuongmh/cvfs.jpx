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

	privateRouter.POST(DEFAULT_API_PREFIX+"/ie/article/reading/:reading_id/submit", tc.SubmitReadingTest)
	privateRouter.GET(DEFAULT_API_PREFIX+"/ie/article/reading/:reading_id/submit", tc.GetTestSubmissionByReadingId)
	privateRouter.GET(DEFAULT_API_PREFIX+"/ie/article/reading/submit/:submit_id", tc.GetTestSubmission)
	privateRouter.DELETE(DEFAULT_API_PREFIX+"/ie/article/reading/submit/:submit_id", tc.DeleteTestSubmission)

	privateRouter.GET(DEFAULT_API_PREFIX+"/ie/article/:id/proposed_vocab", tc.ExtractProposedWordsForArticle)
	privateRouter.POST(DEFAULT_API_PREFIX+"/ie/article/:id/proposed_vocab", tc.HandleVocabProposalSubmit)
	privateRouter.GET(DEFAULT_API_PREFIX+"/ie/article/:id/vocab", tc.GetVocabListByArticleId)
	privateRouter.PUT(DEFAULT_API_PREFIX+"/ie/vocab/:id/anki", tc.GenAnkiDeckForVocabList)

}
