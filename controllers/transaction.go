package controllers

import (
	"github.com/gin-gonic/gin"

	"github.com/dryairship/IITKBucks/models"
)

func pendingTransactionsHandler(c *gin.Context) {
	if models.Blockchain().PendingTransactions != nil {
		c.JSON(200, models.Blockchain().PendingTransactions)
	} else {
		c.JSON(200, make([]int, 0))
	}
}

func newTransactionsHandler(c *gin.Context) {
	var body models.TransactionRequestBody
	err := c.BindJSON(&body)
	if err != nil {
		c.String(400, "Invalid JSON request body")
		return
	}

	txn, err := body.ToTransaction()
	if err != nil {
		c.String(400, "JSON request body could not be converted to a Transaction object")
		return
	}

	valid, _ := models.Blockchain().IsTransactionValid(&txn)
	if !valid {
		c.String(400, "Invalid transaction")
		return
	}

	models.Blockchain().AddTransaction(txn)
	c.String(200, "Transaction successfully added to list, awaiting confirmation")
}
