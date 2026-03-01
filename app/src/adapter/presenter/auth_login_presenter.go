package presenter

import (
	"app/src/usecase"
)

// AuthLoginPresenterは、認証ユースケースの出力DTOをAPI仕様のJSONレスポンス形式に整形する役割を担います。
type authLoginPresenter struct{}

func NewAuthLoginPresenter() usecase.AuthLoginPresenter {
	return &authLoginPresenter{}
}

// Outputはaccess_tokenとuserオブジェクトを統合したレスポンスを返します。
func (p *authLoginPresenter) Output(token string) *usecase.AuthLoginOutput {
	return &usecase.AuthLoginOutput{
		Token: token,
	}
}
