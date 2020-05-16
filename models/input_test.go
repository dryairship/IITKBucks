package models_test

import (
	"testing"

	"github.com/dryairship/IITKBucks/models"
)

func TestInputToByteArray(t *testing.T) {
	txid := [32]byte{82, 253, 252, 7, 33, 130, 101, 79, 22, 63, 95, 15, 154, 98, 29, 114, 149, 102, 199, 77, 16, 3, 124, 77, 123, 187, 4, 7, 209, 226, 198, 73}
	signature, _ := models.SignatureFromHexString("81855ad8681d0d86d1e91e00167939cb6694d2c422acd208a0072939487f6999eb9d18a44784045d")
	expected := []byte{82, 253, 252, 7, 33, 130, 101, 79, 22, 63, 95, 15, 154, 98, 29, 114, 149, 102, 199, 77, 16, 3, 124, 77, 123, 187, 4, 7, 209, 226, 198,
		73, 0, 0, 0, 3, 0, 0, 0, 40, 129, 133, 90, 216, 104, 29, 13, 134, 209, 233, 30, 0, 22, 121, 57, 203, 102, 148, 210, 196, 34, 172, 210, 8, 160, 7, 41,
		57, 72, 127, 105, 153, 235, 157, 24, 164, 71, 132, 4, 93}

	input := models.Input{
		TransactionId: txid,
		OutputIndex:   3,
		Signature:     signature,
	}

	ans := input.ToByteArray()

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

func TestInputListToByteArray(t *testing.T) {
	txid1 := [32]byte{82, 253, 252, 7, 33, 130, 101, 79, 22, 63, 95, 15, 154, 98, 29, 114, 149, 102, 199, 77, 16, 3, 124, 77, 123,
		187, 4, 7, 209, 226, 198, 73}
	signature1, _ := models.SignatureFromHexString("81855ad8681d0d86d1e91e00167939cb6694d2c422acd208a0072939487f6999eb9d18a44784045d")
	txid2 := [32]byte{135, 243, 198, 124, 242, 39, 70, 233, 149, 175, 90, 37, 54, 121, 81, 186, 162, 255, 108, 212, 113, 196, 131,
		241, 95, 185, 11, 173, 179, 124, 88, 33}
	signature2, _ := models.SignatureFromHexString("b6d95526a41a9504680b4e7c8b763a1b1d49d4955c8486216325253fec738dd7a9e28bf921119c160f0702448615bbda08313f6a8eb668d20bf50598")
	expected := []byte{0, 0, 0, 2, 82, 253, 252, 7, 33, 130, 101, 79, 22, 63, 95, 15, 154, 98, 29, 114, 149, 102, 199, 77, 16, 3,
		124, 77, 123, 187, 4, 7, 209, 226, 198, 73, 0, 0, 0, 3, 0, 0, 0, 40, 129, 133, 90, 216, 104, 29, 13, 134, 209, 233, 30, 0,
		22, 121, 57, 203, 102, 148, 210, 196, 34, 172, 210, 8, 160, 7, 41, 57, 72, 127, 105, 153, 235, 157, 24, 164, 71, 132, 4,
		93, 135, 243, 198, 124, 242, 39, 70, 233, 149, 175, 90, 37, 54, 121, 81, 186, 162, 255, 108, 212, 113, 196, 131, 241, 95,
		185, 11, 173, 179, 124, 88, 33, 0, 0, 12, 160, 0, 0, 0, 60, 182, 217, 85, 38, 164, 26, 149, 4, 104, 11, 78, 124, 139, 118,
		58, 27, 29, 73, 212, 149, 92, 132, 134, 33, 99, 37, 37, 63, 236, 115, 141, 215, 169, 226, 139, 249, 33, 17, 156, 22, 15,
		7, 2, 68, 134, 21, 187, 218, 8, 49, 63, 106, 142, 182, 104, 210, 11, 245, 5, 152}

	inputList := models.InputList{
		models.Input{
			TransactionId: txid1,
			OutputIndex:   3,
			Signature:     signature1,
		},
		models.Input{
			TransactionId: txid2,
			OutputIndex:   3232,
			Signature:     signature2,
		},
	}

	ans := inputList.ToByteArray()

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
