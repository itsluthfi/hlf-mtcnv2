package main

import (
	"github.com/itsluthfi/hlf-mtcnv2/rest-api-go/blockchain"
	"github.com/itsluthfi/hlf-mtcnv2/rest-api-go/controllers"
	"github.com/itsluthfi/hlf-mtcnv2/rest-api-go/middlewares"
	"github.com/itsluthfi/hlf-mtcnv2/rest-api-go/models"

	"github.com/gin-gonic/gin"
)

func main() {
	models.ConnectDatabase()

	r := gin.Default()

	public := r.Group("/api", CORSMiddleware())

	public.OPTIONS("/migrate", CORSMiddleware())
	public.GET("/migrate", controllers.Migrate)

	public.OPTIONS("/register", CORSMiddleware())
	public.POST("/register", controllers.Register)

	public.OPTIONS("/login", CORSMiddleware())
	public.POST("/login", controllers.Login)

	protectedUsers := r.Group("/api/user", CORSMiddleware())
	protectedUsers.Use(middlewares.JWTAuthMiddleware())

	protectedUsers.OPTIONS("/detail", CORSMiddleware())
	protectedUsers.GET("/detail", controllers.CurrentUser)

	protectedUsers.OPTIONS("/transfer", CORSMiddleware())
	protectedUsers.POST("/transfer", controllers.Transfer)

	protectedUsers.OPTIONS("/mint", CORSMiddleware())
	protectedUsers.POST("/mint", blockchain.Mint)

	protectedUsers.OPTIONS("/burn", CORSMiddleware())
	protectedUsers.POST("/burn", blockchain.Burn)

	protectedUsers.OPTIONS("/transaction-history", CORSMiddleware())
	protectedUsers.GET("/transaction-history", controllers.GetTransactions)

	protectedUsers.OPTIONS("/balance", CORSMiddleware())
	protectedUsers.GET("/balance", controllers.Balance)

	r.Run(":8080")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
