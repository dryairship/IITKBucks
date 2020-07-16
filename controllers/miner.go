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

	pendingTxns := models.Blockchain().PendingTransactions
	unusedOutputs := models.Blockchain().UnusedTransactionOutputs

	coinbaseTxn := createCoinbaseTransaction()
	currentTxns := models.TransactionList{coinbaseTxn}
	currentSize := len(coinbaseTxn.ToByteArray())

	var pair models.TransactionIdIndexPair
	for _, txn := range pendingTxns {
		inputSum := models.Coins(0)
		isTxnValid := true
		for _, input := range txn.Inputs {
			pair.TransactionId = input.TransactionId
			pair.Index = input.OutputIndex
			input, exists := unusedOutputs[pair]
			if !exists {
				isTxnValid = false
				break
			}
			inputSum += input.Amount
		}

		if isTxnValid {
			if currentSize+len(txn.ToByteArray()) <= 1000000 {
				coinbaseTxn.Outputs[0].Amount += inputSum - txn.Outputs.GetSumOfAmounts()
				currentTxns = append(currentTxns, txn)
			} else {
				break
			}
		}
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
			logger.Printf(logger.MajorEvent, "[Controllers/Miner] [INFO] New block mined! Index: %d, Timestamp: %d, Nonce: %d\n", block.Index, block.Timestamp, block.Nonce)
			select {
			case currentMinerChannel <- true:
				logger.Println(logger.MinorEvent, "[Controllers/Miner] [DEBUG] Someone received the new block message though channel")
			default:
			}
			logger.Printf(logger.MinorEvent, "[Controllers/Miner] [DEBUG] Channel closed.")
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
