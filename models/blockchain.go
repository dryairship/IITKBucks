package models

type blockchain struct {
	Chain                    []Block
	UnusedTransactionOutputs OutputMap
	PendingTransactions      TransactionList
	CurrentTarget            Hash
}

var blockchainInstance *blockchain

func Blockchain() *blockchain {
	if blockchainInstance == nil {
		blockchainInstance = &blockchain{
			Chain: []Block{*NewGenesisBlock()},
		}
	}
	return blockchainInstance
}

func (blockchain *blockchain) AppendBlock(block Block) {
	blockchain.Chain = append(blockchain.Chain, block)
}

func (blockchain *blockchain) IsTransactionValid(transaction *Transaction) bool {
	outputDataHash := transaction.CalculateOutputDataHash()
	sumOfInputs := uint64(0)

	var pair TransactionIdIndexPair
	for _, input := range transaction.Inputs {
		pair.TransactionId = input.TransactionId
		pair.Index = input.OutputIndex

		output, exists := blockchain.UnusedTransactionOutputs[pair]
		if !exists || !input.Signature.Unlock(&output, &outputDataHash) {
			return false
		}

		sumOfInputs += output.Amount
	}

	sumOfOutputs := transaction.Outputs.GetSumOfAmounts()
	return sumOfOutputs <= sumOfInputs
}

func (blockchain *blockchain) IsBlockValid(block *Block) bool {
	if block.Index != uint32(len(blockchain.Chain)+1) {
		return false
	}

	parentIndex := block.Index - 1

	if blockchain.Chain[parentIndex].Timestamp > block.Timestamp {
		return false
	}

	if blockchain.Chain[parentIndex].GetHash() != block.ParentHash {
		return false
	}

	if block.Target != blockchain.CurrentTarget {
		return false
	}

	if block.BodyHash != block.GetBodyHash() {
		return false
	}

	if !block.GetHash().IsLessThan(blockchain.CurrentTarget) {
		return false
	}

	for i := range block.Transactions {
		if !blockchain.IsTransactionValid(&block.Transactions[i]) {
			return false
		}
	}

	return true
}

func (blockchain *blockchain) AddTransaction(transaction Transaction) {
	blockchain.PendingTransactions = append(blockchain.PendingTransactions, transaction)
}
