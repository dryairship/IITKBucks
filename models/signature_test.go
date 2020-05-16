package models_test

import (
	"testing"

	"github.com/dryairship/IITKBucks/models"
)

func TestSignatureToByteArray(t *testing.T) {
	sign := models.Signature([]byte{129, 133, 90, 216, 104, 29, 13, 134, 209, 233, 30, 0, 22, 121, 57, 203, 102, 148, 210, 196, 34, 172, 210, 8,
		160, 7, 41, 57, 72, 127, 105, 153, 235, 157, 24, 164, 71, 132, 4, 93})
	expected := []byte{129, 133, 90, 216, 104, 29, 13, 134, 209, 233, 30, 0, 22, 121, 57, 203, 102, 148, 210, 196, 34, 172, 210, 8,
		160, 7, 41, 57, 72, 127, 105, 153, 235, 157, 24, 164, 71, 132, 4, 93}

	ans := sign.ToByteArray()

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

func TestSignatureFromHexString(t *testing.T) {
	_, err := models.SignatureFromHexString("aaa")
	if err == nil {
		t.Errorf("No Error while parsing odd length string")
	}

	_, err = models.SignatureFromHexString("aaax")
	if err == nil {
		t.Errorf("No Error while parsing non-hexadecimal string")
	}

	_, err = models.SignatureFromHexString("")
	if err == nil {
		t.Errorf("No Error while parsing empty string")
	}

	signature, err := models.SignatureFromHexString("81855ad8681d0d86d1e91e00167939cb6694d2c422acd208a0072939487f6999eb9d18a44784045d")
	if err != nil {
		t.Errorf("Error while parsing correct string: %v", err)
	}

	expected := []byte{129, 133, 90, 216, 104, 29, 13, 134, 209, 233, 30, 0, 22, 121, 57, 203, 102, 148, 210, 196, 34, 172, 210, 8,
		160, 7, 41, 57, 72, 127, 105, 153, 235, 157, 24, 164, 71, 132, 4, 93}

	ans := signature.ToByteArray()

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

func TestSignatureToString(t *testing.T) {
	signature := models.Signature([]byte{182, 217, 85, 38, 164, 26, 149, 4, 104, 11, 78, 124, 139, 118, 58, 27, 29, 73, 212, 149,
		92, 132, 134, 33, 99, 37, 37, 63, 236, 115, 141, 215, 169, 226, 139, 249, 33, 17, 156, 22, 15, 7, 2,
		68, 134, 21, 187, 218, 8, 49, 63, 106, 142, 182, 104, 210, 11, 245, 5, 152})
	expected := "b6d95526a41a9504680b4e7c8b763a1b1d49d4955c8486216325253fec738dd7a9e28bf921119c160f0702448615bbda08313f6a8eb668d20bf50598"

	ans := signature.String()

	if ans != expected {
		t.Errorf("String mismatch: ans = {%s}, expected = {%s}", ans, expected)
	}
}
