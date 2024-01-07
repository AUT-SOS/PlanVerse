package controllers

import (
	"PlanVerse/Tar/configs"
	"PlanVerse/Tar/helpers"
	"PlanVerse/Tar/messages"
	"PlanVerse/Tar/models"
	"PlanVerse/Tar/services"
	"fmt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
)

func RegisterHandler(ctx echo.Context) error {
	req := models.RegisterRequest{}
	err := ctx.Bind(&req)
	if err != nil {
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
		Username:  req.Username,
		Password:  string(hashedPassword),
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
	return ctx.JSON(http.StatusOK, messages.SentEmailSuccessfully)
}
