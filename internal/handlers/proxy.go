package handlers

import (
	"api-gateway/internal/config"
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func ForwardToExamService(cfg *config.Config) gin.HandlerFunc {
	return func(context *gin.Context) {
		targetURL := cfg.EXAM_SERVICE_URL + context.Request.URL.Path

		newBody, err := io.ReadAll(context.Request.Body)
		if err != nil {
			context.JSON(500, gin.H{"error": "failed to read request body"})
			return
		}

		newReq, err := http.NewRequest(context.Request.Method, targetURL, bytes.NewReader(newBody))
		if err != nil {
			context.JSON(500, gin.H{"error": "failed to create request"})
			context.Abort()
			return
		}

		newReq.Header.Set("Content-Type", context.GetHeader("Content-Type"))
		newReq.Header.Set("x-user-id", context.GetString("user_id"))
		newReq.Header.Set("x-user-role", context.GetString("user_role"))

		client := &http.Client{}
		resp, err := client.Do(newReq)
		if err != nil {
			context.JSON(500, gin.H{"error": "request failed"})
			return
		}
		defer resp.Body.Close()

		// Copy response status
		context.Status(resp.StatusCode)

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				context.Header(key, value)
			}
		}

		// Copy response newBody
		io.Copy(context.Writer, resp.Body)
	}
}
