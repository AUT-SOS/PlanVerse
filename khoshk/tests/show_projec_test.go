package tests

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

func TestShowProject(t *testing.T) {
	t.Run("should return 200 status ok and the project with id = 40 as response", func(t *testing.T) {
		reqBody := models.ShowProjectRequest{
			Link: "http://localhost:5173/project?test-project1=bATiZewQLT",
		}
		buf, _ := json.Marshal(reqBody)
		req, errReq := http.NewRequest(http.MethodPost, "http://localhost:8080/show-project", bytes.NewBuffer(buf))
		req.Header.Set("Authorization", auth)
		res, errRes := http.DefaultClient.Do(req)
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)
		var resBody models.ShowProjectResponse
		json.Unmarshal(body, &resBody)
		assert.NoError(t, errReq)
		assert.NoError(t, errRes)
		assert.Equal(t, fmt.Sprint("200 OK"), res.Status)
		assert.Equal(t, models.ShowProjectResponse{
			ProjectID:     40,
			Title:         "test-project1",
			BackGroundPic: "pic1",
			MembersNumber: 1,
			Members: []models.MemberInfo{
				{Username: "test-user", ProfilePic: ""},
			},
		}, resBody)
	})
	t.Run("should return 401 status unauthorized and message: Unauthorized", func(t *testing.T) {
		reqBody := models.ShowProjectRequest{
			Link: "http://localhost:5173/project?test-project1=wtgaQmAPSH",
		}
		buf, _ := json.Marshal(reqBody)
		req, errReq := http.NewRequest(http.MethodPost, "http://localhost:8080/show-project", bytes.NewBuffer(buf))
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
