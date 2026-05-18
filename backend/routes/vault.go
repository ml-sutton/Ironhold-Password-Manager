package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ml-sutton/Ironhold-Password-Manager/backend/models"
)

func CreateVaultEntry(ctx *gin.Context) {
	// Checks headers
	var requestHeaders models.CreateVaultEntryHeaders
	var requestHeaderBindError = ctx.ShouldBindHeader(&requestHeaders)
	if nil != requestHeaderBindError {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"": "",
		})
	}

	var requestBody models.CreateVaultEntryRequest
	var requestBodyBindError = ctx.ShouldBindBodyWithJSON(&requestBody)
	if nil != requestBodyBindError {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"": "",
		})
	}

}
