package contract

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/pubnub/go/v7/crypto"
	"io"
	"os"
	"strings"
)

func cryptor(ctx context.Context, cryptor string) error {
	cryptoState := getCryptoState(ctx)
	cryptoState.cryptorNames = append(cryptoState.cryptorNames, cryptor)
	return nil
}

func cryptoModuleWithRegisteredcryptors(ctx context.Context, cryptor1 string, cryptor2 string) error {
	cryptoState := getCryptoState(ctx)
	cryptoState.cryptorNames = append(cryptoState.cryptorNames, cryptor1, cryptor2)
	return nil
}

func decryptedFileContentEqualToFileContent(ctx context.Context, filename string) error {
	cryptoState := getCryptoState(ctx)

	file, e := os.Open(cryptoState.cryptoFeaturePath + "/assets/" + filename)
	if e != nil {
		return e
	}

	fileContent, e := io.ReadAll(file)
	if e != nil {
		return e
	}

	if cryptoState.result != nil {
		if !bytes.Equal(cryptoState.result, fileContent) {
			return errors.New("decrypted file content not equal to file content")
		}
	} else if cryptoState.resultReader != nil {
		resultContent, e := io.ReadAll(cryptoState.resultReader)
		if e != nil {
			return e
		}
		if !bytes.Equal(resultContent, fileContent) {
			return errors.New("decrypted file content not equal to file content")
		}
	} else {
		return errors.New("no result")
	}
	return nil
}

func encryptedFileSuccessfullyDecryptedByLegacyCodeWithCipherKeyAndVector(ctx context.Context, cipherKey string, iv string) error {
	riv := randomIv(iv)
	cryptoState := getCryptoState(ctx)

	cryptor, e := createCryptor("legacy", cipherKey, riv)
	if e != nil {
		return e
	}
	module := crypto.NewCryptoModule(cryptor, []crypto.Cryptor{cryptor})
	if e != nil {
		return e
	}
	if cryptoState.result != nil {
		_, err := module.Decrypt(cryptoState.result)
		if err != nil {
			return err
		}
	} else if cryptoState.resultReader != nil {
		_, err := module.DecryptStream(bufio.NewReader(cryptoState.resultReader))
		if err != nil {
			return err
		}
	} else {
		return errors.New("no result")
	}

	return nil
}

func iDecryptFileAs(ctx context.Context, filename string, decryptionType string) error {

	cryptoState := getCryptoState(ctx)
	module, e := cryptoState.createModule()
	if e != nil {
		return e
	}
	file, e := os.Open(cryptoState.cryptoFeaturePath + "/assets/" + filename)
	if e != nil {
		return e
	}

	if decryptionType == "stream" {
		cryptoState.resultReader, e = module.DecryptStream(bufio.NewReader(file))
		if e != nil {
			return e
		}

	} else {
		fileContent, e := io.ReadAll(file)
		if e != nil {
			return e
		}
		cryptoState.result, e = module.Decrypt(fileContent)
		if e != nil {
			return e
		}
	}
	return nil
}

func iDecryptFile(ctx context.Context, filename string) error {

	cryptoState := getCryptoState(ctx)

	module, e := cryptoState.createModule()
	if e != nil {
		return e
	}
	file, e := cryptoState.openAssetFile(filename)
	if e != nil {
		return e
	}

	_, e = module.DecryptStream(file)
	if e != nil {
		cryptoState.err = e
	}
	_ = file.Close()
	return nil
}

func createCryptor(name string, cipherKey string, randomIv bool) (crypto.Cryptor, error) {
	if name == "acrh" {
		return crypto.NewAesCbcCryptor(cipherKey)
	} else if name == "legacy" {
		return crypto.NewLegacyCryptor(cipherKey, randomIv)
	} else {
		return nil, fmt.Errorf("unknown crypto algorithm %s", name)
	}
}

func iEncryptFileAs(ctx context.Context, filename string, encryptionType string) error {
	cryptoState := getCryptoState(ctx)
	module, e := cryptoState.createModule()
	if e != nil {
		return e
	}
	file, e := cryptoState.openAssetFile(filename)
	if e != nil {
		return e
	}
	if encryptionType == "stream" {
		cryptoState.resultReader, e = module.EncryptStream(bufio.NewReader(file))
		if e != nil {
			return e
		}
	} else {
		content, e := io.ReadAll(file)
		if e != nil {
			return e
		}
		cryptoState.result, e = module.Encrypt(content)
		if e != nil {
			return e
		}

	}
	_ = file.Close()
	return nil
}

func iReceiveDecryptionError(ctx context.Context) error {
	cryptoState := getCryptoState(ctx)
	if cryptoState.err != nil {
		if strings.HasPrefix(cryptoState.err.Error(), "decryption error") {
			return nil
		} else {
			return cryptoState.err
		}
	} else {
		return errors.New("expected error")
	}
}

func iReceiveSuccess(ctx context.Context) error {
	cryptoState := getCryptoState(ctx)
	if cryptoState.err != nil {
		return cryptoState.err
	}
	return nil
}

func iReceiveUnknownCryptorError(ctx context.Context) error {
	cryptoState := getCryptoState(ctx)
	if cryptoState.err != nil {
		if strings.Contains(cryptoState.err.Error(), "unknown crypto error") {
			return nil
		} else {
			return cryptoState.err
		}
	} else {
		return errors.New("expected error")
	}
}

func withCipherKey(ctx context.Context, cipherKey string) error {
	cryptoState := getCryptoState(ctx)
	cryptoState.cipherKey = cipherKey
	return nil
}

func randomIv(iv string) bool {
	if iv == "constant" {
		return false
	} else if iv == "random" {
		return true
	} else {
		return false
	}
}

func withVector(ctx context.Context, iv string) error {
	cryptoState := getCryptoState(ctx)
	cryptoState.randomIv = randomIv(iv)
	return nil
}
