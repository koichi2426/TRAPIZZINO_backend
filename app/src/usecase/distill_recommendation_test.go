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

// --- MOCK 定義 ---

type MockRecommendationService struct{ mock.Mock }

type DistillMockAuthService struct{ mock.Mock }

func (m *DistillMockAuthService) VerifyToken(ctx context.Context, token string) (*entities.User, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *DistillMockAuthService) HashPassword(password string) (string, error) {
	return "", nil
}

func (m *DistillMockAuthService) VerifyPassword(hashed value_objects.HashedPassword, rawPassword string) error {
	return nil
}

func (m *DistillMockAuthService) IssueToken(ctx context.Context, user *entities.User) (string, error) {
	return "", nil
}

func (m *MockRecommendationService) Distill(ctx context.Context, u *entities.User, lat value_objects.Latitude, lng value_objects.Longitude) (*entities.Spot, value_objects.TotalScore, value_objects.ResonanceCount, value_objects.DensityScore, value_objects.Reason, []*entities.Post, error) {
	args := m.Called(ctx, u, lat, lng)
	// 戻り値が多いので、それぞれ慎重にキャストして返します
	spot, _ := args.Get(0).(*entities.Spot)
	total, _ := args.Get(1).(value_objects.TotalScore)
	res, _ := args.Get(2).(value_objects.ResonanceCount)
	den, _ := args.Get(3).(value_objects.DensityScore)
	reason, _ := args.Get(4).(value_objects.Reason)
	posts, _ := args.Get(5).([]*entities.Post)
	return spot, total, res, den, reason, posts, args.Error(6)
}

// MockPresenter はレスポンス形状のモック
type MockDistillPresenter struct{}

func (p *MockDistillPresenter) Output(s *entities.Spot, ts value_objects.TotalScore, rc value_objects.ResonanceCount, ds value_objects.DensityScore, r value_objects.Reason, ps []*entities.Post) *usecase.DistillRecommendationResponse {
	return &usecase.DistillRecommendationResponse{
		Recommendation: &usecase.RecommendationResult{
			Spot: usecase.SpotOutput{ID: s.ID.Value(), Name: s.Name.String()},
			DistillationAnalysis: usecase.AnalysisOutput{
				TotalScore: ts.Float64(),
				ResonanceScore: rc.Int(),
				DensityScore: ds.Int(),
				Reason: r.String(),
			},
		},
	}
}

// --- TEST 本体 ---

func TestDistillRecommendation_Execute(t *testing.T) {
	// 共通データ準備
	malloy, _ := entities.NewUser(2, "local_malloy", "malloy@example.com", "hashed_password")
	bobSpot, _ := entities.NewSpot(1, "ボブの隠れ家", 35.6467, 139.7101, 1)
	
	// Value Objects の準備（エラーハンドリングは省略）
	ts, _ := value_objects.NewTotalScore(2.5)
	rc, _ := value_objects.NewResonanceCount(1)
	ds, _ := value_objects.NewDensityScore(1)
	reason, _ := value_objects.NewReason("共鳴による蒸留結果です")
	zeroTS, _ := value_objects.NewTotalScore(0)
	zeroRC, _ := value_objects.NewResonanceCount(0)
	zeroDS, _ := value_objects.NewDensityScore(0)
	noReason, _ := value_objects.NewReason("no recommendation")

	tests := []struct {
		name      string
		input     usecase.DistillRecommendationInput
		setupMock func(am *DistillMockAuthService, rs *MockRecommendationService)
		wantErr   bool
		check     func(t *testing.T, out *usecase.DistillRecommendationResponse)
	}{
		{
			name: "【正常系】共鳴するスポットが見つかり、正しく蒸留結果を返す",
			input: usecase.DistillRecommendationInput{
				Token: "valid_token", Latitude: 35.6467, Longitude: 139.7101,
			},
			setupMock: func(am *DistillMockAuthService, rs *MockRecommendationService) {
				am.On("VerifyToken", mock.Anything, "valid_token").Return(malloy, nil)
				// Domain Service の Distill が呼ばれることを期待
				rs.On("Distill", mock.Anything, malloy, mock.Anything, mock.Anything).
					Return(bobSpot, ts, rc, ds, reason, []*entities.Post{}, nil)
			},
			wantErr: false,
			check: func(t *testing.T, out *usecase.DistillRecommendationResponse) {
				assert.NotNil(t, out.Recommendation)
				assert.Equal(t, "ボブの隠れ家", out.Recommendation.Spot.Name)
				assert.Equal(t, 2.5, out.Recommendation.DistillationAnalysis.TotalScore)
			},
		},
		{
			name: "【正常系】周囲に推奨スポットがない場合は nil を返す",
			input: usecase.DistillRecommendationInput{Token: "valid_token", Latitude: 0, Longitude: 0},
			setupMock: func(am *DistillMockAuthService, rs *MockRecommendationService) {
				am.On("VerifyToken", mock.Anything, "valid_token").Return(malloy, nil)
				rs.On("Distill", mock.Anything, malloy, mock.Anything, mock.Anything).
					Return((*entities.Spot)(nil), zeroTS, zeroRC, zeroDS, noReason, []*entities.Post{}, nil)
			},
			wantErr: false,
			check: func(t *testing.T, out *usecase.DistillRecommendationResponse) {
				assert.Nil(t, out) // Interactor のロジック通り nil が返るか
			},
		},
		{
			name: "【異常系】トークンが不正で認証に失敗する",
			input: usecase.DistillRecommendationInput{Token: "bad_token"},
			setupMock: func(am *DistillMockAuthService, rs *MockRecommendationService) {
				am.On("VerifyToken", mock.Anything, "bad_token").Return((*entities.User)(nil), errors.New("unauthorized"))
			},
			wantErr: true,
		},
		{
			name: "【異常系】不正な座標（緯度）が渡された場合、バリデーションで弾く",
			input: usecase.DistillRecommendationInput{Token: "valid_token", Latitude: 100.0, Longitude: 0},
			setupMock: func(am *DistillMockAuthService, rs *MockRecommendationService) {
				am.On("VerifyToken", mock.Anything, "valid_token").Return(malloy, nil)
			},
			wantErr: true, // Latitude の NewLatitude でエラーになるはず
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am := new(DistillMockAuthService)
			rs := new(MockRecommendationService)
			tt.setupMock(am, rs)
			
			interactor := usecase.NewDistillRecommendationInteractor(&MockDistillPresenter{}, rs, am)

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
			rs.AssertExpectations(t)
		})
	}
}