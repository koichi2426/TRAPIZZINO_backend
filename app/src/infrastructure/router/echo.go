package router

import (
	"database/sql"
	"os"

	"app/adapter/controller"
	"app/adapter/presenter"
	"app/infrastructure/database/postgres"
	impl_services "app/infrastructure/domain_impl/services"
	"app/usecase"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo, db *sql.DB) {
	// 1. インフラ層（Repository / Domain Impl）の初期化
	spotRepo := postgres.NewSpotRepository(db)
	postRepo := postgres.NewPostRepository(db)
	userRepo := postgres.NewUserRepository(db)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "develop_secret_key_change_me"
	}

	authService := impl_services.NewAuthDomainServiceImpl(jwtSecret)
	recommendationService := impl_services.NewRecommendationServiceImpl(spotRepo)

	// 2. プレゼンターの初期化
	authPresenter := presenter.NewAuthPresenter()
	userSignupPresenter := presenter.NewUserSignupPresenter()
	userSpotPresenter := presenter.NewUserSpotPresenter()
	meshSpotPresenter := presenter.NewMeshSpotPresenter()
	recommendationPresenter := presenter.NewRecommendationPresenter()

	// 3. ユースケースの初期化
	authLoginUsecase := usecase.NewAuthLoginInteractor(authPresenter, userRepo, authService)
	userSignupUsecase := usecase.NewUserSignupInteractor(userSignupPresenter, userRepo, authService)
	
	// 【修正】第4引数に authService を追加
	registerSpotUsecase := usecase.NewRegisterSpotPostInteractor(meshSpotPresenter, spotRepo, postRepo, authService)
	
	// 【修正】authService を追加し、トークンからユーザーを特定できるようにします
	listMySpotsUsecase := usecase.NewListMySpotsInteractor(userSpotPresenter, spotRepo, postRepo, authService)
	
	distillRecommendationUsecase := usecase.NewDistillRecommendationInteractor(
		recommendationPresenter, 
		recommendationService, 
		authService,
	)

	// 4. コントローラーの初期化
	authController := controller.NewAuthController(authLoginUsecase)
	userController := controller.NewUserController(userSignupUsecase)
	
	// 【修正】コントローラーはユースケースにトークンを渡すだけにするので、authService の注入は不要に
	meshSpotController := controller.NewMeshSpotController(registerSpotUsecase)
	
	userSpotController := controller.NewUserSpotController(listMySpotsUsecase)
	recommendationController := controller.NewRecommendationController(distillRecommendationUsecase)

	// 5. ルーティング定義
	v1 := e.Group("/v1")

	v1.POST("/auth/login", authController.Execute)
	v1.POST("/users/signup", userController.Execute)

	// PUT メソッドで定義された「情報の蒸留」エンドポイント
	v1.PUT("/mesh/spots", meshSpotController.Execute)
	v1.GET("/users/me/spots", userSpotController.Execute)
	v1.GET("/recommendation/distill", recommendationController.Execute)

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})
}