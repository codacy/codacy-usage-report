package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoiningArrayWithComma(t *testing.T) {
	assert := assert.New(t)

	data := []uint{11, 22, 33, 44}
	expectResult := "11,22,33,44"

	result := JoinUintArray(data, ",")

	assert.Equal(expectResult, result)
}

func TestJoiningEmptyArrayWithComma(t *testing.T) {
	assert := assert.New(t)

	data := []uint{}
	expectResult := ""

	result := JoinUintArray(data, ",")

	assert.Equal(expectResult, result)
}
