package main

import (
	"context"
	"hexagonal/ccgo01/internal/adapter/handler"
	"hexagonal/ccgo01/internal/adapter/handler/middleware"
	"hexagonal/ccgo01/internal/adapter/repository"
	repo "hexagonal/ccgo01/internal/adapter/repository"
	"hexagonal/ccgo01/internal/config"
	"hexagonal/ccgo01/internal/core/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	config.LoadEnv(".env")
	dsn := config.GetEnv("DATABASE_URL", "host=localhost user=postgres password=postgres dbname=walletdb port=5432 sslmode=disable TimeZone=Asia/Bangkok")
	port := config.GetEnv("PORT", "8080")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&repository.WalletGORM{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}
	sqlDB, err := db.DB()

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	walletRepo := repo.NewPostgresWalletRepository(db)
	walletSvc := service.NewWalletService(walletRepo)
	walletHdr := handler.NewWalletHandler(walletSvc)

	r := gin.Default()
	r.Use(middleware.CORS())

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})
	wallets := r.Group("/wallets")
	{
		wallets.POST("", walletHdr.CreateWallet)
		wallets.GET("/:id", walletHdr.GetWallet)
		wallets.POST("/:id/deposit", walletHdr.Deposit)
		wallets.POST("/:id/withdraw", walletHdr.Withdraw)
		wallets.DELETE("/:id", walletHdr.DeleteWallet)
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Printf("received signal %v - shutting down gracefully...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}

	if err == nil {
		sqlDB.Close()
	}

	log.Println("server stopped clenly")
}
