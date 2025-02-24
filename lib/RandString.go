package lib

import math_rand "math/rand"

func RandSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[math_rand.Intn(len(letters))]
	}
	return string(b)
}
