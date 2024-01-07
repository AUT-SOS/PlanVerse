package middlewares

import (
	"PlanVerse/Tar/messages"
	"PlanVerse/Tar/models"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if ctx.Path() == "/register" || ctx.Path() == "/login" {
			return next(ctx)
		}
		tokenString := ctx.Request().Header.Get("Authorization")
		if tokenString == "" {
			return ctx.JSON(http.StatusUnauthorized, messages.Unauthorized)
		}
		token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return os.Getenv("JWTSecret"), nil
		})
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, messages.Unauthorized)
		}
		if !token.Valid {
			return ctx.JSON(http.StatusUnauthorized, messages.Unauthorized)
		}
		claims, _ := token.Claims.(*models.Claims)
		ctx.Set("user_id", claims.UserID)
		return next(ctx)
	}
}
