package controllers

import (
	"github.com/gin-gonic/gin"

	"github.com/dryairship/IITKBucks/models"
)

func pendingTransactionsHandler(c *gin.Context) {
	if models.Blockchain().PendingTransactions != nil {
		c.JSON(200, models.Blockchain().PendingTransactions)
	} else {
		c.JSON(200, gin.H{})
	}
}

func newTransactionsHandler(c *gin.Context) {
	var body models.TransactionRequestBody
	err := c.BindJSON(&body)
	if err != nil {
		_ = c.AbortWithError(400, err)
		return
	}

	txn, err := body.ToTransaction()
	if err != nil {
		_ = c.AbortWithError(400, err)
		return
	}

	valid, _ := models.Blockchain().IsTransactionValid(&txn)
	if !valid {
		c.String(400, "Invalid transaction.")
		return
	}

	models.Blockchain().AddTransaction(txn)
	c.Status(200)
}
