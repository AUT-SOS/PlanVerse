package middlewares

import (
	"PlanVerse/configs"
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
		var projectIDs []int
		result := configs.DB.Table("projects_members").Select("project_id").Where("user_id = ?", userID).Scan(&projectIDs)
		if result.Error != nil {
			return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
		}
		if len(projectIDs) == 0 {
			return ctx.JSON(http.StatusNotAcceptable, messages.UserNoProject)
		}
		exist := false
		for _, id := range projectIDs {
			if id == projectID {
				exist = true
				break
			}
		}
		if !exist {
			return ctx.JSON(http.StatusNotAcceptable, messages.NotMember)
		}
		var isAdmin bool
		result = configs.DB.Table("projects_members").Select("is_admin").Where("project_id = ? and user_id = ?", projectID, userID).Scan(&isAdmin)
		if result.Error != nil {
			return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
		}
		if !isAdmin {
			return ctx.JSON(http.StatusUnauthorized, messages.AdminAccess)
		}
		return next(ctx)
	}
}
