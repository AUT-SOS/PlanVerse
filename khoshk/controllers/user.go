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
	"strings"
	"time"
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
	go func() {
		services.SendMail("PlanVerse Verification", fmt.Sprintf("%s is your PlanVerse verification code", otp), []string{req.Email})
	}()
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
	res := new(models.UserResponse)
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
	res := new(models.UserResponse)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}
	var user models.User
	result := configs.DB.Select([]string{"id", "password", "is_verified"}).Where("email = ?", req.Email).Find(&user)
	if result.Error != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongEmail)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return ctx.JSON(http.StatusNotAcceptable, messages.PasswordIncorrect)
	}
	if !user.IsVerified {
		return ctx.JSON(http.StatusNotAcceptable, messages.UserNotVerified)
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

func ResendEmailHandler(ctx echo.Context) error {
	userIDCtx := ctx.Get("user_id")
	userID := userIDCtx.(int)
	var user models.User
	result := configs.DB.Select("email").Where("id = ?", userID).Find(&user)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	otp, err := helpers.GenerateRandomCode()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.FailedToCreateCode)
	}
	configs.Redis.Set(configs.Ctx, strconv.Itoa(userID), otp, time.Minute*5)
	go func() {
		services.SendMail("PlanVerse Verification", fmt.Sprintf("%s is your PlanVerse verification code", otp), []string{user.Email})
	}()
	return ctx.JSON(http.StatusOK, messages.SentEmail)
}

func GetUserHandler(ctx echo.Context) error {
	res := new(models.GetUserResponse)
	userID, err := strconv.Atoi(ctx.Param("user-id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.WrongUserID)
	}
	result := configs.DB.Table("users").Select([]string{"username", "email", "profile_pic"}).Where("id = ?", userID).Find(res)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	res.ID = userID
	return ctx.JSON(http.StatusOK, res)
}

func GetUserIDHandler(ctx echo.Context) error {
	userIDCtx := ctx.Get("user_id")
	userID := userIDCtx.(int)
	return ctx.JSON(http.StatusOK, userID)
}

func EditUserHandler(ctx echo.Context) error {
	req := new(models.EditUserRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, messages.InvalidRequestBody)
	}
	userIDCtx := ctx.Get("user_id")
	userID := userIDCtx.(int)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.FailedPasswordHashGeneration)
	}
	var user models.User
	result := configs.DB.Where("id = ?", userID).Find(&user)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	user.Username = req.Username
	user.Password = string(hashedPassword)
	user.Email = req.Email
	user.ProfilePic = req.ProfilePic
	result = configs.DB.Save(&user)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	return ctx.JSON(http.StatusOK, messages.UserEdited)
}

func DeleteUserHandler(ctx echo.Context) error {
	userIDCtx := ctx.Get("user_id")
	userID := userIDCtx.(int)
	var user models.User
	result := configs.DB.Where("id = ?", uint(userID)).Preload("Projects").Find(&user)
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	projects := user.Projects
	result = configs.DB.Unscoped().Where("user_id = ?", userID).Delete(&models.ProjectsMembers{})
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	result = configs.DB.Unscoped().Where("id = ?", userID).Delete(&models.User{})
	if result.Error != nil {
		return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
	}
	for _, project := range projects {
		if project.OwnerID == userID {
			newOwnerID, err := helpers.DetectMin(int(project.ID))
			if err != nil {
				return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
			}
			if newOwnerID != 0 {
				result = configs.DB.Table("projects").Where("id = ?", project.ID).Update("owner_id", newOwnerID)
				if result.Error != nil {
					return ctx.JSON(http.StatusInternalServerError, messages.InternalError)
				}
			}
		}
	}
	return ctx.JSON(http.StatusOK, messages.UserAccountDeleted)
}
