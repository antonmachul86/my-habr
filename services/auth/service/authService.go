package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"log"
	"my-habr/services/auth/model"
	"my-habr/services/auth/repository"
	"strconv"
	"time"
)

var (
	jwtKey             = []byte("PASSWORD!!!")
	refreshTokenExpriy = 7 * 24 * time.Hour
	acessTokenExpiry   = 15 * time.Second
)

type AuthService struct {
	userRepo *repository.UserRepository
	redis    *redis.Client
}

func NewAuthService(repo *repository.UserRepository, redis *redis.Client) *AuthService {
	return &AuthService{userRepo: repo, redis: redis}
}

func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *AuthService) CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *AuthService) GenerateTokens(userID int) (string, string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(acessTokenExpiry).Unix(),
	})
	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	refreshTokenBytes := make([]byte, 32)
	rand.Read(refreshTokenBytes)
	refreshToken := hex.EncodeToString(refreshTokenBytes)

	ctx := context.Background()
	s.redis.Set(ctx, refreshToken, userID, refreshTokenExpriy)

	return accessTokenString, refreshToken, nil
}

func (s *AuthService) Register(email, password string) error {
	hashed, err := s.HashPassword(password)
	if err != nil {
		return err
	}

	user := &model.User{Email: email, Password: hashed}
	return s.userRepo.CreateUser(context.Background(), user)

}
func (s *AuthService) Login(email, password string) (string, string, error) {
	log.Printf("🔍 Login attempt for email: %s", email)

	user, err := s.userRepo.FindByEmail(context.Background(), email)
	if err != nil {
		log.Printf("❌ DB error in FindByEmail: %v", err)
		return "", "", fmt.Errorf("database error")
	}

	if user == nil {
		log.Printf("❌ User not found: %s", email)
		return "", "", fmt.Errorf("invalid credentials")
	}

	log.Printf("✅ User found: ID=%d", user.ID)

	if !s.CheckPassword(user.Password, password) {
		log.Printf("❌ Invalid password for user: %s", email)
		return "", "", fmt.Errorf("invalid credentials")
	}

	log.Printf("✅ Password valid, generating tokens for user ID=%d", user.ID)
	return s.GenerateTokens(user.ID)
}

func (s *AuthService) Refresh(refreshToken string) (string, string, error) {
	ctx := context.Background()
	val, err := s.redis.Get(ctx, refreshToken).Result()
	if err != nil {
		return "", "", fmt.Errorf("invalid or expired refresh token")
	}

	s.redis.Del(ctx, refreshToken)

	userID, _ := strconv.Atoi(val)
	return s.GenerateTokens(userID)
}

func (s *AuthService) Logout(refreshToken string) error {
	ctx := context.Background()
	return s.redis.Del(ctx, refreshToken).Err()
}
