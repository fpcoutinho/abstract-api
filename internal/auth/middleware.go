package auth

import (
	"context"
	"net/http"
	"os"
	"sync"

	firebase "firebase.google.com/go/v4"
	firebaseAuth "firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

var (
	client     *firebaseAuth.Client
	clientOnce sync.Once
)

func ensureClient() error {
	var err error
	clientOnce.Do(func() {
		cred := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		if cred == "" {
			err = nil
			return
		}
		app, appErr := firebase.NewApp(context.Background(), nil, option.WithCredentialsFile(cred))
		if appErr != nil {
			err = appErr
			return
		}
		c, authErr := app.Auth(context.Background())
		if authErr != nil {
			err = authErr
			return
		}
		client = c
	})
	return err
}

// FirebaseOrMock returns a middleware that supports a mock mode and real Firebase verification
func FirebaseOrMock() gin.HandlerFunc {
	mode := os.Getenv("AUTH_MODE")
	return func(c *gin.Context) {
		if mode == "mock" {
			uid := c.GetHeader("X-User-Uid")
			if uid == "" {
				auth := c.GetHeader("Authorization")
				if len(auth) > 7 && auth[:7] == "Bearer " {
					uid = auth[7:]
				}
			}
			if uid == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing mock uid (set X-User-Uid or send Bearer <uid>)"})
				return
			}
			c.Set("uid", uid)
			c.Next()
			return
		}

		// real mode: ensure firebase client (requires GOOGLE_APPLICATION_CREDENTIALS)
		if err := ensureClient(); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "firebase init: " + err.Error()})
			return
		}
		if client == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "firebase credentials not configured; set GOOGLE_APPLICATION_CREDENTIALS or use AUTH_MODE=mock"})
			return
		}

		authHeader := c.GetHeader("Authorization")
		if len(authHeader) <= 7 || authHeader[:7] != "Bearer " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Bearer token"})
			return
		}
		idToken := authHeader[7:]
		token, err := client.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token: " + err.Error()})
			return
		}
		c.Set("uid", token.UID)
		c.Next()
	}
}
