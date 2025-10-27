// snippet.includes
// Replace with your package name (usually "main")
package pubnub_samples_test

import (
	"fmt"
	"time"

	pubnub "github.com/pubnub/go/v7"
	"github.com/pubnub/go/v7/crypto"
)

// snippet.end

/*
IMPORTANT NOTE FOR COPYING EXAMPLES:

Throughout this file, you'll see code between "snippet.hide" and "snippet.show" comments.
These sections are used for CI/CD testing and should be SKIPPED if you're copying examples.

Example of what to skip:
	// snippet.hide
	config = setPubnubExampleConfigData(config)  // <- Skip this line (for testing only)
	// snippet.show

When copying examples to your own code:
- Use your own publish/subscribe keys instead of the "demo" keys
- Remove any statements that are between snippet.hide and snippet.show (they're only for testing purposes)
*/

// snippet.publish_encrypted_message
// Example_publishEncryptedMessage demonstrates publishing an encrypted message using CryptoModule
func Example_publishEncryptedMessage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo" // Replace with your subscribe key
	config.PublishKey = "demo"   // Replace with your publish key

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	// Create a crypto module for encryption
	cryptoModule, _ := crypto.NewAesCbcCryptoModule("my-secret-key", true)
	config.CryptoModule = cryptoModule

	pn := pubnub.NewPubNub(config)

	// Publish an encrypted message
	response, status, err := pn.Publish().
		Channel("encrypted-channel").
		Message("This is a secret message!").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 && response.Timestamp > 0 {
		fmt.Println("Encrypted message published successfully")
	}

	// Output:
	// Encrypted message published successfully
}

// snippet.subscribe_encrypted_message
// Example_subscribeEncryptedMessage demonstrates receiving and decrypting messages
func Example_subscribeEncryptedMessage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	// Create a crypto module with the same key used for encryption
	cryptoModule, _ := crypto.NewAesCbcCryptoModule("my-secret-key", true)
	config.CryptoModule = cryptoModule

	pn := pubnub.NewPubNub(config)

	// Create listener to receive messages
	listener := pubnub.NewListener()

	// Create a done channel to stop the goroutine when needed
	done := make(chan bool)

	// snippet.hide
	messageReceived := make(chan bool)
	// snippet.show

	go func() {
		for {
			select {
			case status := <-listener.Status:
				// snippet.hide
				if status.Category == pubnub.PNConnectedCategory {
					// Once connected, publish the encrypted message
					pn.Publish().
						Channel("encrypted-sub-channel").
						Message("Secret data").
						Execute()
				}
				// snippet.show

			case message := <-listener.Message:
				// Message is automatically decrypted
				fmt.Printf("Received decrypted message: %v\n", message.Message)
				// snippet.hide
				messageReceived <- true
				// snippet.show
				return

			case <-done:
				return
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().
		Channels([]string{"encrypted-sub-channel"}).
		Execute()

	fmt.Println("Subscribed to encrypted channel")

	// snippet.hide
	// Wait for message with timeout
	select {
	case <-messageReceived:
		// Message received successfully
	case <-time.After(15 * time.Second):
		fmt.Println("Timeout - prevent hanging")
		// Timeout - prevent hanging
	}
	// snippet.show

	// When done, unsubscribe and stop goroutine
	pn.UnsubscribeAll()
	close(done)

	// Output:
	// Subscribed to encrypted channel
	// Received decrypted message: Secret data
}

// snippet.encrypt_decrypt_string
// Example_encryptDecryptString demonstrates encrypting and decrypting a string directly
func Example_encryptDecryptString() {
	// Create a crypto module
	cryptoModule, err := crypto.NewAesCbcCryptoModule("my-cipher-key", true)
	if err != nil {
		fmt.Printf("Error creating crypto module: %v\n", err)
		return
	}

	// Original message
	originalMessage := "This is a sensitive message"

	// Encrypt the message
	encryptedData, err := cryptoModule.Encrypt([]byte(originalMessage))
	if err != nil {
		fmt.Printf("Error encrypting: %v\n", err)
		return
	}

	fmt.Println("Message encrypted successfully")

	// Decrypt the message
	decryptedData, err := cryptoModule.Decrypt(encryptedData)
	if err != nil {
		fmt.Printf("Error decrypting: %v\n", err)
		return
	}

	fmt.Printf("Decrypted message: %s\n", string(decryptedData))

	// Output:
	// Message encrypted successfully
	// Decrypted message: This is a sensitive message
}

// snippet.legacy_cipher_key
// Example_legacyCipherKey demonstrates using the legacy CipherKey configuration (deprecated)
func Example_legacyCipherKey() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	// Using CipherKey (deprecated - use CryptoModule instead)
	config.CipherKey = "my-cipher-key"

	pn := pubnub.NewPubNub(config)

	// Messages will be automatically encrypted/decrypted
	response, status, err := pn.Publish().
		Channel("legacy-encrypted-channel").
		Message("Legacy encrypted message").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 && response.Timestamp > 0 {
		fmt.Println("Message published with legacy encryption")
	}

	// Output:
	// Message published with legacy encryption
}

// snippet.publish_encrypted_json
// Example_publishEncryptedJSON demonstrates encrypting JSON data
func Example_publishEncryptedJSON() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	// Create crypto module
	cryptoModule, _ := crypto.NewAesCbcCryptoModule("my-secret-key", true)
	config.CryptoModule = cryptoModule

	pn := pubnub.NewPubNub(config)

	// Publish encrypted JSON data
	userData := map[string]interface{}{
		"userId":   "user123",
		"email":    "user@example.com",
		"balance":  1000.50,
		"verified": true,
	}

	response, status, err := pn.Publish().
		Channel("encrypted-json-channel").
		Message(userData).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 && response.Timestamp > 0 {
		fmt.Println("Encrypted JSON published successfully")
	}

	// Output:
	// Encrypted JSON published successfully
}

