package models

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"

	"github.com/dryairship/IITKBucks/logger"
)

type Transaction struct {
	Id      ID         `json:"id"`
	Inputs  InputList  `json:"inputs"`
	Outputs OutputList `json:"outputs"`
}

type TransactionRequestBody struct {
	Inputs  InputListRequestBody  `json:"inputs" binding:"required"`
	Outputs OutputListRequestBody `json:"outputs" binding:"required"`
}

type TransactionList []Transaction
type TransactionMap map[Hash]Transaction

func (txn *Transaction) AddInput(input Input) {
	txn.Inputs = append(txn.Inputs, input)
}

func (txn *Transaction) AddOutput(output Output) {
	txn.Outputs = append(txn.Outputs, output)
}

func (txn Transaction) ToByteArray() []byte {
	return append(txn.Inputs.ToByteArray(), txn.Outputs.ToByteArray()...)
}

func (txnList TransactionList) ToByteArray() []byte {
	var result []byte

	lenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBytes, uint32(len(txnList)))
	result = append(result, lenBytes...)

	for i := range txnList {
		txn := txnList[i].ToByteArray()
		binary.BigEndian.PutUint32(lenBytes, uint32(len(txn)))
		result = append(result, lenBytes...)
		result = append(result, txn...)
	}

	return result
}

func (txn *Transaction) CalculateHash() Hash {
	hash := sha256.Sum256(txn.ToByteArray())
	txn.Id = hash
	return hash
}

func (txn Transaction) CalculateOutputDataHash() Hash {
	return sha256.Sum256(txn.Outputs.ToByteArray())
}

func TransactionFromByteArray(data []byte) (Transaction, error) {
	inputList, bytesRead, err := InputListFromByteArray(data)
	if err != nil {
		logger.Println(logger.RareError, "[Models/Transaction] [ERROR] Error while reading input list from byte array")
		return Transaction{}, err
	}

	outputList, _, err := OutputListFromByteArray(data[bytesRead:])
	if err != nil {
		logger.Println(logger.RareError, "[Models/Transaction] [ERROR] Error while reading output list from byte array")
		return Transaction{}, err
	}

	txn := Transaction{
		Inputs:  inputList,
		Outputs: outputList,
	}

	txn.CalculateHash()
	return txn, nil
}

func TransactionListFromByteArray(data []byte) (TransactionList, error) {
	if len(data) < 4 {
		logger.Println(logger.RareError, "[Models/Transaction] [ERROR] TransactionList has less than 4 bytes")
		return nil, ERROR_INSUFFICIENT_DATA
	}

	numTransactions := binary.BigEndian.Uint32(data[:4])
	currentOffset := 4

	var list TransactionList
	for i := uint32(0); i < numTransactions; i++ {
		if len(data) < currentOffset+4 {
			logger.Printf(logger.RareError, "[Models/Transaction] [ERROR] TransactionList does not have enough data for %d transactions\n", numTransactions)
			return nil, ERROR_INSUFFICIENT_DATA
		}

		txnSize := int(binary.BigEndian.Uint32(data[currentOffset : currentOffset+4]))
		currentOffset += 4
		if len(data) < currentOffset+txnSize {
			logger.Printf(logger.RareError, "[Models/Transaction] [ERROR] TransactionList does not have enough data for %d transactions\n", numTransactions)
			return nil, ERROR_INSUFFICIENT_DATA
		}

		txn, err := TransactionFromByteArray(data[currentOffset : currentOffset+txnSize])
		if err != nil {
			logger.Printf(logger.RareError, "[Models/Transaction] [ERROR] Transaction %d has error. Err: %v\n", i, err)
			return nil, err
		}
		list = append(list, txn)
		currentOffset += txnSize
	}
	return list, nil
}

func (txnRequestBody *TransactionRequestBody) ToTransaction() (Transaction, error) {
	inputs, err := txnRequestBody.Inputs.ToInputList()
	if err != nil {
		logger.Println(logger.RareError, "[Models/Transaction] [ERROR] Could not convert JSON inputs to inputlist")
		return Transaction{}, err
	}

	outputs, err := txnRequestBody.Outputs.ToOutputList()
	if err != nil {
		logger.Println(logger.RareError, "[Models/Transaction] [ERROR] Could not convert JSON outputs to outputlist")
		return Transaction{}, err
	}

	txn := Transaction{
		Inputs:  inputs,
		Outputs: outputs,
	}

	_ = txn.CalculateHash()

	return txn, nil
}

func (txnMap TransactionMap) MarshalJSON() ([]byte, error) {
	if len(txnMap) == 0 {
		return []byte("[]"), nil
	}
	var list TransactionList
	for _, txn := range txnMap {
		list = append(list, txn)
	}
	return json.Marshal(list)
}
