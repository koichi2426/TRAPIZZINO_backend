package presenter

import (
	"time"
	"app/domain/entities"
	"app/domain/value_objects"
	"app/usecase"
)

type distillRecommendationPresenter struct{}

func NewDistillRecommendationPresenter() usecase.DistillRecommendationPresenter {
	return &distillRecommendationPresenter{}
}

// Output はユースケースから渡されたバラバラのドメインオブジェクトを、
// API仕様書通りの JSON 構造（DTO）へ構造化します。
func (p *distillRecommendationPresenter) Output(
	spot *entities.Spot,
	totalScore value_objects.TotalScore,
	resonanceCount value_objects.ResonanceCount,
	density value_objects.DensityScore,
	reason value_objects.Reason,
	posts []*entities.Post,
) *usecase.DistillRecommendationResponse {
	if spot == nil {
		return nil
	}

	// 1. Spot情報の整形 (仕様書の spot ブロックに対応)
	// spot.Location 経由ではなく、Entity のフィールドから直接取得する形に修正
	spotOut := usecase.SpotOutput{
		ID:     spot.ID.Value(),
		Name:   spot.Name.String(),
		MeshID: spot.MeshID.String(),
		Location: usecase.Location{
			Latitude:  spot.Latitude.Value(),
			Longitude: spot.Longitude.Value(),
		},
	}

	// 2. 蒸留分析データの整形 (仕様書の distillation_analysis ブロックに対応)
	analysisOut := usecase.AnalysisOutput{
		ResonanceScore: resonanceCount.Int(),
		DensityScore:   density.Int(),
		TotalScore:     totalScore.Float64(),
		Reason:         reason.String(),
	}

	// 3. 投稿リストの整形 (仕様書の posts ブロックに対応)
	postsOut := make([]usecase.PostOutput, 0, len(posts))
	for _, post := range posts {
		postsOut = append(postsOut, usecase.PostOutput{
			ID:       post.ID.Value(),
			// UserName VO から string を取り出すように修正
			UserName: post.UserName.String(), 
			Caption:  post.Caption.String(),
			ImageURL: post.ImageURL.String(),
			PostedAt: post.PostedAt.Format(time.RFC3339),
		})
	}

	// 4. 仕様書通りの 3 ブロック構造でレスポンスを組み立て
	return &usecase.DistillRecommendationResponse{
		Recommendation: &usecase.RecommendationResult{
			Spot:                 spotOut,
			DistillationAnalysis: analysisOut,
			Posts:                postsOut,
		},
	}
}