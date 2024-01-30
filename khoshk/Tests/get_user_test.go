package Tests

import (
	"PlanVerse/controllers"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUser(t *testing.T) {
	server := echo.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/get-my-user", nil)
	ctx := server.NewContext(req, rec)
	ctx.Set("user_id", 89)
	t.Run("should return 200 status ok and 89 for id", func(t *testing.T) {
		controllers.GetUserIDHandler(ctx)
		result := rec.Result()
		defer result.Body.Close()
		var userID int
		json.Unmarshal(rec.Body.Bytes(), &userID)
		assert.Equal(t, fmt.Sprint("200 OK"), result.Status)
		assert.Equal(t, 89, userID)
	})
}
