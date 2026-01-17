package main

// @title Selfcart API
// @version 1.0
// @description Backend API for Selfcart application
// @host localhost:8080
// @BasePath /api

import (
	"log"
	"os"
	"selfcart/internal/config"
	"selfcart/internal/database"
	"selfcart/internal/handler"
	"selfcart/internal/repository"
	"selfcart/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "selfcart/docs" // ⚠️ WAJIB (sesuai module name)

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           SelfCart API
// @version         1.0
// @description     Backend API for SelfCart
// @termsOfService  https://selfcart.com/terms

// @contact.name   API Support
// @contact.email  support@selfcart.com

// @host 
// @BasePath /api

func main() {
	_ = godotenv.Load()

	cfg := config.Load()
	db := database.NewPostgres(cfg.DBDsn)

	// repositories
	storeRepo := repository.NewStoreRepo(db)
	categoryRepo := repository.NewCategoryRepo(db)
	productRepo := repository.NewProductRepo(db)
	cartRepo := repository.NewCartRepo(db)
	transactionRepo := repository.NewTransactionRepo(db)

	// services
	categoryService := service.NewCategoryService(categoryRepo)
	transactionService := service.NewTransactionService(transactionRepo)
	cartService := service.NewCartService(cartRepo, productRepo, transactionService)
	productService := service.NewProductService(productRepo)
	storeService := service.NewStoreService(storeRepo)

	// handlers
	storeHandler := handler.NewStoreHandler(storeService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	cartHandler := handler.NewCartHandler(cartService)
	productHandler := handler.NewProductHandler(productService)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	api := r.Group("/api")
	{
		api.GET("/store", storeHandler.GetStore)
		api.GET("/store/rack/:store_id", storeHandler.GetRack)
		api.GET("/store/rack/zone/:store_id/:zone", storeHandler.GetRackZone)
		api.GET("/categories", categoryHandler.GetAll)
		api.GET("/products", productHandler.GetAll)
		api.GET("/products/barcode/:barcode", productHandler.GetByBarcode)
		api.GET("/products/category/:category_id", productHandler.GetByCategory)
		api.POST("/cart/create", cartHandler.CreateCart)
		api.POST("/cart/items", cartHandler.AddItem)
		api.POST("/cart/remove", cartHandler.RemoveItem)
		api.GET("/cart/:cart_id", cartHandler.GetCart)
		api.POST("/cart/checkout", cartHandler.CheckOut)
	}
	if os.Getenv("ENABLE_SWAGGER") == "true" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})

	log.Println("Server starting on port", cfg.AppPort)

	r.Run(":" + cfg.AppPort)
}
