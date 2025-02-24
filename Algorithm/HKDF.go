package Algorithm

import (
	"crypto/hmac"
	"crypto/sha256"
	"hash"
)

func HybridHkdfExpand(prikey []byte, salt []byte, info []byte, outLen int) []byte {
	h := hmac.New(sha256.New, prikey)
	h.Write(salt)
	T := h.Sum(nil)
	return HkdfExpand(sha256.New, T, info, outLen)
}

func Hkdf_Expand(h func() hash.Hash, prk, info []byte, outLen int) []byte {
	out := []byte{}
	T := []byte{}
	i := byte(1)
	for len(out) < outLen {
		block := append(T, info...)
		block = append(block, i)

		h := hmac.New(h, prk)
		h.Write(block)

		T = h.Sum(nil)
		out = append(out, T...)
		i++
	}
	return out[:outLen]
}

func HkdfExpand(h func() hash.Hash, prk, info []byte, outLen int) []byte {
	out := []byte{}
	T := []byte{}
	i := byte(1)
	for len(out) < outLen {
		block := append(T, info...)
		block = append(block, i)

		h := hmac.New(h, prk)
		h.Write(block)

		T = h.Sum(nil)
		out = append(out, T...)
		i++
	}
	return out[:outLen]
}
