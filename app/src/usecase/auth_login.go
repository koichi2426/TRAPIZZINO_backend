package usecase

import (
	"context"
	"errors"
	"app/src/domain/entities"
	"app/src/domain/services"
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
	userRepo    entities.UserRepository 
	userService services.AuthDomainService
}

func NewAuthLoginInteractor(
	p AuthLoginPresenter, 
	r entities.UserRepository, 
	s services.AuthDomainService,
) AuthLoginUseCase {
	return &authLoginInteractor{
		presenter:   p,
		userRepo:    r,
		userService: s,
	}
}

func (i *authLoginInteractor) Execute(ctx context.Context, input AuthLoginInput) (*AuthLoginOutput, error) {
	// 1. DBからユーザーを取得する
	user, err := i.userRepo.FindByUsername(ctx, input.Username)
	if err != nil {
		// セキュリティのため、ユーザーの存在有無を特定させないメッセージを返す
		return nil, errors.New("invalid username or password")
	}

	// 2. パスワードを照合する
	// 【修正】インターフェースの変更（VOを受け取る形式）に合わせてキャストを削除
	err = i.userService.VerifyPassword(user.HashedPassword, input.Password)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// 3. 照合成功！JWT トークンを発行する
	token, err := i.userService.IssueToken(ctx, user)
	if err != nil {
		return nil, err
	}

	return i.presenter.Output(token), nil
}