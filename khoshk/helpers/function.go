package helpers

import (
	"PlanVerse/configs"
	"PlanVerse/messages"
	"PlanVerse/models"
	"errors"
	"github.com/golang-jwt/jwt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func GenerateRandomCode() (string, error) {
	otp := ""
	rand.Seed(time.Now().UnixNano())
	random := 10000 + rand.Intn(89999)
	otp = strconv.Itoa(random)
	return otp, nil
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
