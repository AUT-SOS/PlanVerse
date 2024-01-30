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

func TestProjectList(t *testing.T) {
	t.Run("should return 200 status ok and user with id = 93 as response", func(t *testing.T) {
		req, errReq := http.NewRequest(http.MethodGet, "http://localhost:8080/list-project", nil)
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
			{ID: 37, Title: "test-project1", BackGroundPic: "pic1", MembersNumber: 1, IsAdmin: true},
			{ID: 38, Title: "test-project2", BackGroundPic: "pic2", MembersNumber: 1, IsAdmin: true},
		}, projectsList)
	})
}
