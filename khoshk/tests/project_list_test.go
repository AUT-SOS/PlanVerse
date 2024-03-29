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

func TestProjectList(t *testing.T) {
	req, errReq := http.NewRequest(http.MethodGet, "http://localhost:8080/list-project", nil)
	t.Run("should return 200 status ok and projects with 40 and 41 ids", func(t *testing.T) {
		req.Header.Set("Authorization", auth)
		res, errRes := http.DefaultClient.Do(req)
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)
		var projectsList []models.ProjectListResponse
		json.Unmarshal(body, &projectsList)
		assert.NoError(t, errReq)
		assert.NoError(t, errRes)
		assert.Equal(t, fmt.Sprint("200 OK"), res.Status)
		assert.Equal(t, []models.ProjectListResponse{
			{ID: 40, Title: "test-project1", BackGroundPic: "pic1", MembersNumber: 1, IsAdmin: true},
			{ID: 41, Title: "test-project2", BackGroundPic: "pic2", MembersNumber: 1, IsAdmin: true},
		}, projectsList)
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
