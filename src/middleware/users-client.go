package middleware

import (
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type userService struct {
	restyClient *resty.Client
}

func NewUserService() UserService {
	restyClient := resty.New()
	USER_API_URL := os.Getenv("USER_API_URL")
	if USER_API_URL == "" {
		USER_API_URL = "http://127.0.0.1:8080"
	}
	restyClient.SetBaseURL(USER_API_URL).
		SetHeader("Content-Type", "application/json").
		SetTimeout(5 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second)

	slog.Info("Initialized User Service Client", "baseURL", restyClient.BaseURL)
	return &userService{restyClient}
}

type UserService interface {
	GetUser(userId uint) User
}

func (us *userService) GetUser(userId uint) User {
	slog.Info("Fetching user by ID from the database", "id", userId)
	var user User
	resp, err := us.restyClient.R().
		SetResult(&user).
		Get("/users/" + strconv.Itoa(int(userId)))
	if err != nil {
		slog.Error("Error fetching user from User Service", "error", err)
	}
	if resp.StatusCode() != 200 {
		slog.Warn("User not found in User Service", "userId", userId, "statusCode", resp.StatusCode())
	}
	slog.Info("Fetched user from User Service", "user", user)
	return user
}
