package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAesGcm_EncryptAndDecryptMatchStdlib(t *testing.T) {
	cases := []struct {
		name      string
		cipherKey string
		nonce     string
		plaintext string
	}{
		{
			name:      "hello",
			cipherKey: "enigma",
			nonce:     "0102030405060708090a0b0c",
			plaintext: "Hello world",
		},
		{
			name:      "unicode",
			cipherKey: "enigma",
			nonce:     "111111111111111111111111",
			plaintext: "héllo 你好",
		},
		{
			name:      "binary",
			cipherKey: "myCipherKey",
			nonce:     "000102030405060708090a0b",
			plaintext: "\x00\xff\x7f\x80",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			nonce := mustDecodeHex(t, tc.nonce)
			plaintext := []byte(tc.plaintext)
			metadata, ciphertext := mustSealStdlib(t, tc.cipherKey, nonce, plaintext)

			cryptor := mustGcm(t, tc.cipherKey)
			encrypted, e := cryptor.encryptWithNonce(plaintext, nonce)
			require.NoError(t, e)
			assert.Equal(t, metadata, encrypted.Metadata)
			assert.Equal(t, ciphertext, encrypted.Data)
			assert.Equal(t, byte(0x0c), mustHeader(t, encrypted.Metadata)[sizePosition])

			decrypted, e := cryptor.Decrypt(encrypted)
			require.NoError(t, e)
			assert.Equal(t, plaintext, decrypted)

			stream, e := cryptor.DecryptStream(&EncryptedStreamData{
				Metadata: metadata,
				Reader:   bytes.NewReader(ciphertext),
			})
			require.NoError(t, e)
			streamed, e := io.ReadAll(stream)
			require.NoError(t, e)
			assert.Equal(t, plaintext, streamed)
		})
	}
}

func TestAesGcm_KnownCiphertexts(t *testing.T) {
	// Full PNED payloads. Metadata is the 12-byte nonce. The body is
	// ciphertext followed by the 16-byte authentication tag.
	cases := []struct {
		name      string
		cipherKey string
		plaintext string
		payload   string
	}{
		{
			name:      "hello",
			cipherKey: "enigma",
			plaintext: "Hello world",
			payload:   "504e4544014143474d0c0102030405060708090a0b0cea1b5780788b593294f898dabd1d02d7219c77218ad262c4a9aecc",
		},
		{
			name:      "unicode",
			cipherKey: "enigma",
			plaintext: "héllo 你好",
			payload:   "504e4544014143474d0c111111111111111111111111d138fbfc4bd0c337e29d9f2c98a5dd072360f4433e2979f440fc7827b3",
		},
		{
			name:      "binary",
			cipherKey: "myCipherKey",
			plaintext: "\x00\xff\x7f\x80",
			payload:   "504e4544014143474d0c000102030405060708090a0bc411f0f18fefca7becf7235b4d09b970d4d56d19",
		},
		{
			name:      "json-emoji",
			cipherKey: "enigma",
			plaintext: `{"text":"Hello, 世界! 🌍","n":42}`,
			payload:   "504e4544014143474d0cffffffffffffffffffffffff0d54b73ce2509d4e7655043f5b103cb4f1e78738a73a2f56109813955359b4784774de1954f835d4d739805631c5cdcc3467607259",
		},
		{
			name:      "one-block",
			cipherKey: "enigma",
			plaintext: "0123456789abcdef",
			payload:   "504e4544014143474d0c8f5a2d1b6e9c4ef0a37d83b9fd24a488008bd084305aaee605b349028c21b92d48f4dd5c71e99dce881c2ec9",
		},
		{
			name:      "multi-block-json",
			cipherKey: "enigma",
			plaintext: `{"text":"Hello, 世界! 🌍","n":42,"list":[1,2,3],"nested":{"ok":true}}`,
			payload:   "504e4544014143474d0cffffffffffffffffffffffff0d54b73ce2509d4e7655043f5b103cb4f1e78738a73a2f56109813955359b4784774de19050b4f14f7a439a859ac6c327856fe1af874bbe888945c603004176ddab22578e407663346d307669e83fa4c9a8b13743a07a21b7f",
		},
		{
			name:      "unicode-key-multiblock",
			cipherKey: "ключ 🔑 key",
			plaintext: "Line one\n\tLine two with \"quotes\" & symbols <>/\\ — spans several AES blocks.",
			payload:   "504e4544014143474d0c000000000000000000000000326052d9ff1b757ed8810f5d6f87a7c021920eb09b82e47049efecf5094f6f5837691e9fb66539dfc47d434cc01b783527c4d62ae2d8163043f86dc7ac405b90971ceaa07134add2be50dfd60f8390572d9b05255f3097614e888fcca4",
		},
		{
			name:      "single-byte",
			cipherKey: "customKey",
			plaintext: "a",
			payload:   "504e4544014143474d0ca1a2a3a4a5a6a7a8a9aaabacfbeef1ccce8f5a09087348786b99b75076",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			payload := mustDecodeHex(t, tc.payload)
			id, encrypted, e := parseHeader(payload)
			require.NoError(t, e)
			require.Equal(t, aesGcmCryptorID, *id)
			assert.Len(t, encrypted.Metadata, gcmNonceSize)

			decrypted, e := mustGcm(t, tc.cipherKey).Decrypt(encrypted)
			require.NoError(t, e)
			assert.Equal(t, tc.plaintext, string(decrypted))

			for _, module := range mustBothModules(t, tc.cipherKey) {
				out, e := module.Decrypt(payload)
				require.NoError(t, e)
				assert.Equal(t, tc.plaintext, string(out))
			}
		})
	}
}

