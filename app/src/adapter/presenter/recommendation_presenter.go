package presenter

import (
	"src/domain/entities"
	"src/usecase"
)

// RecommendationPresenterは、レコメンデーションユースケースの出力DTOをAPI仕様のJSONレスポンス形式に整形する役割を担います。
type recommendationPresenter struct{}

func NewRecommendationPresenter() usecase.DistillRecommendationPresenter {
	return &recommendationPresenter{}
}

// Outputは最適な1軒のspotとposts配列を構造化して返します。
func (p *recommendationPresenter) Output(spot *entities.Spot) *usecase.DistillRecommendationOutput {
	return &usecase.DistillRecommendationOutput{
		SpotID:   spot.ID.Value(),
		SpotName: spot.Name,
	}
}
