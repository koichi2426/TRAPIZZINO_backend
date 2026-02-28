package domain_impl_services

import (
	"context"
	"errors"
	"time"

	"app/domain/entities"
	"app/domain/services"
	"app/domain/value_objects"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthDomainServiceImpl struct {
	secretKey []byte
}

func NewAuthDomainServiceImpl(secret string) services.AuthDomainService {
	return &AuthDomainServiceImpl{
		secretKey: []byte(secret),
	}
}

func (s *AuthDomainServiceImpl) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// 修正: 引数をインターフェース定義（VO）に合わせ、内部で String() を呼ぶ
func (s *AuthDomainServiceImpl) VerifyPassword(hashedPassword value_objects.HashedPassword, rawPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword.String()), []byte(rawPassword))
}

func (s *AuthDomainServiceImpl) IssueToken(ctx context.Context, user *entities.User) (string, error) {
	if user == nil {
		return "", errors.New("user is nil")
	}

	// 修正: VO から生の値を取り出すメソッド (.Value() や .String()) を使用
	claims := jwt.MapClaims{
		"user_id":  user.ID.Value(),       
		"username": user.Username.String(), 
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

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
		userIDFloat, okID := claims["user_id"].(float64) // JWT の数値は float64 でパースされる
		usernameStr, okName := claims["username"].(string)

		if !okID || !okName {
			return nil, errors.New("failed to parse token claims")
		}

		// 修正: NewID 等のファクトリメソッドを使用して VO を再生成する
		uID, _ := value_objects.NewID(int(userIDFloat))
		uName, _ := value_objects.NewUsername(usernameStr)

		return &entities.User{
			ID:       uID,
			Username: uName,
		}, nil
	}

	return nil, errors.New("could not parse claims")
}