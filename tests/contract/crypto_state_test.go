package contract

import (
	"context"
	"errors"
	"github.com/pubnub/go/v7/crypto"
	"io"
	"os"
)

type cryptoStateKey struct{}

type cryptoState struct {
	cryptorNames        []string
	cryptorName         string
	cipherKey           string
	cryptoFeaturePath   string
	randomIv            bool
	result              []byte
	resultReader        io.Reader
	err                 error
	cryptoModule        crypto.CryptoModule
	legacyCodeCipherKey string
	legacyCodeRandomIv  bool
}

func (c *cryptoState) getCryptoModule() (crypto.CryptoModule, error) {

	if len(c.cryptorNames) > 0 && c.cryptorNames[0] == "legacy" {
		return crypto.NewLegacyCryptoModule(c.cipherKey, c.randomIv)
	} else if len(c.cryptorNames) > 0 && c.cryptorNames[0] == "acrh" {
		return crypto.NewAesCbcCryptoModule(c.cipherKey, c.randomIv)
	} else if c.cryptorName != "" {
		cryptor, e := createCryptor(c.cryptorName, c.cipherKey, c.randomIv)
		c.cryptoModule = crypto.NewCryptoModule(cryptor, []crypto.Cryptor{})
		return c.cryptoModule, e
	} else {
		return nil, errors.New("I don't know how to create this crypto module")
	}
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
