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

type UserSignupMockPresenter struct{}

func (p *UserSignupMockPresenter) Output(user *entities.User, token string) *usecase.UserSignupOutput {
	return &usecase.UserSignupOutput{
		ID:    user.ID.Value(),
		Token: token,
	}
}

type UserSignupMockUserRepository struct{ mock.Mock }

func (m *UserSignupMockUserRepository) Create(user *entities.User) (*entities.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *UserSignupMockUserRepository) FindByID(id value_objects.ID) (*entities.User, error) {
	return nil, nil
}

func (m *UserSignupMockUserRepository) FindByEmail(email value_objects.Email) (*entities.User, error) {
	return nil, nil
}

func (m *UserSignupMockUserRepository) FindByUsername(ctx context.Context, username string) (*entities.User, error) {
	return nil, nil
}

func (m *UserSignupMockUserRepository) Update(user *entities.User) error {
	return nil
}

func (m *UserSignupMockUserRepository) Delete(id value_objects.ID) error {
	return nil
}

type UserSignupMockAuthService struct{ mock.Mock }

func (m *UserSignupMockAuthService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *UserSignupMockAuthService) IssueToken(ctx context.Context, user *entities.User) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *UserSignupMockAuthService) VerifyPassword(hashed value_objects.HashedPassword, rawPassword string) error {
	return nil
}

func (m *UserSignupMockAuthService) VerifyToken(ctx context.Context, token string) (*entities.User, error) {
	return nil, nil
}

func TestUserSignup_Execute(t *testing.T) {
	tests := []struct {
		name      string
		input     usecase.UserSignupInput
		setupMock func(am *UserSignupMockAuthService, ur *UserSignupMockUserRepository)
		wantErr   bool
		check     func(t *testing.T, out *usecase.UserSignupOutput)
	}{
		{
			name: "【正常系】新規ユーザーが登録されトークンが返る",
			input: usecase.UserSignupInput{
				Username: "new_malloy",
				Email:    "malloy@example.com",
				Password: "password123",
			},
			setupMock: func(am *UserSignupMockAuthService, ur *UserSignupMockUserRepository) {
				am.On("HashPassword", "password123").Return("hashed_password_abc", nil)
				ur.On("Create", mock.MatchedBy(func(u *entities.User) bool {
					return u.Username.String() == "new_malloy" &&
						u.Email.String() == "malloy@example.com" &&
						u.HashedPassword.String() == "hashed_password_abc"
				})).Return(func(u *entities.User) *entities.User {
					u.ID, _ = value_objects.NewID(10)
					return u
				}(mustNewUser(t, 0, "new_malloy", "malloy@example.com", "hashed_password_abc")), nil)
				am.On("IssueToken", mock.Anything, mock.AnythingOfType("*entities.User")).Return("welcome_jwt_token", nil)
			},
			wantErr: false,
			check: func(t *testing.T, out *usecase.UserSignupOutput) {
				assert.Equal(t, 10, out.ID)
				assert.Equal(t, "welcome_jwt_token", out.Token)
			},
		},
		{
			name: "【異常系】パスワードハッシュ化に失敗した場合",
			input: usecase.UserSignupInput{Username: "malloy", Email: "malloy@example.com", Password: "password123"},
			setupMock: func(am *UserSignupMockAuthService, ur *UserSignupMockUserRepository) {
				am.On("HashPassword", "password123").Return("", errors.New("hash failed"))
			},
			wantErr: true,
		},
		{
			name: "【異常系】ユーザー保存に失敗した場合",
			input: usecase.UserSignupInput{Username: "malloy", Email: "malloy@example.com", Password: "password123"},
			setupMock: func(am *UserSignupMockAuthService, ur *UserSignupMockUserRepository) {
				am.On("HashPassword", "password123").Return("hashed_password_abc", nil)
				ur.On("Create", mock.Anything).Return((*entities.User)(nil), errors.New("db save error"))
			},
			wantErr: true,
		},
		{
			name: "【異常系】トークン生成に失敗した場合",
			input: usecase.UserSignupInput{Username: "malloy", Email: "malloy@example.com", Password: "password123"},
			setupMock: func(am *UserSignupMockAuthService, ur *UserSignupMockUserRepository) {
				am.On("HashPassword", "password123").Return("hashed_password_abc", nil)
				created := mustNewUser(t, 7, "malloy", "malloy@example.com", "hashed_password_abc")
				ur.On("Create", mock.Anything).Return(created, nil)
				am.On("IssueToken", mock.Anything, created).Return("", errors.New("token creation failed"))
			},
			wantErr: true,
		},
		{
			name: "【異常系】入力メールが不正な場合",
			input: usecase.UserSignupInput{Username: "malloy", Email: "invalid-email", Password: "password123"},
			setupMock: func(am *UserSignupMockAuthService, ur *UserSignupMockUserRepository) {
				am.On("HashPassword", "password123").Return("hashed_password_abc", nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am := new(UserSignupMockAuthService)
			ur := new(UserSignupMockUserRepository)
			presenter := &UserSignupMockPresenter{}
			tt.setupMock(am, ur)

			interactor := usecase.NewUserSignupInteractor(presenter, ur, am)
			out, err := interactor.Execute(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, out)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, out)
				if tt.check != nil {
					tt.check(t, out)
				}
			}

			am.AssertExpectations(t)
			ur.AssertExpectations(t)
		})
	}
}

func mustNewUser(t *testing.T, id int, username, email, hashed string) *entities.User {
	t.Helper()
	u, err := entities.NewUser(id, username, email, hashed)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	return u
}
