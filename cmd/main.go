package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/promotionos/analytics-service/internal/api"
	"github.com/promotionos/analytics-service/internal/application"
	"github.com/promotionos/analytics-service/internal/infrastructure/event"
	infrarepo "github.com/promotionos/analytics-service/internal/infrastructure/repository"
	infraservice "github.com/promotionos/analytics-service/internal/infrastructure/service"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	redisURL := os.Getenv("REDIS_URL")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	// DB connection with analytics schema
	dsn := dbURL + "?search_path=analytics"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// Run migrations
	sqlDB, _ := db.DB()
	runMigrations(sqlDB)

	// Redis connection
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}
	redisClient := redis.NewClient(opt)
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Wire dependencies
	publisher := event.NewPublisher(redisClient)
	repo := infrarepo.NewAnalyticsRepositoryImpl(db)
	liftCalc := infraservice.NewLiftCalculator()
	burnTracker := infraservice.NewBudgetBurnTracker()

	onExhausted := func(tenantID, campaignID string, metrics interface{}) error {
		m := metrics.(*application.MetricsView)
		return publisher.PublishBudgetExhausted(
			tenantID, campaignID,
			m.TotalAmount, m.BurnedAmount, m.BudgetBurnPercent, m.RedemptionCount,
		)
	}

	svc := application.NewAnalyticsService(repo, liftCalc, burnTracker, nil)

	// Start event consumers
	event.StartConsumers(redisClient, svc)

	// Routes
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "analytics-service"})
	})

	h := api.NewAnalyticsHandler(svc)
	r.GET("/analytics/campaigns/:id/report", h.GetReport)
	r.GET("/analytics/campaigns/:id/burn", h.GetBurn)
	r.GET("/analytics/campaigns/:id/lift", h.GetLift)

	log.Printf("Analytics Service starting on port %s", port)
	r.Run(fmt.Sprintf(":%s", port))
}

func runMigrations(db *sql.DB) {
	goose.SetBaseFS(nil)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Goose dialect error: %v", err)
	}
	if err := goose.Up(db, "db/migrations"); err != nil {
		log.Printf("Migration warning: %v", err)
	}
}
