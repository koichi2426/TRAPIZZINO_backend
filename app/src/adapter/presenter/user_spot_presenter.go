package presenter

import (
	"app/usecase"
)

// UserSpotPresenterは、ユーザーが持つ複数のスポットをuser_spots配列としてラップし、API仕様のJSONレスポンス形式に整形する役割を担います。
type userSpotPresenter struct{}

func NewUserSpotPresenter() usecase.ListMySpotsPresenter {
	return &userSpotPresenter{}
}

// OutputはSpotPostPairの配列をuser_spotsとしてラップして返します。
func (p *userSpotPresenter) Output(pairs []usecase.SpotPostPair) *usecase.ListMySpotsOutput {
	return &usecase.ListMySpotsOutput{
		Pairs: pairs,
	}
}
