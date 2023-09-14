package contract

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/pubnub/go/v7/crypto"
	"io"
	"os"
	"strings"
)

func cryptoAlgorithm(ctx context.Context, cryptoAlgorithm string) error {
	cryptoState := getCryptoState(ctx)
	cryptoState.cryptoAlgorithm = cryptoAlgorithm
	return nil
}

func cryptorModuleWithRegisteredCryptoAlgorithms() error {
	return godog.ErrPending
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

	algorithm, e := createCryptoAlgorithm("legacy", cipherKey, riv)
	if e != nil {
		return e
	}
	cryptor := crypto.NewCryptor(algorithm)
	if cryptoState.result != nil {
		_, err := cryptor.Decrypt(cryptoState.result)
		if err != nil {
			return err
		}
	} else if cryptoState.resultReader != nil {
		_, err := cryptor.DecryptStream(bufio.NewReader(cryptoState.resultReader))
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

	algorithm, e := createCryptoAlgorithm(cryptoState.cryptoAlgorithm, cryptoState.cipherKey, cryptoState.randomIv)
	if e != nil {
		return e
	}

	cryptor := crypto.NewCryptor(algorithm)
	file, e := os.Open(cryptoState.cryptoFeaturePath + "/assets/" + filename)
	if e != nil {
		return e
	}

	if decryptionType == "stream" {
		cryptoState.resultReader, e = cryptor.DecryptStream(bufio.NewReader(file))
		if e != nil {
			return e
		}

	} else {
		fileContent, e := io.ReadAll(file)
		if e != nil {
			return e
		}
		cryptoState.result, e = cryptor.Decrypt(fileContent)
		if e != nil {
			return e
		}
	}

	return nil
}

func iDecryptFile(ctx context.Context, filename string) error {

	cryptoState := getCryptoState(ctx)

	algorithm, e := createCryptoAlgorithm(cryptoState.cryptoAlgorithm, cryptoState.cipherKey, cryptoState.randomIv)
	if e != nil {
		return e
	}

	cryptor := crypto.NewCryptor(algorithm)
	file, e := os.Open(cryptoState.cryptoFeaturePath + "/assets/" + filename)
	if e != nil {
		return e
	}

	_, e = cryptor.DecryptStream(bufio.NewReader(file))
	if e != nil {
		cryptoState.err = e
	}
	_ = file.Close()
	return nil
}

func createCryptoAlgorithm(name string, cipherKey string, randomIv bool) (crypto.CryptoAlgorithm, error) {
	if name == "acrh" {
		return crypto.NewAesCBCCryptoAlgorithm(cipherKey)
	} else if name == "legacy" {
		return crypto.NewLegacyCryptoAlgorithm(cipherKey, randomIv)
	} else {
		return nil, fmt.Errorf("unknown crypto algorithm %s", name)
	}
}

func iEncryptFileAs(ctx context.Context, filename string, encryptionType string) error {
	cryptoState := getCryptoState(ctx)
	algorithm, e := createCryptoAlgorithm(cryptoState.cryptoAlgorithm, cryptoState.cipherKey, cryptoState.randomIv)
	if e != nil {
		return e
	}

	cryptor := crypto.NewCryptor(algorithm)
	file, e := os.Open(cryptoState.cryptoFeaturePath + "/assets/" + filename)
	if e != nil {
		return e
	}
	if encryptionType == "stream" {
		cryptoState.resultReader, e = cryptor.EncryptStream(bufio.NewReader(file))
		if e != nil {
			return e
		}
	} else {
		content, e := io.ReadAll(file)
		if e != nil {
			return e
		}
		cryptoState.result, e = cryptor.Encrypt(content)
		if e != nil {
			return e
		}

	}
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
		if strings.Contains(cryptoState.err.Error(), "unknown cryptor error") {
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
