package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dryairship/IITKBucks/config"
	"github.com/dryairship/IITKBucks/logger"
	"github.com/dryairship/IITKBucks/models"
)

var currentMinerChannel chan bool

func performNoobInitialization() {
	c := make(chan bool)
	go tryToAddPeers(c)
	addedAPeer := <-c
	if !addedAPeer {
		logger.Fatal("[Controllers/Init] [FATAL] No peers added.")
	}
	blocksReceived := askAndAddBlocks()
	transactionsReceived := askAndAddTransactions()
	if !blocksReceived || !transactionsReceived {
		logger.Fatal("[Controllers/Init] [FATAL] Data init failed.")
	}
	currentMinerChannel = make(chan bool)
	go startMining()
}

func performProInitialization() {
	genesisBlock := models.NewGenesisBlock()

	currentMinerChannel = make(chan bool)
	go mineBlock(genesisBlock)
}

func PerformInitialization() {
	if config.IS_PRO {
		performProInitialization()
	} else {
		performNoobInitialization()
	}
}

func askAndAddBlocks() bool {
	for i := 0; ; i++ {
		response, err := http.Get(fmt.Sprintf("%s/getBlock/%d", peers[0], i))
		if err != nil {
			logger.Println(logger.RareError, "[Controllers/Init] [WARN] go HTTP error while asking for blocks. Peer:", peers[0], ", Block:", i, ", Error:", err)
			return i > 0
		}

		if response.StatusCode != 200 {
			break
		}

		defer response.Body.Close()
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			logger.Println(logger.RareError, "[Controllers/Init] [WARN] Cannot read getBlock response body. Peer:", peers[0], ", Block:", i, ", Error:", err)
			return i > 0
		}

		block, err := models.BlockFromByteArray(bodyBytes)
		if err != nil {
			logger.Println(logger.RareError, "[Controllers/Init] [ERROR] Could not convert received bytes to block. Peer:", peers[0], ", Block:", i, ", Error:", err)
			return i > 0
		}

		isValid, err := models.Blockchain().IsBlockValid(&block)
		if !isValid || err != nil {
			logger.Println(logger.RareError, "[Controllers/Init] [WARN] Received block is not valid. Peer:", peers[0], ", Block:", i, ", Error:", err)
			return i > 0
		}

		models.Blockchain().ProcessBlock(block)
		models.Blockchain().AppendBlock(block)
		logger.Printf(logger.MajorEvent, "[Controllers/Init] [INFO] Successfully added block %d to the blockchain\n", i)
	}
	return true
}

func askAndAddTransactions() bool {
	response, err := http.Get(fmt.Sprintf("%s/getPendingTransactions", peers[0]))
	if err != nil {
		logger.Println(logger.RareError, "[Controllers/Init] [WARN] go HTTP error while asking for pending transactions. Peer:", peers[0], ", Error:", err)
		return false
	}

	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Println(logger.RareError, "[Controllers/Init] [WARN] Cannot read getPendingTransactions response body. Peer:", peers[0], ". Error:", err)
		return false
	}

	var txnList []models.TransactionRequestBody
	err = json.Unmarshal(bodyBytes, &txnList)
	if err != nil {
		logger.Println(logger.RareError, "[Controllers/Init] [WARN] Cannot unmarshal getPeers response body. Peer:", peers[0], ", Error:", err)
		return false
	}

	for _, txnBody := range txnList {
		txn, err := txnBody.ToTransaction()
		if err != nil {
			logger.Println(logger.RareError, "[Controllers/Init] [WARN] Cannot convert transaction body to transaction. Body: ", txnBody)
			continue
		}

		isValid, _ := models.Blockchain().IsTransactionValid(&txn)
		if isValid {
			models.Blockchain().AddTransaction(txn)
		} else {
			logger.Println(logger.RareError, "[Controllers/Init] [WARN] Received transaction is not valid")
		}
	}
	return true
}
