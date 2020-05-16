package models_test

import (
	"testing"

	"github.com/dryairship/IITKBucks/models"
)

func TestUserToByteArray(t *testing.T) {
	user := models.User("dryairship")
	expected := []byte("dryairship")

	ans := user.ToByteArray()

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
