package presenter

import (
	"app/src/usecase"
)

// ListMySpotsPresenterは、ユーザーが持つ複数のスポットをuser_spots配列としてラップし、API仕様のJSONレスポンス形式に整形する役割を担います。
type listMySpotsPresenter struct{}

func NewListMySpotsPresenter() usecase.ListMySpotsPresenter {
	return &listMySpotsPresenter{}
}

// OutputはSpotPostPairの配列をuser_spotsとしてラップして返します。
func (p *listMySpotsPresenter) Output(pairs []usecase.SpotPostPair) *usecase.ListMySpotsOutput {
	return &usecase.ListMySpotsOutput{
		Pairs: pairs,
	}
}
