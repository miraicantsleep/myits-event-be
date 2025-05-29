package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/miraicantsleep/myits-event-be/dto"
	"github.com/miraicantsleep/myits-event-be/service"
	"github.com/miraicantsleep/myits-event-be/utils"
)

func Authenticate(jwtService service.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.MESSAGE_FAILED_TOKEN_NOT_FOUND, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		if !strings.Contains(authHeader, "Bearer ") {
			response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.MESSAGE_FAILED_TOKEN_NOT_VALID, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		authHeader = strings.Replace(authHeader, "Bearer ", "", -1)
		token, err := jwtService.ValidateToken(authHeader)
		if err != nil {
			response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.MESSAGE_FAILED_TOKEN_NOT_VALID, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		if !token.Valid {
			response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.MESSAGE_FAILED_DENIED_ACCESS, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userId, err := jwtService.GetUserIDByToken(authHeader)
		if err != nil {
			response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		role, err := jwtService.GetRoleByToken(authHeader)
		if err != nil {
			response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		ctx.Set("token", authHeader)
		ctx.Set("user_id", userId)
		ctx.Set("role", role)
		ctx.Next()
	}
}

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole, exists := ctx.Get("role")
		if !exists {
			response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, "failed to get user role", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		if userRole != requiredRole {
			response := utils.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, "insufficient permission", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		ctx.Next()
	}
}
