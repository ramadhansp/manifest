package middleware

import (
	"log"
	"net/http"
	"time"

	"manifest-api/dto"
	"manifest-api/utils"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		latency := time.Since(t)
		log.Printf("method=%s path=%s status=%d latency=%s", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), latency)
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, dto.APIResponse{
					Success: false,
					Message: "Internal Server Error",
					Errors:  []string{"A critical error occurred"},
				})
			}
		}()
		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: "Unauthorized",
				Errors:  []string{"Missing or invalid token"},
			})
			return
		}
		
		tokenString := authHeader[7:]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.APIResponse{
				Success: false,
				Message: "Unauthorized",
				Errors:  []string{err.Error()},
			})
			return
		}
		
		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])
		c.Next()
	}
}
