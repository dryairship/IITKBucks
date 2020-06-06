package models

type Blockchain struct {
	Chain                    []Block
	UnusedTransactionOutputs OutputMap
	PendingTransactions      TransactionList
	CurrentTarget            Hash
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		Chain: []Block{*NewGenesisBlock()},
	}
}

func (blockchain *Blockchain) AppendBlock(block Block) {
	blockchain.Chain = append(blockchain.Chain, block)
}

func (blockchain Blockchain) IsTransactionValid(transaction *Transaction) bool {
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

func (blockchain Blockchain) IsBlockValid(block *Block) bool {
	if block.Index != len(blockchain.Chain)+1 {
		return false
	}

	parentIndex := block.Index - 1

	if blockchain.Chain[parentIndex].Timestamp > block.Timestamp {
		return false
	}

	if blockchain.Chain[parentIndex].GetHash() != block.ParentHash {
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
