package controller

import (
	"net/http"
	"app/usecase"
	"github.com/labstack/echo/v4"
)

// AuthControllerは、POST /v1/auth/login のリクエストを受け取り、
// email, password をユースケース層のDTOに変換して認証処理を行う役割を担います。
type AuthController struct {
	usecase usecase.AuthLoginUseCase
}

func NewAuthController(u usecase.AuthLoginUseCase) *AuthController {
	return &AuthController{usecase: u}
}

// Executeはリクエストボディをバインドし、ユースケースを呼び出してレスポンスを返します。
func (ctrl *AuthController) Execute(c echo.Context) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	input := usecase.AuthLoginInput{
		Username: req.Username,
		Password: req.Password,
	}
	output, err := ctrl.usecase.Execute(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, output)
}
