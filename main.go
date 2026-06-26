// @title           Manifest API
// @version         1.0
// @description     REST API untuk manajemen manifest kapal, BC11, dan NPE
// @host            localhost:8080
// @BasePath        /api
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     Masukkan token dengan format: Bearer {token}

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
	_ "manifest-api/docs"
	"manifest-api/middleware"
	"manifest-api/models"
	"manifest-api/repository"
	"manifest-api/routes"
	"manifest-api/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using OS env vars")
	}

	db := config.SetupDatabase()

	var count int64
	db.Model(&models.User{}).Count(&count)
	if count == 0 {
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		db.Create(&models.User{Username: "admin", Password: string(hash), Role: "Administrator"})
	}

	repoManager := repository.NewDBManager(db)
	appService := service.NewAppService(repoManager)
	appController := controller.NewAppController(appService)

	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	routes.RegisterRoutes(r, appController)

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
