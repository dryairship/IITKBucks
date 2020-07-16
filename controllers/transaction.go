package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dryairship/IITKBucks/logger"
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
	go propagateTransactionToPeers(body)
	c.String(200, "Transaction successfully added to list, awaiting confirmation")
}

func propagateTransactionToPeers(txnBody models.TransactionRequestBody) {
	jsonString, err := json.Marshal(txnBody)
	if err != nil {
		logger.Println(logger.RareError, "[Controllers/Transaction] [ERROR] Could not marshal txn body to json string, Error: ", err)
		return
	}

	buffer := bytes.NewBuffer(jsonString)
	client := &http.Client{}

	count := 0

	for _, peer := range peers {
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/newTransaction", peer), buffer)
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			logger.Println(logger.RareError, "[Controllers/Transaction] [WARN] go HTTP error while builing newTransaction request. Peer: ", peer, ", Error: ", err)
			continue
		}
		req.Header.Set("Content-Type", "application/octet-stream")

		resp, err := client.Do(req)
		if err != nil {
			logger.Println(logger.CommonError, "[Controllers/Transaction] [WARN] go HTTP error while making newTransaction request. Peer: ", peer, ", Error: ", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			var reply []byte
			_, _ = resp.Body.Read(reply)
			logger.Printf(logger.CommonError, "[Controllers/Transaction] [WARN] Peer %s gave %d response on newTransaction request. %s\n", peer, resp.StatusCode, reply)
		} else {
			count++
		}
	}
	logger.Printf(logger.MinorEvent, "[Controllers/Transaction] [DEBUG] Successfully sent newTransaction requests to %d out of %d peers.\n", count, len(peers))
}
