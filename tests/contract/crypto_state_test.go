package contract

import (
	"context"
	"io"
)

type cryptoStateKey struct{}

type cryptoState struct {
	cryptoAlgorithm   string
	cipherKey         string
	cryptoFeaturePath string
	randomIv          bool
	result            []byte
	resultReader      io.Reader
	err               error
}

func getCryptoState(ctx context.Context) *cryptoState {
	return ctx.Value(cryptoStateKey{}).(*cryptoState)
}

func newCryptoState(cryptoFeaturePath string) *cryptoState {
	return &cryptoState{cryptoFeaturePath: cryptoFeaturePath, randomIv: false}
}
