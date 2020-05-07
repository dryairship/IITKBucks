package models

type Output struct {
	Id        ID    `json:"id"`
	Recipient User  `json:"recipient"`
	Amount    Coins `json:"amount"`
}
type OutputList []Output

func (outputList OutputList) FindOutputIndex(outputId ID) int {
	for i := range outputList {
		if outputList[i].Id == outputId {
			return i
		}
	}
	return -1
}

func (outputList *OutputList) DeleteOutputAtIndex(i int) {
	(*outputList)[i] = (*outputList)[len(*outputList)-1]
	*outputList = (*outputList)[:len(*outputList)-1]
}

func (outputList OutputList) GetSumOfAmounts() Coins {
	var totalCoins Coins = 0
	for i := range outputList {
		totalCoins += outputList[i].Amount
	}
	return totalCoins
}
