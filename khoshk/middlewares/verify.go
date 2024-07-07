package middlewares

import (
	"PlanVerse/configs"
	"PlanVerse/messages"
	"PlanVerse/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func VerifyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if ctx.Path() == "/register" || ctx.Path() == "/login" || ctx.Path() == "/refresh" || ctx.Path() == "/verify" || ctx.Path() == "/resend-email" || ctx.Path() == "/metrics" || ctx.Path() == "/ws" {
			return next(ctx)
		}
		userIDCtx := ctx.Get("user_id")
		userID := userIDCtx.(int)
		var user models.User
		result := configs.DB.Select([]string{"is_verified"}).Where("id = ?", userID).Find(&user)
		if result.Error != nil {
			return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
		}
		if !user.IsVerified {
			return ctx.JSON(http.StatusNotAcceptable, messages.UserNotVerified)
		}
		return next(ctx)
	}
}
