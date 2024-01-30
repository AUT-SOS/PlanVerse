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
	t.Run("should return 200 status ok and 93 as id and User profile edited successfully", func(t *testing.T) {
		reqBody := models.EditUserRequest{
			Username:   "arshiabp",
			Email:      "arshia.bahar@gmail.com",
			Password:   "Ab25121381",
			ProfilePic: "",
		}
		buf, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/edit-profile", bytes.NewBuffer(buf))
		req.Header.Set("Authorization", auth)
		res, _ := http.DefaultClient.Do(req)
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)
		var resBody string
		json.Unmarshal(body, &resBody)
		assert.Equal(t, fmt.Sprint("200 OK"), res.Status)
		assert.Equal(t, messages.UserEdited, resBody)
	})
}
