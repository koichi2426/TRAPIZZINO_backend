package domain_impl_services

import (
	"context"
	"app/domain/entities"
	"app/domain/services"
)

// AuthDomainServiceImplはAuthDomainServiceインターフェースの具象実装です。
type AuthDomainServiceImpl struct{}

func NewAuthDomainServiceImpl() services.AuthDomainService {
	return &AuthDomainServiceImpl{}
}

// IssueTokenはユーザー情報からトークンを発行します（ダミー実装）。
func (s *AuthDomainServiceImpl) IssueToken(ctx context.Context, user *entities.User) (string, error) {
	return "dummy-token", nil
}

// VerifyTokenはトークンを検証しユーザー情報を返します（ダミー実装）。
func (s *AuthDomainServiceImpl) VerifyToken(ctx context.Context, token string) (*entities.User, error) {
	return &entities.User{ID: 1, Username: "dummy", Email: "dummy@example.com", HashedPassword: "hashed"}, nil
}
