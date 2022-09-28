package ed25519

type ExpendedPrivateKey struct {
	Prefix     Ed25519Scalar
	PrivateKey Ed25519Scalar
}

type Keypair struct {
	PublicKey          Ed25519Point
	ExtendedPrivateKey ExpendedPrivateKey
}
