package domain_impl_services

import (
	"context"
	"errors"
	"time"

	"app/domain/entities"
	"app/domain/services"
	"app/domain/value_objects" // 追加：型変換に必要

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthDomainServiceImpl struct {
	secretKey []byte
}

// NewAuthDomainServiceImpl は環境変数から渡された秘密鍵で初期化します
func NewAuthDomainServiceImpl(secret string) services.AuthDomainService {
	return &AuthDomainServiceImpl{
		secretKey: []byte(secret),
	}
}

// HashPassword は bcrypt でパスワードをハッシュ化します
func (s *AuthDomainServiceImpl) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword はハッシュと生のパスワードを比較します
func (s *AuthDomainServiceImpl) VerifyPassword(hashedPassword, rawPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword))
}

// IssueToken は JWT を発行します
func (s *AuthDomainServiceImpl) IssueToken(ctx context.Context, user *entities.User) (string, error) {
	if user == nil {
		return "", errors.New("user is nil")
	}

	claims := jwt.MapClaims{
		"user_id":  int(user.ID),           // 独自型を int にキャストして保存
		"username": string(user.Username), // 独自型を string にキャストして保存
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// VerifyToken は JWT を検証し、ドメインモデル(User)を復元します
func (s *AuthDomainServiceImpl) VerifyToken(ctx context.Context, tokenString string) (*entities.User, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userID, okID := claims["user_id"].(float64)
		username, okName := claims["username"].(string)

		if !okID || !okName {
			return nil, errors.New("failed to parse token claims")
		}

		// 重要：基本型からドメインの独自型（Value Object）にキャストして戻す
		return &entities.User{
			ID:       value_objects.ID(int(userID)),
			Username: value_objects.Username(username),
		}, nil
	}

	return nil, errors.New("could not parse claims")
}