// snippet.multiple_crypto_modules
// Example_multipleCryptoModules demonstrates using different encryption for different channels
func Example_multipleCryptoModules() {
	// Create first PubNub instance with one encryption key
	config1 := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user-1"))
	config1.SubscribeKey = "demo"
	config1.PublishKey = "demo"

	// snippet.hide
	config1 = setPubnubExampleConfigData(config1)
	// snippet.show

	cryptoModule1, _ := crypto.NewAesCbcCryptoModule("key-for-channel-1", true)
	config1.CryptoModule = cryptoModule1

	pn1 := pubnub.NewPubNub(config1)

	// Create second PubNub instance with different encryption key
	config2 := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user-2"))
	config2.SubscribeKey = "demo"
	config2.PublishKey = "demo"

	// snippet.hide
	config2 = setPubnubExampleConfigData(config2)
	// snippet.show

	cryptoModule2, _ := crypto.NewAesCbcCryptoModule("key-for-channel-2", true)
	config2.CryptoModule = cryptoModule2

	pn2 := pubnub.NewPubNub(config2)

	// Publish to different channels with different encryption
	response1, status1, err1 := pn1.Publish().
		Channel("secure-channel-1").
		Message("Message for channel 1").
		Execute()

	response2, status2, err2 := pn2.Publish().
		Channel("secure-channel-2").
		Message("Message for channel 2").
		Execute()

	if err1 == nil && err2 == nil &&
		status1.StatusCode == 200 && status2.StatusCode == 200 &&
		response1.Timestamp > 0 && response2.Timestamp > 0 {
		fmt.Println("Messages published with different encryption keys")
	}

	// Output:
	// Messages published with different encryption keys
}

// snippet.history_with_encryption
// Example_historyWithEncryption demonstrates fetching encrypted message history
func Example_historyWithEncryption() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	// Create crypto module
	cryptoModule, _ := crypto.NewAesCbcCryptoModule("my-history-key", true)
	config.CryptoModule = cryptoModule

	pn := pubnub.NewPubNub(config)

	// Publish an encrypted message
	pn.Publish().
		Channel("history-encrypted-channel").
		Message("Encrypted historical message").
		Execute()

	// Fetch history - messages will be automatically decrypted
	response, status, err := pn.History().
		Channel("history-encrypted-channel").
		Count(10).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 && len(response.Messages) > 0 {
		fmt.Println("Fetched encrypted message history")
		// Messages are automatically decrypted
		fmt.Printf("First message: %v\n", response.Messages[0].Message)
	}

	// Output:
	// Fetched encrypted message history
	// First message: Encrypted historical message
}

// snippet.end
