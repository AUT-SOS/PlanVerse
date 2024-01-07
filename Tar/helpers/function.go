package helpers

import (
	"PlanVerse/Tar/configs"
	"PlanVerse/Tar/models"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"math/rand"
	"os"
	"strconv"
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
	var users []models.User
	result := configs.DB.Find(&users)
	if result.Error != nil {
		return "", result.Error
	}
	for _, u := range users {
		val, err := configs.Redis.Get(configs.Ctx, strconv.Itoa(int(u.ID))).Result()
		if err != nil {
			return "", err
		}
		if val == otp {
			return "", errors.New("this otp is used")
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
	signedToken, err := token.SignedString(os.Getenv("JWTSecret"))
	if err != nil {
		log.Println("(GenerateToken) Error :", err)
		return "", err
	}
	return signedToken, nil
}
