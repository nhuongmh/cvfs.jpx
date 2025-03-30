package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nhuongmh/cfvs.jpx/api/router"
	"github.com/nhuongmh/cfvs.jpx/bootstrap"
	"github.com/nhuongmh/cfvs.jpx/pkg/logger"
)

func main() {
	logger.InitLog()
	app := bootstrap.Init()

	timeout := time.Duration(app.Env.ContextTimeout) * time.Second
	gine := gin.Default()

	logger.Log.Info().Msg("Setting up router...")
	gine.Use(CORSMiddleware())
	router.SetupPostgres(&app, timeout, gine)

	logger.Log.Info().Msg("Starting server...")
	gine.Run(app.Env.ServerAddress)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
