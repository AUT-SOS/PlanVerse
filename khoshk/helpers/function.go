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