func TestAesGcm_RejectsInvalidPayloads(t *testing.T) {
	nonce := mustDecodeHex(t, "0102030405060708090a0b0c")
	plaintext := []byte("Hello world")
	metadata, ciphertext := mustSealStdlib(t, "enigma", nonce, plaintext)
	cryptor := mustGcm(t, "enigma")

	tampered := append([]byte{}, ciphertext...)
	tampered[len(tampered)-1] ^= 0xff
	_, e := cryptor.Decrypt(&EncryptedData{Metadata: metadata, Data: tampered})
	assert.EqualError(t, e, "decryption error")

	other := mustGcm(t, "other-key")
	_, e = other.Decrypt(&EncryptedData{Metadata: metadata, Data: ciphertext})
	assert.EqualError(t, e, "decryption error")

	_, e = cryptor.Decrypt(&EncryptedData{Metadata: metadata[:gcmNonceSize-1], Data: ciphertext})
	assert.EqualError(t, e, "decryption error")

	_, e = cryptor.Decrypt(&EncryptedData{Metadata: metadata, Data: ciphertext[:gcmTagSize-1]})
	assert.EqualError(t, e, "decryption error")

	_, e = cryptor.Decrypt(nil)
	assert.EqualError(t, e, "decryption error")

	_, e = cryptor.Encrypt(nil)
	assert.EqualError(t, e, errAesGcmEncryptionUnavailable.Error())

	_, e = cryptor.Encrypt([]byte("Hello world"))
	assert.EqualError(t, e, errAesGcmEncryptionUnavailable.Error())

	_, e = cryptor.EncryptStream(bytes.NewReader([]byte("Hello world")))
	assert.EqualError(t, e, errAesGcmEncryptionUnavailable.Error())
}

func TestAesGcm_EmptyPlaintextDecrypts(t *testing.T) {
	nonce := mustDecodeHex(t, "090807060504030201000102")
	metadata, ciphertext := mustSealStdlib(t, "enigma", nonce, nil)
	assert.Equal(t, nonce, metadata)
	assert.Len(t, ciphertext, gcmTagSize)

	cryptor := mustGcm(t, "enigma")
	decrypted, e := cryptor.Decrypt(&EncryptedData{Metadata: metadata, Data: ciphertext})
	require.NoError(t, e)
	assert.Empty(t, decrypted)

	payload := append(mustHeader(t, metadata), ciphertext...)
	for _, module := range mustBothModules(t, "enigma") {
		out, e := module.Decrypt(payload)
		require.NoError(t, e)
		assert.Empty(t, out)

		reader, e := module.DecryptStream(bytes.NewReader(payload))
		require.NoError(t, e)
		streamed, e := io.ReadAll(reader)
		require.NoError(t, e)
		assert.Empty(t, streamed)
	}
}

