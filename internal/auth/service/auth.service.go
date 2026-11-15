package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"eci-backend/internal/auth/model"
	"eci-backend/internal/platform/security"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"log"
	"strings"
	"time"

	"eci-backend/internal/auth/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("email atau password salah")
)

type AuthService struct {
	userRepository    *repository.UserRepository
	sessionRepository *repository.SessionRepository
	turnstileService  *security.TurnstileService
}

const sessionDuration = 24 * time.Hour

func NewAuthService(
	userRepository *repository.UserRepository,
	turnstileService *security.TurnstileService,
	sessionRepository *repository.SessionRepository,
) *AuthService {
	return &AuthService{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		turnstileService:  turnstileService,
	}
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
	turnstileToken string,
) (*model.User, string, time.Time, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	user, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, "", time.Time{}, ErrInvalidCredentials
		}

		log.Println("LOGIN DEBUG: database error:", err)
		return nil, "", time.Time{}, err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	); err != nil {
		log.Println("LOGIN DEBUG: password mismatch")
		return nil, "", time.Time{}, ErrInvalidCredentials
	}

	if err := s.turnstileService.Verify(ctx, turnstileToken); err != nil {
		log.Println("LOGIN DEBUG: turnstile failed:", err)
		return nil, "", time.Time{}, err
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, "", time.Time{}, err
	}

	token := base64.RawURLEncoding.EncodeToString(tokenBytes)

	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	expiresAt := time.Now().Add(sessionDuration)

	session := &model.Session{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	if err := s.sessionRepository.Create(ctx, session); err != nil {
		return nil, "", time.Time{}, err
	}

	return user, token, expiresAt, nil
}

func (s *AuthService) Logout(
	ctx context.Context,
	token string,
) error {
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	return s.sessionRepository.DeleteByTokenHash(ctx, tokenHash)
}

func (s *AuthService) GetProfile(
	ctx context.Context,
	userID string,
) (*model.User, error) {
	id, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	return s.userRepository.FindByID(ctx, id)
}
