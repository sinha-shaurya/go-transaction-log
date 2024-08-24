package controller

import (
	"go-transaction-log/database"
	"go-transaction-log/transaction/dto"
	"go-transaction-log/transaction/service"
	"go-transaction-log/transaction/transformer"

	"go-transaction-log/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, nil)
}

func AddTransaction(ctx *gin.Context) {
	var request dto.AddTransactionRequest

	// bind and validate request
	validationErr := utils.BindAndValidate(ctx, &request)
	if validationErr != nil {
		ctx.IndentedJSON(http.StatusBadRequest, validationErr.Error())
	}

	db, err := database.GetDbContext()
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	createErr := service.AddTransaction(ctx, db, dto.AddTransactionRequestDto(request))
	if createErr != nil {
		ctx.IndentedJSON(createErr.ErrorCode, createErr.GetErrorResponse())
		return
	}
	ctx.IndentedJSON(http.StatusAccepted, dto.AddTransactionResponse{Status: "ok"})
}

func GetTransactionById(ctx *gin.Context) {
	var request dto.GetTransactionByIdRequest
	validationErr := utils.BindAndValidate(ctx, &request)
	if validationErr != nil {
		ctx.IndentedJSON(http.StatusBadRequest, validationErr.Error())
	}

	db, err := database.GetDbContext()
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	transaction, apiErr := service.GetTransactionById(db, request.TranscationId)
	if apiErr != nil {
		ctx.IndentedJSON(apiErr.ErrorCode, apiErr.GetErrorResponse())
		return
	}

	response := transformer.BuildGetTransactionResponse(*transaction)
	ctx.IndentedJSON(http.StatusOK, response)
}

func GetTransactionByType(ctx *gin.Context) {
	var request dto.GetTransactionByTypeRequest
	validationErr := utils.BindAndValidate(ctx, &request)
	if validationErr != nil {
		ctx.IndentedJSON(http.StatusBadRequest, validationErr.Error())
		return
	}

	db, err := database.GetDbContext()
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	transactions, apiErr := service.GetTransactionByType(db, request.TranscationType)
	if apiErr != nil {
		ctx.IndentedJSON(apiErr.ErrorCode, apiErr.GetErrorResponse())
		return
	}

	response := transformer.BuildGetTransactionByTypeResponse(transactions)
	ctx.IndentedJSON(http.StatusOK, response)
}

func SumTransactions(ctx *gin.Context) {
	var request dto.GetTransactionByIdRequest
	validationErr := utils.BindAndValidate(ctx, &request)
	if validationErr != nil {
		ctx.IndentedJSON(http.StatusBadRequest, validationErr.Error())
		return
	}

	db, err := database.GetDbContext()
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	sum, apiErr := service.SumTransactions(db, request.TranscationId)
	if apiErr != nil {
		ctx.IndentedJSON(apiErr.ErrorCode, apiErr.GetErrorResponse())
		return
	}

	ctx.IndentedJSON(http.StatusOK, dto.SumTransactionResponse{Sum: sum})
}
