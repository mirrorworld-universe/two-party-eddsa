package eddsa

type KeyAgg struct {
	Apk  Ed25519Point
	Hash Ed25519Scalar
}

func (ka *KeyAgg) ToString() string {
	return "{" +
		"\"Apk\": " + ka.Apk.ToString() + ", " +
		"\"Hash\": " + ka.Hash.ToString() +
		"}"
}
