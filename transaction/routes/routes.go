package routes

import (
	"go-transaction-log/transaction/controller"

	"github.com/gin-gonic/gin"
)

func TransactionRoutes(router *gin.RouterGroup) {
	router = router.Group("transcationservice")
	router.GET("v1/health", controller.HealthCheck)

	// Define routes
	router.PUT("/transaction/:transactionId", controller.AddTransaction)
	router.GET("/transaction/:transactionId", controller.GetTransactionById)
	router.GET("/types/:transactionType", controller.GetTransactionByType)
	router.GET("/sum/:transactionId", controller.SumTransactions)
}
