package controller

import (
	"net/http"
	"strconv"
	"strings"

	"app/src/usecase"
	"github.com/labstack/echo/v4"
)

type DistillRecommendationController struct {
	usecase usecase.DistillRecommendationUseCase
}

func NewDistillRecommendationController(u usecase.DistillRecommendationUseCase) *DistillRecommendationController {
	return &DistillRecommendationController{usecase: u}
}

func (ctrl *DistillRecommendationController) Execute(c echo.Context) error {
	// 1. Authorization ヘッダーから生のトークン文字列を取得
	// Bearer <token> の形式で送られてくることを想定します
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authorization header is required"})
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid authorization format"})
	}
	token := tokenParts[1]

	// 2. クエリパラメータ（現在地）のパース
	latStr := c.QueryParam("latitude")
	lngStr := c.QueryParam("longitude")

	if latStr == "" || lngStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Latitude and longitude are required"})
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid latitude format"})
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid longitude format"})
	}

	// 3. ユースケースの実行（UserIDではなくTokenを渡す）
	input := usecase.DistillRecommendationInput{
		Token:     token,
		Latitude:  lat,
		Longitude: lng,
	}

	output, err := ctrl.usecase.Execute(c.Request().Context(), input)
	if err != nil {
		// 認証エラーなどのドメインエラーを適切にハンドリング
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	// 4. 結果が空の場合のハンドリング
	if output == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "No recommendation found in your resonance circle"})
	}

	// 5. 成功レスポンス
	return c.JSON(http.StatusOK, output)
}