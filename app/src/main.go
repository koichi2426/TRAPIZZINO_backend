package main

import (
	"database/sql"
	"log"
	"os"
	"fmt"
	"src/infrastructure/database/postgres" // MySQLからPostgreSQLに読み替え
	"src/infrastructure/router"

	_ "github.com/lib/pq" // PostgreSQLドライバー
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 1. 環境変数から設定を読み込み
	config := postgres.NewConfigFromEnv()
	
	// 2. データベース接続 (PostgreSQL)
	db, err := sql.Open("postgres", config.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 接続確認
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// 3. Echo インスタンスの生成
	e := echo.New()

	// ミドルウェアの設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	// 4. ルーターの初期化 (DIの実行場所)
	router.InitRoutes(e, db)

	// 5. サーバー起動 (ポートはガイドラインに従い 8000)
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
	}
	
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}