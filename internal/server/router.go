package server

import "github.com/gin-gonic/gin"

func New() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	engine.GET("/healthz", func(context *gin.Context) {
		context.JSON(200, gin.H{"status": "ok"})
	})

	v1 := engine.Group("/api/v1")
	v1.GET("/health", func(context *gin.Context) {
		context.JSON(200, gin.H{"status": "ok"})
	})

	return engine
}