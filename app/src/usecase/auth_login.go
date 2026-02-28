package usecase

import (
	"context"
	"errors"
	"app/domain/entities"
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
	userRepo    entities.UserRepository // 追加：DBからユーザーを探すために必要
	userService services.AuthDomainService
}

// コンストラクタも userRepo を受け取るように修正
func NewAuthLoginInteractor(p AuthLoginPresenter, r entities.UserRepository, s services.AuthDomainService) AuthLoginUseCase {
	return &authLoginInteractor{
		presenter:   p,
		userRepo:    r,
		userService: s,
	}
}

func (i *authLoginInteractor) Execute(ctx context.Context, input AuthLoginInput) (*AuthLoginOutput, error) {
	// 1. DBからユーザーを取得する（ユーザー名で検索）
	user, err := i.userRepo.FindByUsername(ctx, input.Username)
	if err != nil {
		// ユーザーが見つからない場合も「ユーザー名またはパスワードが違います」と出すのがセキュリティの鉄則
		return nil, errors.New("invalid username or password")
	}

	// 2. パスワードを照合する（生の入力 vs DBのハッシュ）
	// userService.VerifyPassword は一致しない場合にエラーを返します
	err = i.userService.VerifyPassword(string(user.HashedPassword), input.Password)
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