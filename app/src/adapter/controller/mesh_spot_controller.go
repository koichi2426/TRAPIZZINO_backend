package controller

import (
	"net/http"
	"strings"
	"app/usecase"
	"app/domain/services"
	"github.com/labstack/echo/v4"
)

type MeshSpotController struct {
	usecase     usecase.RegisterSpotPostUseCase
	authService services.AuthDomainService
}

func NewMeshSpotController(u usecase.RegisterSpotPostUseCase, a services.AuthDomainService) *MeshSpotController {
	return &MeshSpotController{
		usecase:     u,
		authService: a,
	}
}

func (ctrl *MeshSpotController) Execute(c echo.Context) error {
	// 1. Authorization ヘッダーから Bearer トークンを取得
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing or invalid authorization header"})
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// 2. トークンを検証し、ユーザー情報を復元
	user, err := ctrl.authService.VerifyToken(c.Request().Context(), tokenString)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized: invalid token"})
	}

	// 3. リクエストボディのバインド
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

	// 4. トークンから取得した UserID と Username をセット
	// 修正ポイント: input.Username に user.Username.String() を代入
	input := usecase.RegisterSpotPostInput{
		UserID:    user.ID.Value(),
		Username:  user.Username.String(), 
		SpotName:  req.SpotName,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		ImageURL:  req.ImageURL,
		Caption:   req.Caption,
		Overwrite: req.Overwrite,
	}

	// 5. ユースケース実行
	output, err := ctrl.usecase.Execute(c.Request().Context(), input)
	if err != nil {
		// すでにデータが存在する場合などは StatusConflict(409) を返す
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, output)
}