package controller

import (
	"net/http"
	"strings"
	"src/usecase"
	"github.com/labstack/echo/v4"
)

// UserSpotControllerは、GET /v1/users/me/spots のリクエストを受け取り、
// AuthorizationヘッダーからBearerトークンを抽出し、ユースケース層のDTOに変換する役割を担います。
type UserSpotController struct {
	usecase usecase.ListMySpotsUseCase
}

func NewUserSpotController(u usecase.ListMySpotsUseCase) *UserSpotController {
	return &UserSpotController{usecase: u}
}

// ExecuteはAuthorizationヘッダーを解析し、ユースケースを呼び出してレスポンスを返します。
func (ctrl *UserSpotController) Execute(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	input := usecase.ListMySpotsInput{
		UserID: 0, // トークンからユーザーIDを抽出する処理は後続で追加
	}
	output, err := ctrl.usecase.Execute(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, output)
}
