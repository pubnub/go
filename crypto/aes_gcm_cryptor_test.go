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
			assert.Equal(t, byte(0x1c), mustHeader(t, encrypted.Metadata)[sizePosition])

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
	// Full PNED payloads. Metadata is the 12-byte nonce followed by the
	// 16-byte tag; the remaining bytes are ciphertext.
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
			payload:   "504e4544014143474d1c0102030405060708090a0b0cdabd1d02d7219c77218ad262c4a9aeccea1b5780788b593294f898",
		},
		{
			name:      "unicode",
			cipherKey: "enigma",
			plaintext: "héllo 你好",
			payload:   "504e4544014143474d1c111111111111111111111111a5dd072360f4433e2979f440fc7827b3d138fbfc4bd0c337e29d9f2c98",
		},
		{
			name:      "binary",
			cipherKey: "myCipherKey",
			plaintext: "\x00\xff\x7f\x80",
			payload:   "504e4544014143474d1c000102030405060708090a0b8fefca7becf7235b4d09b970d4d56d19c411f0f1",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			payload := mustDecodeHex(t, tc.payload)
			id, encrypted, e := parseHeader(payload)
			require.NoError(t, e)
			require.Equal(t, aesGcmCryptorID, *id)
			assert.Len(t, encrypted.Metadata, gcmMetadataSize)

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

	tampered := append([]byte{}, metadata...)
	tampered[gcmNonceSize] ^= 0xff
	_, e := cryptor.Decrypt(&EncryptedData{Metadata: tampered, Data: ciphertext})
	assert.EqualError(t, e, "decryption error")

	other := mustGcm(t, "other-key")
	_, e = other.Decrypt(&EncryptedData{Metadata: metadata, Data: ciphertext})
	assert.EqualError(t, e, "decryption error")

	_, e = cryptor.Decrypt(&EncryptedData{Metadata: metadata[:gcmMetadataSize-1], Data: ciphertext})
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

func TestAesGcm_EmptyPlaintextDecryptsWhenTagIsInMetadata(t *testing.T) {
	nonce := mustDecodeHex(t, "090807060504030201000102")
	metadata, ciphertext := mustSealStdlib(t, "enigma", nonce, nil)
	assert.Empty(t, ciphertext)
	assert.Len(t, metadata, gcmMetadataSize)

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
	ciphertext = append([]byte{}, sealed[:len(sealed)-gcmTagSize]...)
	tag := sealed[len(sealed)-gcmTagSize:]
	metadata = make([]byte, 0, gcmMetadataSize)
	metadata = append(metadata, nonce...)
	metadata = append(metadata, tag...)
	return metadata, ciphertext
}
