package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"app/src/domain/entities"
	"app/src/domain/value_objects"
	"app/src/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- MOCK 定義 ---

type MockAuthService struct{ mock.Mock }
func (m *MockAuthService) VerifyToken(ctx context.Context, token string) (*entities.User, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*entities.User), args.Error(1)
}
func (m *MockAuthService) HashPassword(p string) (string, error) { return "", nil }
func (m *MockAuthService) IssueToken(ctx context.Context, u *entities.User) (string, error) { return "", nil }
func (m *MockAuthService) VerifyPassword(h value_objects.HashedPassword, p string) error { return nil }

type MockSpotRepository struct{ mock.Mock }
func (m *MockSpotRepository) FindByLocation(ctx context.Context, lat, lon float64) (*entities.Spot, error) {
	args := m.Called(ctx, lat, lon)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*entities.Spot), args.Error(1)
}
func (m *MockSpotRepository) Create(s *entities.Spot) (*entities.Spot, error) {
	args := m.Called(s)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*entities.Spot), args.Error(1)
}
func (m *MockSpotRepository) FindByID(ctx context.Context, id value_objects.ID) (*entities.Spot, error) { return nil, nil }
func (m *MockSpotRepository) FindByMeshID(mID value_objects.MeshID) ([]*entities.Spot, error) { return nil, nil }
func (m *MockSpotRepository) Update(s *entities.Spot) error { return nil }
func (m *MockSpotRepository) Delete(id value_objects.ID) error { return nil }
func (m *MockSpotRepository) FindResonantUsersWithMatchCount(ctx context.Context, uID value_objects.ID) ([]entities.ResonantUser, error) { return nil, nil }
func (m *MockSpotRepository) FindSpotsByMeshAndUsers(ctx context.Context, mIDs []value_objects.MeshID, uIDs []value_objects.ID) ([]*entities.Spot, error) { return nil, nil }
func (m *MockSpotRepository) GetDensityScoreByMesh(ctx context.Context, mID value_objects.MeshID) (value_objects.DensityScore, error) { return value_objects.NewDensityScore(0) }
func (m *MockSpotRepository) FindPostsBySpot(ctx context.Context, sID value_objects.ID) ([]*entities.Post, error) { return nil, nil }

type MockPostRepository struct{ mock.Mock }
func (m *MockPostRepository) Create(p *entities.Post) (*entities.Post, error) {
	args := m.Called(p)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*entities.Post), args.Error(1)
}
func (m *MockPostRepository) FindBySpotID(sID value_objects.ID) ([]*entities.Post, error) { return nil, nil }
func (m *MockPostRepository) FindByID(id value_objects.ID) (*entities.Post, error) { return nil, nil }
func (m *MockPostRepository) Update(p *entities.Post) error { return nil }
func (m *MockPostRepository) Delete(id value_objects.ID) error { return nil }

type MockPresenter struct{}
func (p *MockPresenter) Output(s *entities.Spot, post *entities.Post) *usecase.RegisterSpotPostOutput {
	return &usecase.RegisterSpotPostOutput{SpotID: s.ID.Value(), PostID: post.ID.Value()}
}

// --- TEST 本体 ---

