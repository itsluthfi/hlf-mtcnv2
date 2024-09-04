package controllers

import (
	"net/http"

	"github.com/itsluthfi/hlf-mtcnv2/rest-api-go/blockchain"
	"github.com/itsluthfi/hlf-mtcnv2/rest-api-go/models"
	"github.com/itsluthfi/hlf-mtcnv2/rest-api-go/utils/token"

	"github.com/gin-gonic/gin"
)

func Migrate(c *gin.Context) {
	models.Migrate()
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func GetTransactions(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, err := models.GetUsernameByID(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	accountID, _ := blockchain.AccountID(username)
	Transactions := models.GetTransactions(username, accountID)

	c.JSON(http.StatusOK, gin.H{"data": Transactions})
}
