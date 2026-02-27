package controller

import (
	"net/http"
	"app/usecase"
	"github.com/labstack/echo/v4"
)

// MeshSpotControllerは、PUT /v1/mesh/spots のリクエストを受け取り、
// リクエストボディのバインドとヘッダーからのBearerトークン抽出を行い、ユースケース層のDTOに変換する役割を担います。
type MeshSpotController struct {
	usecase usecase.RegisterSpotPostUseCase
}

func NewMeshSpotController(u usecase.RegisterSpotPostUseCase) *MeshSpotController {
	return &MeshSpotController{usecase: u}
}

// ExecuteはリクエストボディとAuthorizationヘッダーを解析し、ユースケースを呼び出してレスポンスを返します。
func (ctrl *MeshSpotController) Execute(c echo.Context) error {
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
	input := usecase.RegisterSpotPostInput{
		UserID:    0, // トークンからユーザーIDを抽出する処理は後続で追加
		SpotName:  req.SpotName,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		ImageURL:  req.ImageURL,
		Caption:   req.Caption,
		Overwrite: req.Overwrite,
	}
	output, err := ctrl.usecase.Execute(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, output)
}
