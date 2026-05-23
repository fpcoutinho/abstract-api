package server

import (
	"github.com/gin-gonic/gin"

	"abstract-api/internal/auth"
	"abstract-api/internal/handlers"
)

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

	// protected routes
	protected := v1.Group("")
	protected.Use(auth.FirebaseOrMock())
	protected.GET("/profile", handlers.GetProfile)
	protected.PATCH("/profile", handlers.PatchProfile)
	protected.PATCH("/profile/avatar", handlers.PatchAvatar)

	v1.GET("/trails", handlers.ListTrails)
	v1.GET("/trails/:trailId", handlers.GetTrail)
	v1.GET("/trails/:trailId/missions", handlers.ListTrailMissions)
	v1.GET("/missions/:missionId", handlers.GetMission)

	return engine
}
