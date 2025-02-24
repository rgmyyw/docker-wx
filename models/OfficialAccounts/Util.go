package OfficialAccounts

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
)

type DefaultParam struct {
	Wxid  string
	Appid string
}

type ReadParam struct {
	Wxid string
	Url  string
}

type GetkeyParam struct {
	Wxid string
	Url  string
	Appid string
}

func CreateRandomNumber() uint64 {
	var numbers = []byte{1, 2, 3, 4, 5, 7, 8, 9}
	var container string
	length := bytes.NewReader(numbers).Len()

	for i := 1; i <= 7; i++ {
		random, err := rand.Int(rand.Reader, big.NewInt(int64(length)))
		if err != nil {

		}
		container += fmt.Sprintf("%d", numbers[random.Int64()])
	}

	B := "340" + container
	b, _ := strconv.Atoi(B)
	return uint64(b)
}
