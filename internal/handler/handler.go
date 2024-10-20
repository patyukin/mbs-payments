package handler

import (
	"auth-telegram/internal/handler/dto"
	"auth-telegram/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

const (
	HeaderAuthorization = "Authorization"
	HeaderUserID        = "User-ID"
	HeaderUserRole      = "User-Role"
)

type UseCase interface {
	SignUp(ctx context.Context, loginData model.SignUpData) (dto.SignUpResponse, error)
	SignUpV2(ctx context.Context, in model.SignUpV2Data) (dto.SignUpV2Response, error)
	ResendCode(ctx context.Context, loginData model.SignUpData) (dto.SignUpResponse, error)
	SignIn(ctx context.Context, signInData model.SignInData) error
	SignInV2(ctx context.Context, signInData model.SignInData) (dto.TokensResponse, error)
	SignInVerify(ctx context.Context, code string) (dto.TokensResponse, error)
	GetUserAuthInfoByToken(ctx context.Context, id string) (dto.UserAuthInfo, error)
	GetUserFullInfo(ctx context.Context, userID uuid.UUID) (dto.User, error)
	GetUserInfoByUUID(ctx context.Context, userID uuid.UUID) (dto.User, error)
	GenerateTokens(ctx context.Context, refreshToken string) (dto.TokensResponse, error)
	GetTelegramBot() string
	GetJWTToken() []byte
	GetTokenByName(ctx context.Context, name string) (string, error)
	ValidateToken(token string) (string, error)
}

type Handler struct {
	uc UseCase
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func New(uc UseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) HandleError(w http.ResponseWriter, code int, message string) {
	log.Error().Msgf("Error: %s", message)

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(ErrorResponse{Error: message})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get(HeaderAuthorization)
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is missing")
	}

	// Проверяем, что заголовок начинается с "Bearer "
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", fmt.Errorf("authorization header does not start with 'Bearer '")
	}

	// Извлекаем токен, удаляя префикс "Bearer "
	token := strings.TrimPrefix(authHeader, bearerPrefix)
	return token, nil
}
