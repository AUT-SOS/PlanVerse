package helpers

import (
	"PlanVerse/configs"
	"PlanVerse/messages"
	"PlanVerse/models"
	"errors"
	"github.com/golang-jwt/jwt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const fixedLink = "http://localhost:5173/project"

func GenerateRandomCode() (string, error) {
	otp := ""
	rand.Seed(time.Now().UnixNano())
	random := 10000 + rand.Intn(89999)
	otp = strconv.Itoa(random)
	return otp, nil
}

func GenerateRandomLink(title string) (string, error) {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	code := make([]rune, 10)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}
	baseUrl, err := url.Parse(fixedLink)
	if err != nil {
		return "", err
	}
	params := url.Values{}
	params.Add(title, string(code))
	baseUrl.RawQuery = params.Encode()
	return baseUrl.String(), nil
}

func GenerateToken(userID int, duration time.Duration) (string, error) {
	claims := models.Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWTSecret")))
	if err != nil {
		log.Println("(GenerateToken) Error :", err)
		return "", err
	}
	return signedToken, nil
}

func CheckMail(email string) error {
	var users []models.User
	result := configs.DB.Select([]string{"id", "email", "is_verified"}).Find(&users)
	if result.Error != nil {
		return errors.New(strings.ToLower(messages.InternalError))
	}
	for _, otherUser := range users {
		if otherUser.Email == email {
			if otherUser.IsVerified {
				return errors.New(strings.ToLower(messages.DuplicateEmail))
			} else {
				configs.DB.Unscoped().Where("id = ?", otherUser.ID).Delete(&models.User{})
				return nil
			}
		}
	}
	return nil
}

func DetectMin(projectID int) (int, error) {
	var admins []Admin
	result := configs.DB.Table("projects_members").Select([]string{"promotion_time", "user_id"}).Where("project_id = ?", projectID).Scan(&admins)
	if result.Error != nil {
		return 0, errors.New(strings.ToLower(messages.InternalError))
	}
	if len(admins) == 0 {
		result = configs.DB.Table("projects").Where("id = ?", projectID).Update("members_number", 0)
		if result.Error != nil {
			return 0, errors.New(strings.ToLower(messages.InternalError))
		}
		result = configs.DB.Unscoped().Where("project_id = ?", projectID).Delete(&models.InvitedMembers{})
		if result.Error != nil {
			return 0, errors.New(strings.ToLower(messages.InternalError))
		}
		var stateIDs []int
		result = configs.DB.Table("states").Select("id").Where("project_id = ?", projectID).Scan(&stateIDs)
		if result.Error != nil {
			return 0, errors.New(strings.ToLower(messages.InternalError))
		}
		var wg *sync.WaitGroup
		for i := range stateIDs {
			wg.Add(1)
			go func(index int, wg *sync.WaitGroup) {
				defer wg.Done()
				var taskIDs []int
				result = configs.DB.Table("tasks").Select("id").Where("state_id = ?", stateIDs[index]).Scan(&taskIDs)
				for j := range taskIDs {
					configs.DB.Unscoped().Where("task_id = ?", taskIDs[j]).Delete(&models.TasksPerformers{})
				}
				result = configs.DB.Unscoped().Where("state_id = ?", stateIDs[index]).Delete(&models.Task{})
			}(i, wg)
		}
		wg.Wait()
		result = configs.DB.Unscoped().Where("project_id = ?", projectID).Delete(&models.State{})
		if result.Error != nil {
			return 0, errors.New(strings.ToLower(messages.InternalError))
		}
		result = configs.DB.Unscoped().Where("project_id = ?", projectID).Delete(&models.JoinLink{})
		if result.Error != nil {
			return 0, errors.New(strings.ToLower(messages.InternalError))
		}
		result = configs.DB.Unscoped().Where("id = ?", projectID).Delete(&models.Project{})
		if result.Error != nil {
			return 0, errors.New(strings.ToLower(messages.InternalError))
		}
		return 0, nil
	}
	minTime := Admin{
		UserID:        0,
		PromotionTime: time.Now().Add(time.Hour),
	}
	counter := 0
	for _, admin := range admins {
		if !admin.PromotionTime.Equal(time.Time{}) {
			counter++
			if admin.PromotionTime.Before(minTime.PromotionTime) {
				minTime = admin
			}
		}
	}
	if counter == 0 {
		minTime = admins[0]
		result = configs.DB.Table("projects_members").Where("project_id = ? and user_id = ?", projectID, minTime.UserID).Update("is_admin", true)
		if result.Error != nil {
			return 0, errors.New(strings.ToLower(messages.InternalError))
		}
		result = configs.DB.Table("projects_members").Where("project_id = ? and user_id = ?", projectID, minTime.UserID).Update("promotion_time", time.Now())
		if result.Error != nil {
			return 0, errors.New(strings.ToLower(messages.InternalError))
		}
		result = configs.DB.Table("projects").Where("id = ?", projectID).Update("members_number", len(admins))
		if result.Error != nil {
			return 0, errors.New(strings.ToLower(messages.InternalError))
		}
		return minTime.UserID, nil
	} else {
		result = configs.DB.Table("projects").Where("id = ?", projectID).Update("members_number", len(admins))
		if result.Error != nil {
			return 0, errors.New(strings.ToLower(messages.InternalError))
		}
		return minTime.UserID, nil
	}
}
