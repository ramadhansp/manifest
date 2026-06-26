package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"manifest-api/config"
	"manifest-api/controller"
	"manifest-api/middleware"
	"manifest-api/repository"
	"manifest-api/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using OS env vars")
	}

	db := config.SetupDatabase()

	repoManager := repository.NewDBManager(db)
	appService := service.NewAppService(repoManager)
	appController := controller.NewAppController(appService)

	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	api := r.Group("/api")
	{
		api.POST("/shipping-agents", appController.CreateShippingAgent)
		api.GET("/shipping-agents", appController.GetShippingAgents)
		api.GET("/shipping-agents/:id", appController.GetShippingAgent)

		api.POST("/vessels", appController.CreateVessel)
		api.GET("/vessels", appController.GetVessels)
		api.GET("/vessels/:id", appController.GetVessel)

		api.POST("/manifests", appController.CreateManifest)
		api.GET("/manifests", appController.GetManifests)
		api.GET("/manifests/:id", appController.GetManifest)
		api.POST("/manifests/:id/details", appController.AddManifestDetail)

		api.POST("/bc11", appController.CreateBC11)
		api.POST("/npe", appController.CreateNPE)

		api.GET("/summary", appController.GetSummary)
		api.POST("/seed", appController.SeedData)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Printf("Server listening on port %s", port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
