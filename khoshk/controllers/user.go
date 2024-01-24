package controllers

import (
	"PlanVerse/configs"
	"PlanVerse/helpers"
	"PlanVerse/messages"
	"PlanVerse/models"
	"PlanVerse/services"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(ctx echo.Context) error {
	req := new(models.RegisterRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}
	err := helpers.CheckMail(req.Email)
	if err != nil {
		if err.Error() == strings.ToLower(messages.InternalError) {
			return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
		} else if err.Error() == strings.ToLower(messages.DuplicateEmail) {
			return ctx.JSON(http.StatusNotAcceptable, messages.DuplicateEmail)
		}
	}
	otp, err := helpers.GenerateRandomCode()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.FailedToCreateCode)
	}
	go func() {
		services.SendMail("PlanVerse Verification", fmt.Sprintf("%s is your PlanVerse verification code", otp), []string{req.Email})
	}()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.FailedPasswordHashGeneration)
	}
	newUser := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}
	if err = configs.DB.Create(&newUser).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.FailedToCreateUser)
	}
	configs.Redis.Set(configs.Ctx, strconv.Itoa(int(newUser.ID)), otp, time.Minute*5)
	accessToken, err := helpers.GenerateToken(int(newUser.ID), time.Hour)
	if err != nil {
		configs.DB.Unscoped().Where("id = ?", newUser.ID).Delete(&models.User{})
		return ctx.JSON(http.StatusInternalServerError, messages.FailedToCreateAccessToken)
	}
	refreshToken, err := helpers.GenerateToken(int(newUser.ID), time.Hour*24*7)
	if err != nil {
		configs.DB.Unscoped().Where("id = ?", newUser.ID).Delete(&models.User{})
		return ctx.JSON(http.StatusInternalServerError, messages.FailedToCreateRefreshToken)
	}
	ctx.Response().Header().Set("Authorization", accessToken)
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/refresh",
		HttpOnly: true,
	}
	ctx.SetCookie(cookie)
	return ctx.JSON(http.StatusOK, messages.AddedToDatabase)
}

func VerifyHandler(ctx echo.Context) error {
	req := new(models.VerifyRequest)
	res := new(models.Response)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}
	userIDCtx := ctx.Get("user_id")
	userID := userIDCtx.(int)
	val, err := configs.Redis.Get(configs.Ctx, strconv.Itoa(userID)).Result()
	if err != nil {
		configs.DB.Unscoped().Where("id = ?", userID).Delete(&models.User{})
		return ctx.JSON(http.StatusBadRequest, messages.OTPExpired)
	}
	if val != req.OTP {
		return ctx.JSON(http.StatusBadRequest, messages.WrongOTP)
	}
	result := configs.DB.Model(&models.User{}).Where("id = ?", userID).Update("is_verified", true)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	configs.Redis.Del(configs.Ctx, strconv.Itoa(userID))
	res.UserID = userID
	res.Message = messages.RegisteredSuccessfully
	return ctx.JSON(http.StatusOK, res)
}

func RefreshHandler(ctx echo.Context) error {
	cookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, messages.RefreshTokenExpired)
	}
	refreshToken, err := jwt.ParseWithClaims(cookie.Value, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWTSecret")), nil
	})
	if err != nil || !refreshToken.Valid {
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
	res := new(models.Response)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}
	var user models.User
	result := configs.DB.Select([]string{"id", "password", "is_verified"}).Where("email = ?", req.Email).Find(&user)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return ctx.JSON(http.StatusUnauthorized, messages.EmailOrPasswordIncorrect)
	}
	if !user.IsVerified {
		return ctx.JSON(http.StatusUnauthorized, messages.UserNotVerified)
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
	res.UserID = int(user.ID)
	res.Message = messages.LoggedInSuccessfully
	return ctx.JSON(http.StatusOK, res)
}

func GetUserHandler(ctx echo.Context) error {
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongUserID)
	}
	var user models.User
	result := configs.DB.Where("id = ?", userID).Find(&user)
	if result.Error != nil {
		return ctx.JSON(http.StatusBadRequest, messages.UnknownUser)
	}
	return ctx.JSON(http.StatusOK, user)
}
