package usecase

import (
	"context"
	"app/domain/entities"
	"app/domain/services"
)

type UserSignupInput struct {
	Username string
	Email    string
	Password string
}

type UserSignupOutput struct {
	ID    int
	Token string
}

type UserSignupPresenter interface {
	Output(user *entities.User, token string) *UserSignupOutput
}

type UserSignupUseCase interface {
	Execute(ctx context.Context, input UserSignupInput) (*UserSignupOutput, error)
}

type userSignupInteractor struct {
	presenter   UserSignupPresenter
	userRepo    entities.UserRepository
	userService services.AuthDomainService
}

func NewUserSignupInteractor(p UserSignupPresenter, r entities.UserRepository, s services.AuthDomainService) UserSignupUseCase {
	return &userSignupInteractor{
		presenter:   p,
		userRepo:    r,
		userService: s,
	}
}

func (i *userSignupInteractor) Execute(ctx context.Context, input UserSignupInput) (*UserSignupOutput, error) {
	// 1. パスワードをハッシュ化（bcrypt変換）する
	hashedPassword, err := i.userService.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// 2. ハッシュ化されたパスワードを使ってユーザーエンティティを作成する
	user, err := entities.NewUser(0, input.Username, input.Email, hashedPassword)
	if err != nil {
		return nil, err
	}

	// 3. データベースに保存する（Repository を経由）
	created, err := i.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	// 4. JWT トークンを発行する
	token, err := i.userService.IssueToken(ctx, created)
	if err != nil {
		return nil, err
	}

	return i.presenter.Output(created, token), nil
}