package controller

import (
	"net/http"
	"strings"
	"app/usecase"
	"github.com/labstack/echo/v4"
)

// UserSpotControllerは、GET /v1/users/me/spots のリクエストを受け取り、
// AuthorizationヘッダーからBearerトークンを抽出し、ユースケース層のDTOに変換します。
type UserSpotController struct {
	usecase usecase.ListMySpotsUseCase
}

func NewUserSpotController(u usecase.ListMySpotsUseCase) *UserSpotController {
	return &UserSpotController{usecase: u}
}

// ExecuteはAuthorizationヘッダーを解析し、ユースケースを呼び出してレスポンスを返します。
func (ctrl *UserSpotController) Execute(c echo.Context) error {
	// 1. Authorizationヘッダーからトークンを抽出
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "authorization header is required"})
	}

	// "Bearer <token>" の形式からトークン部分のみを取得
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid authorization format"})
	}
	token := parts[1]

	// 2. ユースケース層のDTOにトークンをセット
	// ListMySpotsInput側に Token string フィールドがある前提です
	input := usecase.ListMySpotsInput{
		Token: token,
	}

	// 3. ユースケースの呼び出し
	output, err := ctrl.usecase.Execute(c.Request().Context(), input)
	if err != nil {
		// 認証失敗なら401、その他なら500などエラー内容に応じて調整
		// ここでは簡略化のためInternalServerErrorにしていますが、
		// 実際はエラー型を判定してステータスコードを分けるのが理想的です
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// 4. 成功レスポンスを返却
	return c.JSON(http.StatusOK, output)
}