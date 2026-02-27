package controller

import (
	"net/http"
	"strconv"
	"app/usecase"
	"github.com/labstack/echo/v4"
)

// RecommendationControllerは、GET /v1/recommendation/distill のリクエストを受け取り、
// クエリパラメータ lat, lng をfloat64に変換しユースケース層のDTOに渡す役割を担います。
type RecommendationController struct {
	usecase usecase.DistillRecommendationUseCase
}

func NewRecommendationController(u usecase.DistillRecommendationUseCase) *RecommendationController {
	return &RecommendationController{usecase: u}
}

// Executeはクエリパラメータを解析し、ユースケースを呼び出してレスポンスを返します。
func (ctrl *RecommendationController) Execute(c echo.Context) error {
	userIDStr := c.QueryParam("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user_id"})
	}
	input := usecase.DistillRecommendationInput{
		UserID: userID,
	}
	output, err := ctrl.usecase.Execute(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, output)
}
