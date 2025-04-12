package handlers

import (
	"github.com/gin-gonic/gin"
)

func init() {
	// Set Gin to release mode for tests
	gin.SetMode(gin.ReleaseMode)
}
