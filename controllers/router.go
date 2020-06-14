package controllers

import (
	"github.com/gin-gonic/gin"
)

func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func SetUpRoutes(router *gin.Engine) {
	router.GET("/ping", pingHandler)
	router.Static("/getBlock", "./blocks")
	router.GET("/getPendingTransactions", pendingTransactionsHandler)
	router.POST("/newTransaction", newTransactionsHandler)
}
