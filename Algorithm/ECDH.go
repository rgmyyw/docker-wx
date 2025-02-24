package Algorithm

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
)

func GetECDH415Key() (privKey []byte, pubKey []byte) {
	privKey = nil
	pubKey = nil
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub := &priv.PublicKey
	pubKey = elliptic.Marshal(pub.Curve, pub.X, pub.Y)
	privKey = priv.D.Bytes()
	return
}

func DoECDH415Key(privD, pubData []byte) []byte {
	X, Y := elliptic.Unmarshal(elliptic.P256(), pubData)
	if X == nil || Y == nil {
		return nil
	}
	x, _ := elliptic.P256().ScalarMult(X, Y, privD)
	return x.Bytes()
}

func DoECDH713Key(privD, pubData []byte) []byte {
	X, Y := elliptic.Unmarshal(elliptic.P224(), pubData)
	if X == nil || Y == nil {
		return []byte{}
	}
	x, _ := elliptic.P224().ScalarMult(X, Y, privD)
	return x.Bytes()
}

func DoECDH713(pub, priv []byte) []byte {
	curve := elliptic.P224()
	x, y := elliptic.Unmarshal(curve, pub)
	if x == nil {
		return nil
	}
	xShared, _ := curve.ScalarMult(x, y, priv)
	sharedKey := make([]byte, (curve.Params().BitSize+7)>>3)
	xBytes := xShared.Bytes()
	copy(sharedKey[len(sharedKey)-len(xBytes):], xBytes)

	dh := md5.Sum(sharedKey)
	return dh[:]
}

func GetEcdh713Key() (privKey []byte, pubKey []byte) {
	privKey = nil
	pubKey = nil
	priv, _ := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	pub := &priv.PublicKey
	pubKey = elliptic.Marshal(pub.Curve, pub.X, pub.Y)
	privKey = priv.D.Bytes()
	return
}

func Ecdh(curve elliptic.Curve, pub, priv []byte) []byte {

	x, y := elliptic.Unmarshal(curve, pub)
	if x == nil {
		return nil
	}

	xShared, _ := curve.ScalarMult(x, y, priv)
	sharedKey := make([]byte, (curve.Params().BitSize+7)>>3)
	xBytes := xShared.Bytes()
	copy(sharedKey[len(sharedKey)-len(xBytes):], xBytes)

	dh := sha256.Sum256(sharedKey)
	return dh[:]
}
