package controllers

import (
	"PlanVerse/configs"
	"PlanVerse/helpers"
	"PlanVerse/messages"
	"PlanVerse/models"
	"PlanVerse/services"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func RegisterHandler(ctx echo.Context) error {
	req := new(models.RegisterRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.FailedPasswordHashGeneration)
	}
	otp, err := helpers.GenerateRandomCode()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.FailedToCreateCode)
	}
	newUser := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}
	var users []models.User
	result := configs.DB.Find(&users)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	for _, otherUser := range users {
		if otherUser.Username == newUser.Username {
			return ctx.JSON(http.StatusNotAcceptable, messages.DuplicateUsername)
		} else if otherUser.Email == newUser.Email {
			return ctx.JSON(http.StatusNotAcceptable, messages.DuplicateEmail)
		}
	}
	err = configs.DB.Create(&newUser).Error
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.FailedToCreateUser)
	}
	configs.Redis.Set(configs.Ctx, strconv.Itoa(int(newUser.ID)), otp, time.Minute*5)
	err = services.SendMail("PlanVerse Verification", fmt.Sprintf("%s is your PlanVerse verification code", otp), []string{newUser.Email})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.FailedToSendEmail)
	}
	accessToken, err := helpers.GenerateToken(int(newUser.ID), time.Hour)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.FailedToCreateAccessToken)
	}
	refreshToken, err := helpers.GenerateToken(int(newUser.ID), time.Hour*24*7)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.FailedToCreateRefreshToken)
	}
	ctx.Response().Header().Set("Authorization", accessToken)
	cookie := &http.Cookie{
		Name:  "refresh_token",
		Value: refreshToken,
		Path:  "/refresh",
	}
	ctx.SetCookie(cookie)
	return ctx.JSON(http.StatusOK, messages.SentEmailSuccessfully)
}

func VerifyHandler(ctx echo.Context) error {
	req := new(models.VerifyRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}
	userIDCtx := ctx.Get("user_id")
	userID := userIDCtx.(int)
	val, err := configs.Redis.Get(configs.Ctx, strconv.Itoa(userID)).Result()
	if err != nil {
		configs.DB.Delete(&models.User{}, userID)
		return ctx.JSON(http.StatusBadRequest, messages.OTPExpired)
	}
	if val != req.OTP {
		return ctx.JSON(http.StatusBadRequest, messages.WrongOTP)
	}
	configs.Redis.Del(configs.Ctx, strconv.Itoa(userID))
	return ctx.JSON(http.StatusOK, messages.RegisteredSuccessfully)
}

func RefreshHandler(ctx echo.Context) error {
	cookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, messages.RefreshTokenExpired)
	}
	refreshToken, err := jwt.ParseWithClaims(cookie.Value, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWTSecret")), nil
	})
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, messages.InvalidRefreshToken)
	}
	if !refreshToken.Valid {
		deleteCookie := &http.Cookie{
			Name:    "refresh_token",
			Value:   "",
			Path:    "/refresh",
			Expires: time.Unix(0, 0),
		}
		ctx.SetCookie(deleteCookie)
		return ctx.JSON(http.StatusUnauthorized, messages.InvalidRefreshToken)
	}
	claims, _ := refreshToken.Claims.(*models.Claims)
	accessToken, err := helpers.GenerateToken(claims.UserID, time.Hour)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.FailedToCreateAccessToken)
	}
	ctx.Response().Header().Set("Authorization", accessToken)
	return ctx.JSON(http.StatusOK, messages.NewAccessToken)
}

func LoginHandler(ctx echo.Context) error {
	req := new(models.LoginRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}
	var user models.User
	result := configs.DB.Where("user_name = ?", req.Username).First(&user)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return ctx.JSON(http.StatusUnauthorized, messages.UsernameOrPasswordIncorrect)
	}
	accessToken, err := helpers.GenerateToken(int(user.ID), time.Hour)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.FailedToCreateAccessToken)
	}
	refreshToken, err := helpers.GenerateToken(int(user.ID), time.Hour*24*7)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.FailedToCreateRefreshToken)
	}
	ctx.Response().Header().Set("Authorization", accessToken)
	cookie := &http.Cookie{
		Name:  "refresh_token",
		Value: refreshToken,
		Path:  "/refresh",
	}
	ctx.SetCookie(cookie)
	return ctx.JSON(http.StatusOK, messages.LoggedInSuccessfully)
}
