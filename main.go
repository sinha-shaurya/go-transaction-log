package main

import (
	"fmt"
	"go-transaction-log/database"
	"go-transaction-log/dto"
	"go-transaction-log/transaction/routes"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("v1/health", healthCheck)

	// setup router group
	apiGroup := router.Group("/api")
	routes.TransactionRoutes(apiGroup)

	
	router.Run("localhost:8080")
}

func healthCheck(ctx *gin.Context) {
	healthCheckResponse := dto.HealthCheckResponse{
		Status: "Healthy",
	}
	db, err := database.GetDbContext()
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
	}
	fmt.Println(db)


	// err = db.AutoMigrate(&models.Transaction{})
    // if err != nil {
    //     panic("failed to migrate database")
    // }


	ctx.IndentedJSON(http.StatusOK, healthCheckResponse)
}
