package models

import (
	"github.com/dryairship/IITKBucks/config"
	"github.com/dryairship/IITKBucks/logger"
)

type blockchain struct {
	Chain                    []Block
	UnusedTransactionOutputs OutputMap
	PendingTransactions      TransactionMap
	CurrentTarget            Hash
	CurrentBlockReward       Coins
	TransactionAdded         chan bool
	UserOutputs              map[User][]TransactionIdIndexPair
}

var blockchainInstance *blockchain

func Blockchain() *blockchain {
	if blockchainInstance == nil {
		target, err := HashFromHexString(config.INITIAL_TARGET)
		if err != nil {
			logger.Fatal("[Models/Blockchain] [FATAL] Could not read initial target. Err:", err)
		}
		blockchainInstance = &blockchain{
			Chain:                    make([]Block, 0),
			UnusedTransactionOutputs: make(OutputMap),
			PendingTransactions:      make(TransactionMap),
			CurrentTarget:            target,
			CurrentBlockReward:       Coins(config.INITIAL_BLOCK_REWARD),
			TransactionAdded:         make(chan bool),
			UserOutputs:              make(map[User][]TransactionIdIndexPair),
		}
	}
	return blockchainInstance
}

func (blockchain *blockchain) AppendBlock(block Block) {
	block.SaveToFile()
	block.Transactions = nil
	blockchain.Chain = append(blockchain.Chain, block)
}

func (blockchain *blockchain) IsTransactionPending(transactionHash Hash) bool {
	_, exists := blockchain.PendingTransactions[transactionHash]
	return exists
}

func (blockchain *blockchain) IsTransactionValid(transaction *Transaction) (bool, Coins) {
	outputDataHash := transaction.CalculateOutputDataHash()
	sumOfInputs := Coins(0)

	var pair TransactionIdIndexPair
	for _, input := range transaction.Inputs {
		pair.TransactionId = input.TransactionId
		pair.Index = input.OutputIndex

		output, exists := blockchain.UnusedTransactionOutputs[pair]
		if !exists {
			logger.Println(logger.RareError, "[Models/Blockchain] [ERROR] TXIDInputPair does not exist in unused outputs. Pair: ", pair)
			return false, 0
		}
		if !input.Signature.Unlock(&output, &pair, &outputDataHash) {
			logger.Println(logger.RareError, "[Models/Blockchain] [ERROR] Signature verification failed. Input:", input, ", Output:", output)
			return false, 0
		}

		sumOfInputs += output.Amount
	}

	sumOfOutputs := transaction.Outputs.GetSumOfAmounts()
	return sumOfOutputs <= sumOfInputs, sumOfInputs - sumOfOutputs
}

func (blockchain *blockchain) IsBlockValid(block *Block) (bool, error) {
	if block.Index != uint32(len(blockchain.Chain)+1) {
		logger.Println(logger.RareError, "[Models/Blockchain] [ERROR] Block has incorrect index. Expected:", len(blockchain.Chain)+1, ", Found:", block.Index)
		return false, ERROR_INCORRECT_INDEX
	}

	parentIndex := block.Index - 1

	if blockchain.Chain[parentIndex].Timestamp > block.Timestamp {
		logger.Println(logger.RareError, "[Models/Blockchain] [ERROR] Block's timestamp is less than parent's timestamp. Parent:", blockchain.Chain[parentIndex].Timestamp, ", Block:", block.Timestamp)
		return false, ERROR_LESS_TIMESTAMP
	}

	if blockchain.Chain[parentIndex].GetHash() != block.ParentHash {
		logger.Println(logger.RareError, "[Models/Blockchain] [ERROR] Parent hash mismatch. Expected:", blockchain.Chain[parentIndex].GetHash(), ", Found:", block.ParentHash)
		return false, ERROR_PARENT_HASH_MISMATCH
	}

	if block.Target != blockchain.CurrentTarget {
		logger.Println(logger.RareError, "[Models/Blockchain] [ERROR] Target mismatch. Expected:", blockchain.CurrentTarget, ", Found:", block.Target)
		return false, ERROR_TARGET_MISMATCH
	}

	if block.BodyHash != block.GetBodyHash() {
		logger.Println(logger.RareError, "[Models/Blockchain] [ERROR] Body Hash mismatch. Calculated hash:", block.GetBodyHash(), ", Stored hash:", block.BodyHash)
		return false, ERROR_INCORRECT_BODY_HASH
	}

	if !block.GetHash().IsLessThan(blockchain.CurrentTarget) {
		logger.Println(logger.RareError, "[Models/Blockchain] [ERROR] Block hash is not less than target. Block hash:", block.GetHash(), ", Target:", blockchain.CurrentTarget)
		return false, ERROR_ABOVE_TARGET
	}

	var totalFee Coins
	for i := range block.Transactions {
		if i == 0 {
			continue
		}
		valid, txnFee := blockchain.IsTransactionValid(&block.Transactions[i])
		if !valid {
			logger.Printf(logger.RareError, "[Models/Blockchain] [ERROR] Transaction %d in the block is invalid", i)
			return false, ERROR_CONTAINS_INVALID_TRANSACTION
		}
		totalFee += txnFee
	}

	totalFee += blockchain.CurrentBlockReward

	if totalFee < block.Transactions[0].Outputs.GetSumOfAmounts() {
		logger.Println(logger.RareError, "[Models/Blockchain] [ERROR] Coinbase transaction has more coins than expected. Expected:", totalFee, ", Found:", block.Transactions[0].Outputs.GetSumOfAmounts())
		return false, ERROR_INVALID_COINBASE_TRANSCTION
	}

	return true, nil
}

func (blockchain *blockchain) AddTransaction(transaction Transaction) {
	blockchain.TransactionAdded <- true
	blockchain.PendingTransactions[transaction.CalculateHash()] = transaction
}

func (blockchain *blockchain) ProcessBlock(block Block) {
	for _, txn := range block.Transactions {
		delete(blockchain.PendingTransactions, txn.CalculateHash())

		for _, input := range txn.Inputs {
			txidIndexPair := TransactionIdIndexPair{
				TransactionId: input.TransactionId,
				Index:         input.OutputIndex,
			}
			delete(blockchain.UnusedTransactionOutputs, txidIndexPair)
		}

		txidIndexPair := TransactionIdIndexPair{
			TransactionId: txn.Id,
			Index:         0,
		}
		for i, output := range txn.Outputs {
			txidIndexPair.Index = uint32(i)
			blockchain.UnusedTransactionOutputs[txidIndexPair] = output
			if blockchain.UserOutputs[output.Recipient] == nil {
				blockchain.UserOutputs[output.Recipient] = make([]TransactionIdIndexPair, 0)
			}
			blockchain.UserOutputs[output.Recipient] = append(blockchain.UserOutputs[output.Recipient], txidIndexPair)
		}
	}
}
