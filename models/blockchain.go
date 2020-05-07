package models

type Blockchain struct {
	Chain                    []Block
	UnusedTransactionOutputs OutputList
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
	i := blockchain.UnusedTransactionOutputs.FindOutputIndex(transaction.Input)
	if i == -1 {
		return false
	}

	totalAmount := transaction.Outputs.GetSumOfAmounts()
	return totalAmount <= blockchain.UnusedTransactionOutputs[i].Amount
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
