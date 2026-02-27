package presenter

import (
	"app/domain/entities"
	"app/usecase"
)

// MeshSpotPresenterは、スポットと投稿の登録・更新結果をAPI仕様のJSONレスポンス形式に整形する役割を担います。
type meshSpotPresenter struct{}

func NewMeshSpotPresenter() usecase.RegisterSpotPostPresenter {
	return &meshSpotPresenter{}
}

// Outputはspotとpostのペアをレスポンスとして返します。
func (p *meshSpotPresenter) Output(spot *entities.Spot, post *entities.Post) *usecase.RegisterSpotPostOutput {
	return &usecase.RegisterSpotPostOutput{
		SpotID: spot.ID.Value(),
		PostID: post.ID.Value(),
	}
}
