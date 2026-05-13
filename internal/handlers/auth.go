package handlers

import (
	"api-gateway/internal/config"
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func ForwardToAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetURL := cfg.AUTH_SERVICE_URL + c.Request.URL.Path

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to read request body",
			})
			return
		}

		req, err := http.NewRequest(
			c.Request.Method,
			targetURL,
			bytes.NewBuffer(body),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create request",
			})
			return
		}

		req.Header = c.Request.Header.Clone()

		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error": "auth service unavailable",
			})
			return
		}
		defer resp.Body.Close()

		for key, values := range resp.Header {
			for _, value := range values {
				c.Writer.Header().Add(key, value)
			}
		}

		c.Status(resp.StatusCode)

		io.Copy(c.Writer, resp.Body)
	}
}
