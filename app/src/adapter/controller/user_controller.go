package controller

import (
	"net/http"
	"app/usecase"
	"github.com/labstack/echo/v4"
)

// UserControllerは、POST /v1/users/signup のリクエストを受け取り、
// username, email, password をユースケース層のDTOに変換して新規ユーザー登録処理を行う役割を担います。
type UserController struct {
	usecase usecase.UserSignupUseCase
}

func NewUserController(u usecase.UserSignupUseCase) *UserController {
	return &UserController{usecase: u}
}

// Executeはリクエストボディをバインドし、ユースケースを呼び出してレスポンスを返します。
func (ctrl *UserController) Execute(c echo.Context) error {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	input := usecase.UserSignupInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	output, err := ctrl.usecase.Execute(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, output)
}
