package presenter

import (
	"src/domain/entities"
	"src/usecase"
)

// AuthPresenterは、認証ユースケースの出力DTOをAPI仕様のJSONレスポンス形式に整形する役割を担います。
type authPresenter struct{}

func NewAuthPresenter() usecase.AuthLoginPresenter {
	return &authPresenter{}
}

// Outputはaccess_tokenとuserオブジェクトを統合したレスポンスを返します。
func (p *authPresenter) Output(token string) *usecase.AuthLoginOutput {
	return &usecase.AuthLoginOutput{
		Token: token,
	}
}
