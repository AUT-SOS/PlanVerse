package Tests

import (
	"PlanVerse/messages"
	"PlanVerse/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestVerifyEmail(t *testing.T) {
	t.Run("should return 200 status ok and User profile edited successfully", func(t *testing.T) {
		reqBody := models.EditUserRequest{
			Username:   "test-user",
			Email:      "example.planverse@gmail.com",
			Password:   "Ab25121381",
			ProfilePic: "",
		}
		buf, _ := json.Marshal(reqBody)
		req, errReq := http.NewRequest(http.MethodPost, "http://localhost:8080/edit-profile", bytes.NewBuffer(buf))
		req.Header.Set("Authorization", auth)
		res, errRes := http.DefaultClient.Do(req)
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)
		var resBody string
		json.Unmarshal(body, &resBody)
		assert.NoError(t, errReq)
		assert.NoError(t, errRes)
		assert.Equal(t, fmt.Sprint("200 OK"), res.Status)
		assert.Equal(t, messages.UserEdited, resBody)
	})
}
