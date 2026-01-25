package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/eenemeene/kitamanager-go/internal/apperror"
)

// parseID extracts and validates ID from URL parameter
func parseID(c *gin.Context, param string) (uint, error) {
	id, err := strconv.ParseUint(c.Param(param), 10, 32)
	if err != nil {
		return 0, apperror.BadRequest("invalid " + param)
	}
	return uint(id), nil
}

// respondError sends consistent error response
func respondError(c *gin.Context, err error) {
	code := apperror.HTTPStatus(err)
	c.JSON(code, gin.H{"error": err.Error()})
}
