package controllers

import (
	"github.com/dryairship/IITKBucks/models"
	"github.com/gin-gonic/gin"
)

type unusedOutput struct {
	TransactionId models.Hash  `json:"transactionId"`
	Index         uint32       `json:"index"`
	Amount        models.Coins `json:"amount"`
}

func getUnusedOutputsHandler(c *gin.Context) {
	var body aliasRequestBody
	err := c.BindJSON(&body)
	if err != nil {
		c.String(400, "Invalid JSON request body")
		return
	}

	var publicKey = body.PublicKey
	if body.Alias != "" {
		key, exists := aliasMap[body.Alias]
		if !exists {
			c.String(400, "Alias not asscociated with any public key")
			return
		}
		publicKey = key
	}

	txidIndexPairs, exists := models.Blockchain().UserOutputs[models.User(publicKey)]
	if !exists {
		c.JSON(200, gin.H{
			"unusedOutputs": make([]int, 0),
		})
		return
	}

	var unusedOutputs []unusedOutput
	var validPairs []models.TransactionIdIndexPair

	for _, txidIndexPair := range txidIndexPairs {
		output, exists := models.Blockchain().UnusedTransactionOutputs[txidIndexPair]
		if exists {
			unusedOutputs = append(unusedOutputs,
				unusedOutput{
					Amount:        output.Amount,
					Index:         txidIndexPair.Index,
					TransactionId: txidIndexPair.TransactionId,
				},
			)
			validPairs = append(validPairs, txidIndexPair)
		}
	}

	models.Blockchain().UserOutputs[models.User(publicKey)] = validPairs
	c.JSON(200, gin.H{
		"unusedOutputs": unusedOutputs,
	})
}
