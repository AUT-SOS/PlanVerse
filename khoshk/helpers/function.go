package helpers

import (
	"PlanVerse/configs"
	"PlanVerse/messages"
	"PlanVerse/models"
	"errors"
	"fmt"
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
	for i := 0; i < 5; i++ {
		random := rand.Intn(9 - 1)
		randomDigit := strconv.Itoa(random)
		otp = fmt.Sprint(otp + randomDigit)
	}
	keys, _, err := configs.Redis.Scan(configs.Ctx, 0, "", 1000).Result()
	if err != nil {
		return "", err
	}
	for _, key := range keys {
		val, newErr := configs.Redis.Get(configs.Ctx, key).Result()
		if newErr != nil {
			return "", newErr
		}
		if val == otp {
			return "", errors.New(strings.ToLower(messages.RepeatedOTP))
		}
	}
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

func CheckDuplicate(email string) error {
	var users []models.User
	result := configs.DB.Select([]string{"email", "is_verified"}).Find(&users)
	if result.Error != nil {
		return errors.New(strings.ToLower(messages.InternalError))
	}
	for _, otherUser := range users {
		if otherUser.Email == email {
			if otherUser.IsVerified {
				return errors.New(strings.ToLower(messages.DuplicateEmail))
			} else {
				configs.DB.Raw("delete from users where id = ?", otherUser.ID)
				break
			}
		}
	}
	return nil
}
