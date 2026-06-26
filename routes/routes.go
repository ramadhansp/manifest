package routes

import (
	"manifest-api/controller"
	"manifest-api/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine, ctrl *controller.AppController) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := r.Group("/api")

	api.POST("/login", ctrl.Login)
	api.POST("/register", ctrl.Register)

	protected := api.Group("/")
	protected.Use(middleware.JWTAuth())

	protected.POST("/shipping-agents", ctrl.CreateShippingAgent)
	protected.GET("/shipping-agents", ctrl.GetShippingAgents)
	protected.GET("/shipping-agents/:id", ctrl.GetShippingAgent)

	protected.POST("/vessels", ctrl.CreateVessel)
	protected.GET("/vessels", ctrl.GetVessels)
	protected.GET("/vessels/:id", ctrl.GetVessel)

	protected.POST("/manifests", ctrl.CreateManifest)
	protected.GET("/manifests", ctrl.GetManifests)
	protected.GET("/manifests/:id", ctrl.GetManifest)
	protected.POST("/manifests/:id/details", ctrl.AddManifestDetail)

	protected.POST("/bc11", ctrl.CreateBC11)
	protected.POST("/npe", ctrl.CreateNPE)

	protected.GET("/summary", ctrl.GetSummary)
}
