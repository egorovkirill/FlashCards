package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"Error": "Empty authorization header",
		})
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"Error": "Invalid authorization header",
		})
		return
	}

	if len(headerParts[1]) == 0 {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"Error": "Token is empty",
		})
		return
	}

	userId, err := h.service.Authenthication.ParseToken(headerParts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"Error": err.Error(),
		})
		return
	}

	c.Set(userCtx, userId)
}
