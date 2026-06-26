package middleware

import (
	"log"
	"net/http"
	"time"

	"manifest-api/dto"

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