func TestRegisterSpotPost_Execute(t *testing.T) {
	// 共通データ
	malloy, _ := entities.NewUser(2, "local_malloy", "malloy@example.com", "hashed_password")
	existingSpot, _ := entities.NewSpot(1, "恵比寿うどん", 35.6467, 139.7101, 1)
	newlyCreatedSpot, _ := entities.NewSpot(99, "新規店", 35.0, 135.0, 2)
	dummyPost, _ := entities.NewPost(100, 2, 1, "local_malloy", "http://example.com/post.jpg", "caption", time.Now())

	tests := []struct {
		name      string
		input     usecase.RegisterSpotPostInput
		setupMock func(am *MockAuthService, sm *MockSpotRepository, pm *MockPostRepository)
		wantErr   bool
		check     func(t *testing.T, out *usecase.RegisterSpotPostOutput)
	}{
		{
			name: "【正常系】既存スポットがある場合、自動的に合流する（Merge）",
			input: usecase.RegisterSpotPostInput{
				Token: "valid_token", Latitude: 35.6467, Longitude: 139.7101, ImageURL: "http://example.com/merge.jpg", Caption: "合流！",
			},
			setupMock: func(am *MockAuthService, sm *MockSpotRepository, pm *MockPostRepository) {
				am.On("VerifyToken", mock.Anything, "valid_token").Return(malloy, nil)
				sm.On("FindByLocation", mock.Anything, 35.6467, 139.7101).Return(existingSpot, nil)
				pm.On("Create", mock.MatchedBy(func(p *entities.Post) bool {
					return p.SpotID.Value() == 1 // 既存IDが使われているか検証
				})).Return(dummyPost, nil)
			},
			wantErr: false,
			check: func(t *testing.T, out *usecase.RegisterSpotPostOutput) {
				assert.Equal(t, 1, out.SpotID)
			},
		},
		{
			name: "【正常系】新規地点の場合、スポットを新規作成して投稿する",
			input: usecase.RegisterSpotPostInput{
				Token: "valid_token", SpotName: "新規店", Latitude: 35.0, Longitude: 135.0, ImageURL: "http://example.com/new.jpg",
			},
			setupMock: func(am *MockAuthService, sm *MockSpotRepository, pm *MockPostRepository) {
				am.On("VerifyToken", mock.Anything, "valid_token").Return(malloy, nil)
				sm.On("FindByLocation", mock.Anything, 35.0, 135.0).Return((*entities.Spot)(nil), nil)
				sm.On("Create", mock.Anything).Return(newlyCreatedSpot, nil)
				pm.On("Create", mock.Anything).Return(dummyPost, nil)
			},
			wantErr: false,
			check: func(t *testing.T, out *usecase.RegisterSpotPostOutput) {
				assert.Equal(t, 99, out.SpotID)
			},
		},
		{
			name: "【異常系】トークンが不正な場合、エラーを返す",
			input: usecase.RegisterSpotPostInput{Token: "invalid_token"},
			setupMock: func(am *MockAuthService, sm *MockSpotRepository, pm *MockPostRepository) {
				am.On("VerifyToken", mock.Anything, "invalid_token").Return((*entities.User)(nil), errors.New("unauthorized"))
			},
			wantErr: true,
		},
		{
			name: "【異常系】Spot検索時にDBエラーが発生した場合",
			input: usecase.RegisterSpotPostInput{Token: "valid_token", Latitude: 35.6, Longitude: 139.7},
			setupMock: func(am *MockAuthService, sm *MockSpotRepository, pm *MockPostRepository) {
				am.On("VerifyToken", mock.Anything, "valid_token").Return(malloy, nil)
				sm.On("FindByLocation", mock.Anything, 35.6, 139.7).Return((*entities.Spot)(nil), errors.New("db find error"))
			},
			wantErr: true,
		},
		{
			name: "【異常系】Post保存時にDBエラーが発生した場合",
			input: usecase.RegisterSpotPostInput{Token: "valid_token", Latitude: 35.6, Longitude: 139.7, ImageURL: "http://example.com/error.jpg"},
			setupMock: func(am *MockAuthService, sm *MockSpotRepository, pm *MockPostRepository) {
				am.On("VerifyToken", mock.Anything, "valid_token").Return(malloy, nil)
				sm.On("FindByLocation", mock.Anything, 35.6, 139.7).Return(existingSpot, nil)
				pm.On("Create", mock.Anything).Return((*entities.Post)(nil), errors.New("db insert error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am, sm, pm := new(MockAuthService), new(MockSpotRepository), new(MockPostRepository)
			tt.setupMock(am, sm, pm)
			interactor := usecase.NewRegisterSpotPostInteractor(&MockPresenter{}, sm, pm, am)

			out, err := interactor.Execute(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.check != nil {
					tt.check(t, out)
				}
			}
			am.AssertExpectations(t)
			sm.AssertExpectations(t)
			pm.AssertExpectations(t)
		})
	}
}