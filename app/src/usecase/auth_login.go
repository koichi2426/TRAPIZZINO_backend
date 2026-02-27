package usecase

import (
	"context"
	"app/domain/services"
)

type AuthLoginInput struct {
	Username string
	Password string
}

type AuthLoginOutput struct {
	Token string
}

type AuthLoginPresenter interface {
	Output(token string) *AuthLoginOutput
}

type AuthLoginUseCase interface {
	Execute(ctx context.Context, input AuthLoginInput) (*AuthLoginOutput, error)
}

type authLoginInteractor struct {
	presenter   AuthLoginPresenter
	userService services.AuthDomainService
}

func NewAuthLoginInteractor(p AuthLoginPresenter, s services.AuthDomainService) AuthLoginUseCase {
	return &authLoginInteractor{
		presenter:   p,
		userService: s,
	}
}

func (i *authLoginInteractor) Execute(ctx context.Context, input AuthLoginInput) (*AuthLoginOutput, error) {
	user, err := i.userService.VerifyToken(ctx, input.Username+":"+input.Password)
	if err != nil {
		return nil, err
	}
	token, err := i.userService.IssueToken(ctx, user)
	if err != nil {
		return nil, err
	}
	return i.presenter.Output(token), nil
}
