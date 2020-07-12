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
	router.POST("/newBlock", newBlockHandler)

	router.GET("/getPeers", getPeersHandler)
	router.POST("/newPeer", newPeerHandler)

	router.POST("/addAlias", addAliasHandler)
	router.POST("/getPublicKey", getPublicKeyHandler)

	router.GET("/getPendingTransactions", pendingTransactionsHandler)
	router.POST("/newTransaction", newTransactionsHandler)

	router.POST("/getUnusedOutputs", getUnusedOutputsHandler)
}
