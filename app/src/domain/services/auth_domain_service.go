package services

import (
	"context"
	"src/domain/entities"
)

type AuthDomainService interface {
	IssueToken(ctx context.Context, user *entities.User) (string, error)
	VerifyToken(ctx context.Context, token string) (*entities.User, error)
}
