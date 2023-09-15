package contract

import (
	"context"
	"github.com/pubnub/go/v7/crypto"
	"io"
	"os"
)

type cryptoStateKey struct{}

type cryptoState struct {
	cryptorNames      []string
	cipherKey         string
	cryptoFeaturePath string
	randomIv          bool
	result            []byte
	resultReader      io.Reader
	err               error
}

func (c *cryptoState) createModule() (crypto.CryptoModule, error) {
	cryptors := make([]crypto.Cryptor, len(c.cryptorNames))
	var e error
	for i, cryptor := range c.cryptorNames {
		cryptors[i], e = createCryptor(cryptor, c.cipherKey, c.randomIv)
		if e != nil {
			return nil, e
		}
	}
	module := crypto.NewCryptoModule(cryptors[0], cryptors[1:])
	return module, nil
}

func (c *cryptoState) openAssetFile(filename string) (io.ReadCloser, error) {
	return os.Open(c.cryptoFeaturePath + "/assets/" + filename)
}

func getCryptoState(ctx context.Context) *cryptoState {
	return ctx.Value(cryptoStateKey{}).(*cryptoState)
}

func newCryptoState(cryptoFeaturePath string) *cryptoState {
	return &cryptoState{cryptoFeaturePath: cryptoFeaturePath, randomIv: false}
}