func TestAesGcm_FactoriesDecryptItAndDoNotEncryptWithIt(t *testing.T) {
	nonce := mustDecodeHex(t, "0102030405060708090a0b0c")
	plaintext := []byte("Hello world")
	metadata, ciphertext := mustSealStdlib(t, "enigma", nonce, plaintext)
	payload := append(mustHeader(t, metadata), ciphertext...)

	cbc, e := NewAesCbcCryptoModule("enigma", true)
	require.NoError(t, e)
	legacy, e := NewLegacyCryptoModule("enigma", false)
	require.NoError(t, e)

	for _, module := range []CryptoModule{cbc, legacy} {
		out, e := module.Decrypt(payload)
		require.NoError(t, e)
		assert.Equal(t, plaintext, out)

		reader, e := module.DecryptStream(bytes.NewReader(payload))
		require.NoError(t, e)
		streamed, e := io.ReadAll(reader)
		require.NoError(t, e)
		assert.Equal(t, plaintext, streamed)
	}

	cbcEncrypted, e := cbc.Encrypt(plaintext)
	require.NoError(t, e)
	cbcID, e := peekHeaderCryptorId(cbcEncrypted)
	require.NoError(t, e)
	assert.Equal(t, crivId, *cbcID)

	legacyEncrypted, e := legacy.Encrypt(plaintext)
	require.NoError(t, e)
	legacyID, e := peekHeaderCryptorId(legacyEncrypted)
	require.NoError(t, e)
	assert.Equal(t, legacyId, *legacyID)

	roundTrip, e := cbc.Decrypt(cbcEncrypted)
	require.NoError(t, e)
	assert.Equal(t, plaintext, roundTrip)
	roundTrip, e = legacy.Decrypt(legacyEncrypted)
	require.NoError(t, e)
	assert.Equal(t, plaintext, roundTrip)
}

func TestAesGcm_EncryptIsDisabled(t *testing.T) {
	cryptor := mustGcm(t, "enigma")
	_, e := cryptor.Encrypt([]byte("Hello world"))
	assert.EqualError(t, e, errAesGcmEncryptionUnavailable.Error())

	_, e = cryptor.EncryptStream(bytes.NewReader([]byte("Hello world")))
	assert.EqualError(t, e, errAesGcmEncryptionUnavailable.Error())
}

func mustGcm(t *testing.T, cipherKey string) *aesGcmCryptor {
	t.Helper()
	cryptor, e := newAesGcmCryptor(cipherKey)
	require.NoError(t, e)
	return cryptor.(*aesGcmCryptor)
}

func mustBothModules(t *testing.T, cipherKey string) []CryptoModule {
	t.Helper()
	cbc, e := NewAesCbcCryptoModule(cipherKey, true)
	require.NoError(t, e)
	legacy, e := NewLegacyCryptoModule(cipherKey, true)
	require.NoError(t, e)
	return []CryptoModule{cbc, legacy}
}

func mustHeader(t *testing.T, metadata []byte) []byte {
	t.Helper()
	header, e := headerV1(aesGcmCryptorID, metadata)
	require.NoError(t, e)
	return header
}

func mustDecodeHex(t *testing.T, value string) []byte {
	t.Helper()
	decoded, e := hex.DecodeString(value)
	require.NoError(t, e)
	return decoded
}

func mustSealStdlib(t *testing.T, cipherKey string, nonce, plaintext []byte) (metadata, ciphertext []byte) {
	t.Helper()
	sum := sha256.Sum256([]byte(cipherKey))
	block, e := aes.NewCipher(sum[:])
	require.NoError(t, e)
	aead, e := cipher.NewGCM(block)
	require.NoError(t, e)
	sealed := aead.Seal(nil, nonce, plaintext, nil)
	metadata = append([]byte{}, nonce...)
	ciphertext = sealed
	return metadata, ciphertext
}
