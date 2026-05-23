package handlers

import (
	"net/http"
	"strconv"

	"abstract-api/internal/store"
	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	uidVal, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}

	profile, err := store.GetProfile(c.Request.Context(), uidVal.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

func PatchProfile(c *gin.Context) {
	uidVal, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}

	var input store.ProfileUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := store.UpdateProfile(c.Request.Context(), uidVal.(string), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

func PatchAvatar(c *gin.Context) {
	uidVal, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}

	var input struct {
		AvatarIdx int `json:"avatarIdx"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := store.UpdateAvatar(c.Request.Context(), uidVal.(string), input.AvatarIdx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

func ParseIntParam(c *gin.Context, name string) (int, error) {
	return strconv.Atoi(c.Param(name))
}
