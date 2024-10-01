package controllers

import (
	"net/http"

	"github.com/itsluthfi/hlf-mtcnv2/rest-api-go/blockchain"
	"github.com/itsluthfi/hlf-mtcnv2/rest-api-go/models"
	"github.com/itsluthfi/hlf-mtcnv2/rest-api-go/utils/token"

	"github.com/gin-gonic/gin"
)

func CurrentUser(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := models.GetUserByID(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
}

type RegisterInput struct {
	Email     string `json:"email" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := blockchain.Register(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	u.Email = input.Email
	u.Username = input.Username
	u.Password = input.Password
	u.Firstname = input.Firstname
	u.Lastname = input.Lastname
	u.IsSuper = false

	accountID, err := blockchain.AccountID(input.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u.WalletID = accountID

	_, err = u.SaveUser()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration success"})
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}

	u.Username = input.Username
	u.Password = input.Password

	// Generate Token
	token, err := models.LoginCheck(u.Username, u.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

type VerifyInput struct {
	Username string `json:"username" binding:"required"`
	Code     string `json:"code" binding:"required"`
}
