package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	businesslogic "github.com/ml-sutton/Ironhold-Password-Manager/backend/business-logic"
	dataaccesslayer "github.com/ml-sutton/Ironhold-Password-Manager/backend/data-access-layer"
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
	// checks authentication
	var authenticationError = businesslogic.AuthenticateUser(&requestHeaders)
	if nil != authenticationError {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"": "",
		})
	}
	// checks request validity
	var requestBody models.CreateVaultEntryRequest
	var requestBodyBindError = ctx.ShouldBindBodyWithJSON(&requestBody)
	if nil != requestBodyBindError {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"": "",
		})
	}
	// Checks for database clash
	var attemptedVaultEntry *models.VaultEntry = models.CreateVaultEntryFromRequest("", requestBody)
	var blobExistsError error = dataaccesslayer.CheckExists(attemptedVaultEntry)
	if nil != blobExistsError {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"": "",
		})
	}
	// writes to database
	var databaseWriteError error = dataaccesslayer.WriteVaultEntry(attemptedVaultEntry)
	if nil != databaseWriteError {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"": "",
		})
	} else {
		ctx.JSON(http.StatusCreated, gin.H{
			"": "",
		})
	}
}
