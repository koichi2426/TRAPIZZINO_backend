package router

import (
	"database/sql"
	"src/adapter/controller"
	"src/adapter/presenter"
	"src/infrastructure/database/postgres"
	impl_services "src/infrastructure/domain_impl/services"
	"src/usecase"

	"github.com/labstack/echo/v4"
)

// InitRoutesは、すべての依存関係を手動で注入し、Echoのルーティングを構成します。
func InitRoutes(e *echo.Echo, db *sql.DB) {
	// 1. インフラ層（Repository / Domain Impl）の初期化
	spotRepo := postgres.NewSpotRepository(db)
	postRepo := postgres.NewPostRepository(db)
	userRepo := postgres.NewUserRepository(db)

	authService := impl_services.NewAuthDomainServiceImpl()
	meshService := impl_services.NewMeshCalculationServiceImpl()
	recommendationService := impl_services.NewRecommendationServiceImpl()

	// 2. プレゼンターの初期化
	authPresenter := presenter.NewAuthPresenter()
	userSpotPresenter := presenter.NewUserSpotPresenter()
	meshSpotPresenter := presenter.NewMeshSpotPresenter()
	recommendationPresenter := presenter.NewRecommendationPresenter()

	// 3. ユースケースの初期化
	authLoginUsecase := usecase.NewAuthLoginInteractor(authPresenter, authService)
	userSignupUsecase := usecase.NewUserSignupInteractor(authPresenter, userRepo, authService)
	registerSpotUsecase := usecase.NewRegisterSpotPostInteractor(meshSpotPresenter, spotRepo, postRepo)
	listMySpotsUsecase := usecase.NewListMySpotsInteractor(userSpotPresenter, spotRepo, postRepo)
	distillRecommendationUsecase := usecase.NewDistillRecommendationInteractor(recommendationPresenter, recommendationService, spotRepo)

	// 4. コントローラーの初期化
	authController := controller.NewAuthController(authLoginUsecase)
	userController := controller.NewUserController(userSignupUsecase)
	meshSpotController := controller.NewMeshSpotController(registerSpotUsecase)
	userSpotController := controller.NewUserSpotController(listMySpotsUsecase)
	recommendationController := controller.NewRecommendationController(distillRecommendationUsecase)

	// 5. ルーティング定義
	v1 := e.Group("/v1")

	v1.POST("/auth/login", authController.Execute)
	v1.POST("/users/signup", userController.Execute)
	v1.PUT("/mesh/spots", meshSpotController.Execute)
	v1.GET("/users/me/spots", userSpotController.Execute)
	v1.GET("/recommendation/distill", recommendationController.Execute)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})
}
