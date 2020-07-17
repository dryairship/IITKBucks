package controllers

import (
	"time"

	"github.com/dryairship/IITKBucks/logger"
	"github.com/dryairship/IITKBucks/models"
)

func createCoinbaseTransaction() models.Transaction {
	output := models.Output{
		Recipient: models.MyUser,
		Amount:    models.Blockchain().CurrentBlockReward,
	}
	outputList := models.OutputList{output}
	return models.Transaction{
		Outputs: outputList,
	}
}

func createCandidateBlock() models.Block {
	if len(models.Blockchain().PendingTransactions) == 0 {
		logger.Println(logger.MajorEvent, "[Controllers/Miner] [TRACE] Waiting for a transaction to appear.")
		<-models.Blockchain().TransactionAdded
	}

	if len(models.Blockchain().PendingTransactions) == 0 {
		logger.Println(logger.RareError, "[Controllers/Miner] [ERROR] Tried to proceed with block creation without a transaction")
		return createCandidateBlock()
	}

	usedOutputs := make(models.OutputMap)

	coinbaseTxn := createCoinbaseTransaction()
	currentTxns := models.TransactionList{coinbaseTxn}
	currentSize := len(coinbaseTxn.ToByteArray())

	var pair models.TransactionIdIndexPair
	for key, txn := range models.Blockchain().PendingTransactions {
		inputSum := models.Coins(0)
		isTxnValid := true
		for _, input := range txn.Inputs {
			pair.TransactionId = input.TransactionId
			pair.Index = input.OutputIndex
			input, exists := models.Blockchain().UnusedTransactionOutputs[pair]
			_, used := usedOutputs[pair]
			if !exists || used {
				isTxnValid = false
				break
			}
			usedOutputs[pair] = input
			inputSum += input.Amount
		}

		if isTxnValid {
			if currentSize+len(txn.ToByteArray()) <= 1000000 {
				coinbaseTxn.Outputs[0].Amount += inputSum - txn.Outputs.GetSumOfAmounts()
				currentTxns = append(currentTxns, txn)
			} else {
				break
			}
		} else {
			delete(models.Blockchain().PendingTransactions, key)
		}
	}

	if len(currentTxns) == 1 {
		return createCandidateBlock()
	}

	index := len(models.Blockchain().Chain)
	parentHash := models.Blockchain().Chain[index-1].GetHash()
	target := models.Blockchain().CurrentTarget

	logger.Println(logger.MajorEvent, "[Controllers/Miner] [TRACE] Candidate block created with index", index)

	return models.Block{
		Index:        uint32(index),
		ParentHash:   parentHash,
		Target:       target,
		Transactions: currentTxns,
	}
}

func mineBlock(block models.Block) {
	_ = block.CalculateHeader(true)
	target := models.Blockchain().CurrentTarget
	logger.Println(logger.MajorEvent, "[Controllers/Miner] [TRACE] Mining started. Estimated income =", block.Transactions[0].Outputs[0].Amount)
	for i := int64(0); ; i++ {
		block.Nonce = i
		block.Timestamp = time.Now().UnixNano()
		if block.GetHash().IsLessThan(target) {
			logger.Printf(logger.MajorEvent, "[Controllers/Miner] [INFO] New block mined! Index: %d, Timestamp: %d, Nonce: %d, Number of Transaction: %d\n",
				block.Index, block.Timestamp, block.Nonce, len(block.Transactions))
			performPostNewBlockSteps(block)
			return
		}

		select {
		case <-currentMinerChannel:
			logger.Println(logger.MajorEvent, "[Controllers/Miner] [WARN] Mining Interrupted.")
			return
		default:
			continue
		}
	}
}

func startMining() {
	newCandidateBlock := createCandidateBlock()
	mineBlock(newCandidateBlock)
}
