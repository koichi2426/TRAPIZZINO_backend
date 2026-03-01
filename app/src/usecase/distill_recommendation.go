package usecase

import (
	"context"
	"app/src/domain/entities"
	"app/src/domain/services"
	"app/src/domain/value_objects"
)

// DistillRecommendationInput はコントローラーから渡される入力データです
type DistillRecommendationInput struct {
	Token     string
	Latitude  float64
	Longitude float64
}

// DistillRecommendationResponse はフロントエンドへ返す最終的なレスポンス形状です
type DistillRecommendationResponse struct {
	Recommendation *RecommendationResult `json:"recommendation"`
}

type RecommendationResult struct {
	Spot                 SpotOutput      `json:"spot"`
	DistillationAnalysis AnalysisOutput  `json:"distillation_analysis"`
	Posts                []PostOutput     `json:"posts"`
}

// ... (SpotOutput, Location, AnalysisOutput, PostOutput の定義は同一のため維持) ...

type SpotOutput struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	MeshID   string   `json:"mesh_id"`
	Location Location `json:"location"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type AnalysisOutput struct {
	ResonanceScore int     `json:"resonance_score"`
	DensityScore   int     `json:"density_score"`
	TotalScore     float64 `json:"total_score"`
	Reason         string  `json:"reason"`
}

type PostOutput struct {
	ID       int    `json:"id"`
	UserName string `json:"user_name"`
	Caption  string `json:"caption"`
	ImageURL string `json:"image_url"`
	PostedAt string `json:"posted_at"`
}

// DistillRecommendationPresenter の引数をバラバラに変更
// これにより value_objects パッケージへの依存を最小限にします
type DistillRecommendationPresenter interface {
	Output(
		spot *entities.Spot,
		totalScore value_objects.TotalScore,
		resonanceCount value_objects.ResonanceCount,
		density value_objects.DensityScore,
		reason value_objects.Reason,
		posts []*entities.Post,
	) *DistillRecommendationResponse
}

type DistillRecommendationUseCase interface {
	Execute(ctx context.Context, input DistillRecommendationInput) (*DistillRecommendationResponse, error)
}

type distillRecommendationInteractor struct {
	presenter      DistillRecommendationPresenter
	recommendation services.RecommendationService
	authService    services.AuthDomainService
}

func NewDistillRecommendationInteractor(
	p DistillRecommendationPresenter,
	r services.RecommendationService,
	a services.AuthDomainService,
) DistillRecommendationUseCase {
	return &distillRecommendationInteractor{
		presenter:      p,
		recommendation: r,
		authService:    a,
	}
}

func (i *distillRecommendationInteractor) Execute(ctx context.Context, input DistillRecommendationInput) (*DistillRecommendationResponse, error) {
	// 1. ユーザーの特定
	user, err := i.authService.VerifyToken(ctx, input.Token)
	if err != nil {
		return nil, err
	}

	// 2. 現在地の Value Object 化
	lat, err := value_objects.NewLatitude(input.Latitude)
	if err != nil {
		return nil, err
	}
	lng, err := value_objects.NewLongitude(input.Longitude)
	if err != nil {
		return nil, err
	}

	// 3. 蒸留アルゴリズム（Domain Service）の実行
	// 戻り値をバラバラで受け取ることにより、ドメイン層内での循環参照を回避します
	spot, totalScore, resonanceCount, density, reason, posts, err := i.recommendation.Distill(ctx, user, lat, lng)
	if err != nil {
		return nil, err
	}

	// 4. 計算結果の空チェック
	if spot == nil {
		return nil, nil
	}

	// 5. プレゼンターへ各ドメインオブジェクトを渡し、出力用 DTO を生成します
	return i.presenter.Output(spot, totalScore, resonanceCount, density, reason, posts), nil
}