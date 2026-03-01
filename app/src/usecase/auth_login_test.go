package usecase_test

import (
	"context"
	"errors"
	"testing"

	"app/src/domain/entities"
	"app/src/domain/value_objects"
	"app/src/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type AuthLoginMockPresenter struct{}

func (p *AuthLoginMockPresenter) Output(token string) *usecase.AuthLoginOutput {
	return &usecase.AuthLoginOutput{Token: token}
}

type AuthLoginMockUserRepository struct{ mock.Mock }

func (m *AuthLoginMockUserRepository) FindByUsername(ctx context.Context, username string) (*entities.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *AuthLoginMockUserRepository) Create(user *entities.User) (*entities.User, error) {
	return nil, nil
}

func (m *AuthLoginMockUserRepository) FindByID(id value_objects.ID) (*entities.User, error) {
	return nil, nil
}

func (m *AuthLoginMockUserRepository) FindByEmail(email value_objects.Email) (*entities.User, error) {
	return nil, nil
}

func (m *AuthLoginMockUserRepository) Update(user *entities.User) error {
	return nil
}

func (m *AuthLoginMockUserRepository) Delete(id value_objects.ID) error {
	return nil
}

type AuthLoginMockAuthService struct{ mock.Mock }

func (m *AuthLoginMockAuthService) VerifyPassword(hashed value_objects.HashedPassword, rawPassword string) error {
	args := m.Called(hashed, rawPassword)
	return args.Error(0)
}

func (m *AuthLoginMockAuthService) IssueToken(ctx context.Context, user *entities.User) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *AuthLoginMockAuthService) HashPassword(password string) (string, error) {
	return "", nil
}

func (m *AuthLoginMockAuthService) VerifyToken(ctx context.Context, token string) (*entities.User, error) {
	return nil, nil
}

func TestAuthLogin_Execute(t *testing.T) {
	bob, _ := entities.NewUser(1, "local_bob", "bob@example.com", "hashed_password")
	hashedPass := bob.HashedPassword

	tests := []struct {
		name      string
		input     usecase.AuthLoginInput
		setupMock func(am *AuthLoginMockAuthService, ur *AuthLoginMockUserRepository)
		wantErr   bool
	}{
		{
			name:  "【正常系】正しいパスワードでログインに成功し、トークンが返る",
			input: usecase.AuthLoginInput{Username: "local_bob", Password: "correct_password"},
			setupMock: func(am *AuthLoginMockAuthService, ur *AuthLoginMockUserRepository) {
				ur.On("FindByUsername", mock.Anything, "local_bob").Return(bob, nil)
				am.On("VerifyPassword", hashedPass, "correct_password").Return(nil)
				am.On("IssueToken", mock.Anything, bob).Return("valid_jwt_token", nil)
			},
			wantErr: false,
		},
		{
			name:  "【異常系】ユーザーが存在しない場合、エラーを返す",
			input: usecase.AuthLoginInput{Username: "none_user", Password: "any"},
			setupMock: func(am *AuthLoginMockAuthService, ur *AuthLoginMockUserRepository) {
				ur.On("FindByUsername", mock.Anything, "none_user").Return((*entities.User)(nil), errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name:  "【異常系】パスワードが間違っている場合、認証失敗",
			input: usecase.AuthLoginInput{Username: "local_bob", Password: "wrong_password"},
			setupMock: func(am *AuthLoginMockAuthService, ur *AuthLoginMockUserRepository) {
				ur.On("FindByUsername", mock.Anything, "local_bob").Return(bob, nil)
				am.On("VerifyPassword", hashedPass, "wrong_password").Return(errors.New("invalid password"))
			},
			wantErr: true,
		},
		{
			name:  "【異常系】トークン生成に失敗した場合",
			input: usecase.AuthLoginInput{Username: "local_bob", Password: "correct_password"},
			setupMock: func(am *AuthLoginMockAuthService, ur *AuthLoginMockUserRepository) {
				ur.On("FindByUsername", mock.Anything, "local_bob").Return(bob, nil)
				am.On("VerifyPassword", hashedPass, "correct_password").Return(nil)
				am.On("IssueToken", mock.Anything, bob).Return("", errors.New("token creation failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am := new(AuthLoginMockAuthService)
			ur := new(AuthLoginMockUserRepository)
			presenter := &AuthLoginMockPresenter{}
			tt.setupMock(am, ur)

			interactor := usecase.NewAuthLoginInteractor(presenter, ur, am)
			output, err := interactor.Execute(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, output)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, "valid_jwt_token", output.Token)
			}

			am.AssertExpectations(t)
			ur.AssertExpectations(t)
		})
	}
}
