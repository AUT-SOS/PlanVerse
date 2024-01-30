package tests

import (
	"PlanVerse/messages"
	"PlanVerse/models"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {
	req, errReq := http.NewRequest(http.MethodGet, "http://localhost:8080/get-user/97", nil)
	t.Run("should return 200 status ok and the user with id = 97 as response", func(t *testing.T) {
		req.Header.Set("Authorization", auth)
		res, errRes := http.DefaultClient.Do(req)
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)
		var resBody models.GetUserResponse
		json.Unmarshal(body, &resBody)
		assert.NoError(t, errReq)
		assert.NoError(t, errRes)
		assert.Equal(t, fmt.Sprint("200 OK"), res.Status)
		assert.Equal(t, models.GetUserResponse{
			ID:         97,
			Username:   "test-user",
			Email:      "example.planverse@gmail.com",
			ProfilePic: "",
		}, resBody)
	})
	t.Run("should return 401 status unauthorized and message: Unauthorized", func(t *testing.T) {
		res, errRes := http.DefaultClient.Do(req)
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)
		var resBody string
		json.Unmarshal(body, &resBody)
		assert.NoError(t, errReq)
		assert.NoError(t, errRes)
		assert.Equal(t, fmt.Sprint("401 Unauthorized"), res.Status)
		assert.Equal(t, messages.Unauthorized, resBody)
	})
}
