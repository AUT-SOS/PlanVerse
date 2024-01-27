package middlewares

import (
	"PlanVerse/configs"
	"PlanVerse/helpers"
	"PlanVerse/messages"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		projectID, err := strconv.Atoi(ctx.Param("project-id"))
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, messages.WrongProjectID)
		}
		userIDCtx := ctx.Get("user_id")
		userID := userIDCtx.(int)
		var showRole helpers.ShowRole
		result := configs.DB.Table("projects_members").Select("is_admin").Where("project_id = ? and user_id = ?", projectID, userID).Scan(&showRole)
		if result.Error != nil {
			return ctx.JSON(http.StatusNotAcceptable, messages.NotMember)
		}
		if !showRole.IsAdmin {
			return ctx.JSON(http.StatusUnauthorized, messages.AdminAccess)
		}
		return next(ctx)
	}
}
