package services

import (
	"context"
	"app/src/domain/entities"
	"app/src/domain/value_objects" // 追加：型定義のため
)

type AuthDomainService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword value_objects.HashedPassword, rawPassword string) error
	IssueToken(ctx context.Context, user *entities.User) (string, error)
	VerifyToken(ctx context.Context, token string) (*entities.User, error)
}