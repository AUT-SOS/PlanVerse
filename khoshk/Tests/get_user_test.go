package Tests

import (
	"PlanVerse/controllers"
	"PlanVerse/models"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUser(t *testing.T) {
	server := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/get-user/10", nil)
	rec := httptest.NewRecorder()
	ctx := server.NewContext(req, rec)
	controllers.GetUserHandler(ctx)
	result := rec.Result()
	defer result.Body.Close()
	var user models.GetUserResponse
	json.Unmarshal(rec.Body.Bytes(), &user)
	assert.Equal(t, http.StatusOK, result.Status)
	assert.Equal(t, "arshiabp", user.Username)
}
