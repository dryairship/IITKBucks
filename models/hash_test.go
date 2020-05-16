package models_test

import (
	"testing"

	"github.com/dryairship/IITKBucks/models"
)

func TestHashToByteArray(t *testing.T) {
	hash := models.Hash([32]byte{129, 133, 90, 216, 104, 29, 13, 134, 209, 233, 30, 0, 22, 121, 57, 203, 102, 148, 210, 196, 34, 172, 210, 8,
		160, 7, 41, 57, 72, 127, 105, 153})
	expected := []byte{129, 133, 90, 216, 104, 29, 13, 134, 209, 233, 30, 0, 22, 121, 57, 203, 102, 148, 210, 196, 34, 172, 210, 8,
		160, 7, 41, 57, 72, 127, 105, 153}

	ans := hash.ToByteArray()

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

func TestHashFromHexString(t *testing.T) {
	_, err := models.HashFromHexString("aaa")
	if err == nil {
		t.Errorf("No Error while parsing odd length string")
	}

	_, err = models.HashFromHexString("aaax")
	if err == nil {
		t.Errorf("No Error while parsing non-hexadecimal string")
	}

	_, err = models.HashFromHexString("")
	if err == nil {
		t.Errorf("No Error while parsing empty string")
	}

	_, err = models.HashFromHexString("81855ad8681d0d86d1e91e00167939cb6694d2c422acd208a0072939487f6999ab")
	if err == nil {
		t.Errorf("No error while parsing string of incorrect length")
	}

	hash, err := models.HashFromHexString("81855ad8681d0d86d1e91e00167939cb6694d2c422acd208a0072939487f6999")
	if err != nil {
		t.Errorf("Error while parsing correct string: %v", err)
	}

	expected := []byte{129, 133, 90, 216, 104, 29, 13, 134, 209, 233, 30, 0, 22, 121, 57, 203, 102, 148, 210, 196, 34, 172, 210, 8,
		160, 7, 41, 57, 72, 127, 105, 153}

	ans := hash.ToByteArray()

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

func TestHashMarshalJSON(t *testing.T) {
	expected := []byte{34, 56, 49, 56, 53, 53, 97, 100, 56, 54, 56, 49, 100, 48, 100, 56, 54, 100, 49, 101, 57, 49, 101, 48,
		48, 49, 54, 55, 57, 51, 57, 99, 98, 54, 54, 57, 52, 100, 50, 99, 52, 50, 50, 97, 99, 100, 50, 48, 56, 97, 48, 48,
		55, 50, 57, 51, 57, 52, 56, 55, 102, 54, 57, 57, 57, 34}
	hash, _ := models.HashFromHexString("81855ad8681d0d86d1e91e00167939cb6694d2c422acd208a0072939487f6999")

	ans, err := hash.MarshalJSON()
	if err != nil {
		t.Errorf("Error while converting hash to JSON: %v", err)
	}

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

func TestHashToString(t *testing.T) {
	hash := models.Hash([32]byte{129, 133, 90, 216, 104, 29, 13, 134, 209, 233, 30, 0, 22, 121, 57, 203, 102, 148, 210, 196,
		34, 172, 210, 8, 160, 7, 41, 57, 72, 127, 105, 153})
	expected := "81855ad8681d0d86d1e91e00167939cb6694d2c422acd208a0072939487f6999"

	ans := hash.String()

	if ans != expected {
		t.Errorf("String mismatch: ans = {%s}, expected = {%s}", ans, expected)
	}
}

func TestHashIsLessThan(t *testing.T) {
	hash0, _ := models.HashFromHexString("0000000000000000000000000000000000000000000000000000000000000000")
	hash1, _ := models.HashFromHexString("0000000000000000000000000000000000000000000000000000000000000001")
	hash1_again, _ := models.HashFromHexString("0000000000000000000000000000000000000000000000000000000000000001")

	if !hash0.IsLessThan(hash1) {
		t.Error("false returned for correct comparison")
	}

	if hash1.IsLessThan(hash0) {
		t.Error("true returned for incorrect comparison")
	}

	if hash1.IsLessThan(hash1_again) {
		t.Error("true returned when comparing equal values")
	}
}
