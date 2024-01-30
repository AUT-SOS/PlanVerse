package Tests

import (
	"PlanVerse/models"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {
	t.Run("should return 200 status ok and user with id = 93 as response", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/get-user/93", nil)
		req.Header.Set("Authorization", auth)
		res, _ := http.DefaultClient.Do(req)
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)
		var user models.GetUserResponse
		json.Unmarshal(body, &user)
		assert.Equal(t, fmt.Sprint("200 OK"), res.Status)
		assert.Equal(t, models.GetUserResponse{
			ID:         93,
			Username:   "arshiabp",
			Email:      "arshia.bahar@gmail.com",
			ProfilePic: "",
		}, user)
	})
}
