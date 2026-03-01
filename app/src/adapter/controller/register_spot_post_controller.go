package controller

import (
	"net/http"
	"strings"
	"app/src/usecase"
	"github.com/labstack/echo/v4"
)

type RegisterSpotPostController struct {
	usecase usecase.RegisterSpotPostUseCase
}

func NewRegisterSpotPostController(u usecase.RegisterSpotPostUseCase) *RegisterSpotPostController {
	return &RegisterSpotPostController{
		usecase: u,
	}
}

func (ctrl *RegisterSpotPostController) Execute(c echo.Context) error {
	// 1. Authorization ヘッダーから Bearer トークンを取得
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing or invalid authorization header"})
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// 2. リクエストボディのバインド
	var req struct {
		SpotName  string  `json:"spot_name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		ImageURL  string  `json:"image_url"`
		Caption   string  `json:"caption"`
		Overwrite bool    `json:"overwrite"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// 3. ユースケース入力を組み立て
	input := usecase.RegisterSpotPostInput{
		Token:     tokenString,
		SpotName:  req.SpotName,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		ImageURL:  req.ImageURL,
		Caption:   req.Caption,
		Overwrite: req.Overwrite,
	}

	// 4. ユースケース実行
	output, err := ctrl.usecase.Execute(c.Request().Context(), input)
	if err != nil {
		// すでにデータが存在する場合などは StatusConflict(409) を返す
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, output)
}