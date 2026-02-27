package presenter

import (
	"app/domain/entities"
	"app/usecase"
)

type userSignupPresenter struct{}

func NewUserSignupPresenter() usecase.UserSignupPresenter {
	return &userSignupPresenter{}
}

func (p *userSignupPresenter) Output(user *entities.User, token string) *usecase.UserSignupOutput {
	return &usecase.UserSignupOutput{
		ID:   int(user.ID),
		Token: token,
	}
}
