package utils

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"testing/quick"
)

func doesEncryptItWork(in []byte) bool {
	cryptor, e := NewCBCCryptor("enigma")
	if e != nil {
		return false
	}

	outputBuff := bytes.NewBuffer(make([]byte, 0, len(in)))
	metadata, err := cryptor.EncryptStream(bytes.NewReader(in), outputBuff)
	if err != nil {
		println(err.Error())
		return false
	}

	decrypted, err := cryptor.Decrypt(EncryptedData{
		Data:     outputBuff.Bytes(),
		Metadata: metadata,
	})
	if err != nil {
		println(err.Error())
		return false
	}
	return bytes.Equal(in, decrypted)
}

func doesDecryptItWork(in []byte) bool {
	cryptor, e := NewCBCCryptor("enigma")
	if e != nil {
		return false
	}

	outputBuff := bytes.NewBuffer(make([]byte, 0, len(in)))
	metadata, err := cryptor.EncryptStream(bytes.NewReader(in), outputBuff)
	if err != nil {
		println(err.Error())
		return false
	}

	decrypted, err := cryptor.Decrypt(EncryptedData{
		Data:     outputBuff.Bytes(),
		Metadata: metadata,
	})
	if err != nil {
		println(err.Error())
		return false
	}
	return bytes.Equal(in, decrypted)
}

func TestEncryptStreamWorks(t *testing.T) {
	if err := quick.Check(doesEncryptItWork, nil); err != nil {
		t.Error(err)
	}
}

func TestDecryptStreamWorks(t *testing.T) {
	if err := quick.Check(doesDecryptItWork, nil); err != nil {
		t.Error(err)
	}
}

func TestTest(t *testing.T) {
	assert.True(t, doesEncryptItWork([]byte{0x96, 0x8f, 0xc6, 0x16, 0x7, 0x60, 0x8, 0xd2, 0xcc, 0x3e, 0xe1, 0x1a, 0x40, 0xbe, 0x75, 0xa9, 0x91, 0x82, 0x40, 0x5d, 0x92, 0xa2, 0x68, 0xd4, 0x68, 0x7, 0xfd, 0x43, 0x7, 0x9b, 0x2d, 0x77, 0x2e, 0xaa, 0xc7, 0x18, 0xe1, 0xb7}))
}

func TestCBCCryptor_EncryptStream(t *testing.T) {
	cryptor, e := NewCBCCryptor("enigma")
	if e != nil {
		t.Errorf("error, %s", e)
	}

	message := "sure i can"
	buff := bytes.NewBuffer(make([]byte, 0, len(message)))
	metadata, err := cryptor.EncryptStream(strings.NewReader(message), buff)
	if err != nil {
		t.Errorf("error, %s", err)
	}

	decrypted, err := cryptor.Decrypt(EncryptedData{Data: buff.Bytes(), Metadata: metadata})
	if err != nil {
		t.Errorf("error, %s", err)
	}

	a := assert.New(t)
	a.Equal(message, string(decrypted))
}

func TestCBCCryptor_DecryptStream(t *testing.T) {
	cryptor, e := NewCBCCryptor("enigma")
	if e != nil {
		t.Errorf("error, %s", e)
	}

	message := "sure i can"
	buff := bytes.NewBuffer(make([]byte, 0, len(message)))
	metadata, err := cryptor.EncryptStream(strings.NewReader(message), buff)
	if err != nil {
		t.Errorf("error, %s", err)
	}

	decryptBuff := bytes.NewBuffer(make([]byte, 0, len(message)))

	err2 := cryptor.DecryptStream(bytes.NewReader(buff.Bytes()), metadata, decryptBuff)
	if err2 != nil {
		t.Errorf("error, %s", err2)
	}

	a := assert.New(t)
	a.Equal(message, decryptBuff.String())
}
