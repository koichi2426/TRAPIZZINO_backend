package usecase_test

import (
	"context"
	"testing"
	"time"

	"app/src/domain/entities"
	"app/src/domain/value_objects"
	"app/src/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/mock"
)

// --- MOCK 定義 ---

type MockAuthService struct {
	mock.Mock
}
func (m *MockAuthService) VerifyToken(ctx context.Context, token string) (*entities.User, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*entities.User), args.Error(1)
}
func (m *MockAuthService) HashPassword(password string) (string, error) { return "", nil }
func (m *MockAuthService) IssueToken(ctx context.Context, user *entities.User) (string, error) { return "", nil }
func (m *MockAuthService) VerifyPassword(hashed value_objects.HashedPassword, password string) error {
	return nil
}

type MockSpotRepository struct {
	mock.Mock
}
func (m *MockSpotRepository) FindByLocation(ctx context.Context, lat, lon float64) (*entities.Spot, error) {
	args := m.Called(ctx, lat, lon)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*entities.Spot), args.Error(1)
}
func (m *MockSpotRepository) Create(s *entities.Spot) (*entities.Spot, error) {
	args := m.Called(s)
	return args.Get(0).(*entities.Spot), args.Error(1)
}
func (m *MockSpotRepository) FindByID(ctx context.Context, id value_objects.ID) (*entities.Spot, error) {
	return nil, nil
}
func (m *MockSpotRepository) FindByMeshID(meshID value_objects.MeshID) ([]*entities.Spot, error) {
	return nil, nil
}
func (m *MockSpotRepository) Update(spot *entities.Spot) error {
	return nil
}
func (m *MockSpotRepository) Delete(id value_objects.ID) error {
	return nil
}
func (m *MockSpotRepository) FindResonantUsersWithMatchCount(ctx context.Context, userID value_objects.ID) ([]entities.ResonantUser, error) {
	return nil, nil
}
func (m *MockSpotRepository) FindSpotsByMeshAndUsers(ctx context.Context, meshIDs []value_objects.MeshID, userIDs []value_objects.ID) ([]*entities.Spot, error) {
	return nil, nil
}
func (m *MockSpotRepository) GetDensityScoreByMesh(ctx context.Context, meshID value_objects.MeshID) (value_objects.DensityScore, error) {
	return value_objects.NewDensityScore(0)
}
func (m *MockSpotRepository) FindPostsBySpot(ctx context.Context, spotID value_objects.ID) ([]*entities.Post, error) {
	return nil, nil
}

type MockPostRepository struct {
	mock.Mock
}
func (m *MockPostRepository) Create(p *entities.Post) (*entities.Post, error) {
	args := m.Called(p)
	return args.Get(0).(*entities.Post), args.Error(1)
}
func (m *MockPostRepository) FindBySpotID(spotID value_objects.ID) ([]*entities.Post, error) {
	return nil, nil
}
func (m *MockPostRepository) FindByID(id value_objects.ID) (*entities.Post, error) {
	return nil, nil
}
func (m *MockPostRepository) Update(post *entities.Post) error {
	return nil
}
func (m *MockPostRepository) Delete(id value_objects.ID) error {
	return nil
}

type MockPresenter struct{}
func (p *MockPresenter) Output(s *entities.Spot, post *entities.Post) *usecase.RegisterSpotPostOutput {
	return &usecase.RegisterSpotPostOutput{
		SpotID: s.ID.Value(),
		PostID: post.ID.Value(),
	}
}

// --- TEST 本体 ---

func TestRegisterSpotPost_Execute(t *testing.T) {
	t.Run("【正常系】既存スポットがある場合、同じSpotIDに自動合流する", func(t *testing.T) {
		ctx := context.Background()
		authMock := new(MockAuthService)
		spotMock := new(MockSpotRepository)
		postMock := new(MockPostRepository)
		presenter := &MockPresenter{}

		interactor := usecase.NewRegisterSpotPostInteractor(presenter, spotMock, postMock, authMock)

		// データ準備
		dummyUser, err := entities.NewUser(2, "local_malloy", "malloy@example.com", "hashed_password")
		require.NoError(t, err)
		dummySpot, err := entities.NewSpot(1, "恵比寿うどん", 35.6467, 139.7101, 1)
		require.NoError(t, err)
		dummyPost, err := entities.NewPost(100, 2, 1, "local_malloy", "http://ex.com/malloy.jpg", "合流！", time.Now())
		require.NoError(t, err)

		// モック設定
		authMock.On("VerifyToken", ctx, "valid_token").Return(dummyUser, nil)
		spotMock.On("FindByLocation", ctx, 35.6467, 139.7101).Return(dummySpot, nil)
		postMock.On("Create", mock.Anything).Return(dummyPost, nil)

		input := usecase.RegisterSpotPostInput{
			Token:     "valid_token",
			SpotName:  "恵比寿うどん",
			Latitude:  35.6467,
			Longitude: 139.7101,
			ImageURL:  "http://ex.com/malloy.jpg",
			Caption:   "合流！",
		}
		
		output, err := interactor.Execute(ctx, input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, 1, output.SpotID)
		assert.Equal(t, 100, output.PostID)
		
		authMock.AssertExpectations(t)
		spotMock.AssertExpectations(t)
		postMock.AssertExpectations(t)
	})
}