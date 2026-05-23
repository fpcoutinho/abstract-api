package handlers

import (
	"errors"
	"net/http"

	"abstract-api/internal/store"
	"github.com/gin-gonic/gin"
)

func ListTrails(c *gin.Context) {
	trails, err := store.ListTrails(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"trails": trails})
}

func GetTrail(c *gin.Context) {
	trailID := c.Param("trailId")
	trail, err := store.GetTrail(c.Request.Context(), trailID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, store.ErrNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, trail)
}

func ListTrailMissions(c *gin.Context) {
	trailID := c.Param("trailId")
	missions, err := store.ListTrailMissions(c.Request.Context(), trailID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"missions": missions})
}

func GetMission(c *gin.Context) {
	missionID := c.Param("missionId")
	mission, err := store.GetMission(c.Request.Context(), missionID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, store.ErrNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mission)
}
