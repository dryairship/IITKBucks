package models_test

import (
	"testing"

	"github.com/dryairship/IITKBucks/models"
)

func TestOutputToByteArray(t *testing.T) {
	output := models.Output{
		Recipient: models.User("dryairship"),
		Amount:    models.Coins(3648),
	}
	expected := []byte{0, 0, 0, 0, 0, 0, 14, 64, 0, 0, 0, 10, 100, 114, 121, 97, 105, 114, 115, 104, 105, 112}

	ans := output.ToByteArray()

	if len(ans) != len(expected) {
		t.Errorf("Byte Array Length Mismatch: Length of result = {%d}, length of expected result = {%d}", len(ans), len(expected))
		return
	}

	for i := range ans {
		if ans[i] != expected[i] {
			t.Errorf("Byte Array Value Mismatch: ans[{%d}] = {%d}, expected[{%d}] = {%d}", i, ans[i], i, expected[i])
			return
		}
	}
}

func TestOutputListToByteArray(t *testing.T) {
	outputList := models.OutputList{
		models.Output{
			Recipient: models.User("dryairship"),
			Amount:    models.Coins(3648),
		},
		models.Output{
			Recipient: models.User("theothershivangi"),
			Amount:    models.Coins(69),
		},
		models.Output{
			Recipient: models.User("AaryanS941"),
			Amount:    models.Coins(1),
		},
	}

	expected := []byte{0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 14, 64, 0, 0, 0, 10, 100, 114, 121, 97, 105, 114, 115, 104, 105, 112,
		0, 0, 0, 0, 0, 0, 0, 69, 0, 0, 0, 16, 116, 104, 101, 111, 116, 104, 101, 114, 115, 104, 105, 118, 97, 110, 103,
		105, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 10, 65, 97, 114, 121, 97, 110, 83, 57, 52, 49}

	ans := outputList.ToByteArray()

	if len(ans) != len(expected) {
		t.Errorf("Byte Array Length Mismatch: Length of result = {%d}, length of expected result = {%d}", len(ans), len(expected))
		return
	}

	for i := range ans {
		if ans[i] != expected[i] {
			t.Errorf("Byte Array Value Mismatch: ans[{%d}] = {%d}, expected[{%d}] = {%d}", i, ans[i], i, expected[i])
			return
		}
	}
}

func TestOutputGetSumOfAmounts(t *testing.T) {
	outputList := models.OutputList{
		models.Output{
			Recipient: models.User("dryairship"),
			Amount:    models.Coins(1000),
		},
		models.Output{
			Recipient: models.User("theothershivangi"),
			Amount:    models.Coins(100),
		},
		models.Output{
			Recipient: models.User("AaryanS941"),
			Amount:    models.Coins(10),
		},
	}

	ans := outputList.GetSumOfAmounts()

	if ans != models.Coins(1110) {
		t.Error("Incorrect sum of amounts")
	}
}
