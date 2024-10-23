package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nhuongmh/cfvs.jpx/bootstrap"
	"github.com/nhuongmh/cfvs.jpx/pkg/service/jpxpractice/repo"
)

const (
	API_V1             = "v1"
	DEFAULT_API_PREFIX = "/api/" + API_V1
)

func Setup(app *bootstrap.Application, timeout time.Duration, gine *gin.Engine) {
	publicRouter := gine.Group("public")
	privateRouter := gine.Group("private")

	publicRouter.Use(cors.Default())
	publicRouter.Static("/data", "./data")

	tr := repo.NewJpxPraticeRepo(app.DB)
	NewJpxServiceRouter(app, tr, timeout, publicRouter, privateRouter)
}
