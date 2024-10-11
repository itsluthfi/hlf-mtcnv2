package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/itsluthfi/hlf-mtcnv2/rest-api-go/blockchain"
	"github.com/itsluthfi/hlf-mtcnv2/rest-api-go/models"
	"github.com/itsluthfi/hlf-mtcnv2/rest-api-go/utils/token"

	"github.com/gin-gonic/gin"
)

type TransferReq struct {
	Value    string `json:"value"`
	Receiver string `json:"receiver"`
}

func Transfer(c *gin.Context) {
	var input TransferReq
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, err := models.GetUsernameByID(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	convertedValue, err := strconv.Atoi(input.Value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	receiverAccountID, err := models.GetAccountIDByUsername(input.Receiver)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	value, err := blockchain.Balance(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	convertedCurrentValue, err := strconv.Atoi(value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if convertedValue > convertedCurrentValue {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not enough coins"})
		return
	}

	err = blockchain.Transfer(username, receiverAccountID, convertedValue)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentTime := time.Now()
	var transactionLog models.Transactions
	transactionLog.Sender = username
	transactionLog.Receiver = input.Receiver
	transactionLog.Value = input.Value
	transactionLog.Date = fmt.Sprint(currentTime.Format("2006-01-02 15:04:05"))
	transactionLog.SaveTransaction()

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func Balance(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, err := models.GetUsernameByID(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	value, err := blockchain.Balance(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"value": value})
}
