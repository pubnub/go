// snippet.includes
// Replace with your package name (usually "main")
package pubnub_samples_test

import (
	"fmt"

	pubnub "github.com/pubnub/go/v7"
)

// snippet.end

/*
IMPORTANT NOTE FOR COPYING EXAMPLES:

Throughout this file, you'll see code between "snippet.hide" and "snippet.show" comments.
These sections are used for CI/CD testing and should be SKIPPED if you're copying examples.

Example of what to skip:
	// snippet.hide
	config = setPubnubExamplePAMConfigData(config)  // <- Skip this line (for testing only)
	// snippet.show

When copying examples to your own code:
- Use your own publish/subscribe/secret keys instead of the "demo" keys
- Remove any statements that are between snippet.hide and snippet.show (they're only for testing purposes)
- Ensure Access Manager is enabled in your PubNub Admin Portal
*/

// snippet.grant_channel
// Example_grantChannelPermissions demonstrates granting read/write permissions on a channel
func Example_grantChannelPermissions() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for PAM)

	// snippet.hide
	config = setPubnubExamplePAMConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Grant read and write permissions on a channel
	_, status, err := pn.Grant().
		Channels([]string{"my-channel"}).
		Read(true).  // Allow reading messages
		Write(true). // Allow publishing messages
		TTL(60).     // Permissions valid for 60 minutes
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("Channel permissions granted successfully")
	}

	// Output:
	// Channel permissions granted successfully
}

// snippet.grant_auth_key
// Example_grantAuthKeyPermissions demonstrates granting permissions to a specific auth key
func Example_grantAuthKeyPermissions() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.SecretKey = "demo-secret"

	// snippet.hide
	config = setPubnubExamplePAMConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Grant permissions to a specific auth key on a channel
	_, status, err := pn.Grant().
		Channels([]string{"my-channel"}).
		AuthKeys([]string{"my-auth-key"}). // Specific user/client auth key
		Read(true).
		Write(false). // Only allow reading, not writing
		TTL(1440).    // Valid for 24 hours (1440 minutes)
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("Auth key permissions granted successfully")
	}

	// Output:
	// Auth key permissions granted successfully
}

// snippet.grant_channel_group
// Example_grantChannelGroupPermissions demonstrates granting permissions on a channel group
func Example_grantChannelGroupPermissions() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.SecretKey = "demo-secret"

	// snippet.hide
	config = setPubnubExamplePAMConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Grant read and manage permissions on a channel group
	_, status, err := pn.Grant().
		ChannelGroups([]string{"my-channel-group"}).
		Read(true).   // Allow reading from channels in the group
		Manage(true). // Allow managing the channel group
		TTL(0).       // TTL=0 means permissions never expire
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("Channel group permissions granted successfully")
	}

	// Output:
	// Channel group permissions granted successfully
}

// snippet.grant_token_channel
// Example_grantTokenChannel demonstrates creating an access token for a channel (PAM v3)
func Example_grantTokenChannel() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.SecretKey = "demo-secret"

	// snippet.hide
	config = setPubnubExamplePAMConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Define channel permissions
	channelPermissions := map[string]pubnub.ChannelPermissions{
		"my-channel": {
			Read:  true,
			Write: true,
		},
	}

	// Grant token with channel permissions
	response, status, err := pn.GrantToken().
		TTL(60). // Token valid for 60 minutes
		Channels(channelPermissions).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 && response.Data.Token != "" {
		fmt.Println("Access token granted successfully")
	}

	// Output:
	// Access token granted successfully
}

// snippet.grant_token_authorized_user
// Example_grantTokenAuthorizedUser demonstrates creating a token for a specific user
func Example_grantTokenAuthorizedUser() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.SecretKey = "demo-secret"

	// snippet.hide
	config = setPubnubExamplePAMConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Define channel permissions
	channelPermissions := map[string]pubnub.ChannelPermissions{
		"user-channel": {
			Read:  true,
			Write: true,
		},
	}

	// Grant token authorized for a specific user UUID
	response, status, err := pn.GrantToken().
		TTL(1440). // Token valid for 24 hours
		AuthorizedUUID("user-123").
		Channels(channelPermissions).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 && response.Data.Token != "" {
		fmt.Println("User-authorized token granted successfully")
	}

	// Output:
	// User-authorized token granted successfully
}

// snippet.grant_token_pattern
// Example_grantTokenPattern demonstrates creating a token with pattern-based permissions
func Example_grantTokenPattern() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.SecretKey = "demo-secret"

	// snippet.hide
	config = setPubnubExamplePAMConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Define pattern-based permissions (matches multiple channels)
	channelPatterns := map[string]pubnub.ChannelPermissions{
		"room.*": { // Matches room.1, room.2, room.lobby, etc.
			Read:  true,
			Write: true,
		},
	}

	// Grant token with pattern-based permissions
	response, status, err := pn.GrantToken().
		TTL(60).
		ChannelsPattern(channelPatterns).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 && response.Data.Token != "" {
		fmt.Println("Pattern-based token granted successfully")
	}

	// Output:
	// Pattern-based token granted successfully
}

// snippet.grant_token_multi_resource
// Example_grantTokenMultiResource demonstrates creating a token with multiple resource types
func Example_grantTokenMultiResource() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.SecretKey = "demo-secret"

	// snippet.hide
	config = setPubnubExamplePAMConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Define permissions for channels
	channelPermissions := map[string]pubnub.ChannelPermissions{
		"channel-1": {
			Read:  true,
			Write: true,
		},
	}

	// Define permissions for channel groups
	groupPermissions := map[string]pubnub.GroupPermissions{
		"group-1": {
			Read: true,
		},
	}

	// Define permissions for UUIDs (for user metadata operations)
	uuidPermissions := map[string]pubnub.UUIDPermissions{
		"user-123": {
			Get:    true,
			Update: true,
		},
	}

	// Grant token with multiple resource types
	response, status, err := pn.GrantToken().
		TTL(60).
		Channels(channelPermissions).
		ChannelGroups(groupPermissions).
		UUIDs(uuidPermissions).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 && response.Data.Token != "" {
		fmt.Println("Multi-resource token granted successfully")
	}

	// Output:
	// Multi-resource token granted successfully
}

// snippet.parse_token
// Example_parseToken demonstrates parsing a token to view its permissions
func Example_parseToken() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.SecretKey = "demo-secret"

	// snippet.hide
	config = setPubnubExamplePAMConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	// First, grant a token to parse
	channelPermissions := map[string]pubnub.ChannelPermissions{
		"my-channel": {
			Read:  true,
			Write: true,
		},
	}
	grantResp, _, _ := pn.GrantToken().
		TTL(60).
		Channels(channelPermissions).
		Execute()
	token := grantResp.Data.Token
	// snippet.show

	// Parse the token to view its permissions
	parsedToken, err := pubnub.ParseToken(token)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Display token information
	fmt.Printf("Token TTL: %d minutes\n", parsedToken.TTL)

	// Show channel permissions
	for channelName, perms := range parsedToken.Resources.Channels {
		fmt.Printf("Channel '%s' - Read: %v, Write: %v\n",
			channelName, perms.Read, perms.Write)
	}

	// Output:
	// Token TTL: 60 minutes
	// Channel 'my-channel' - Read: true, Write: true
}

// snippet.set_token
// Example_setToken demonstrates setting an auth token on the client
func Example_setToken() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExamplePAMConfigData(config)

	// Create a token first for demonstration
	configWithSecret := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	configWithSecret = setPubnubExamplePAMConfigData(configWithSecret)
	pnAdmin := pubnub.NewPubNub(configWithSecret)

	channelPermissions := map[string]pubnub.ChannelPermissions{
		"my-channel": {
			Read:  true,
			Write: true,
		},
	}
	grantResp, _, _ := pnAdmin.GrantToken().
		TTL(60).
		AuthorizedUUID("demo-user").
		Channels(channelPermissions).
		Execute()
	token := grantResp.Data.Token
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// Set the auth token on the client
	pn.SetToken(token)

	// Now the client can access authorized resources
	response, status, err := pn.Publish().
		Channel("my-channel").
		Message("Hello with token!").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 && response.Timestamp > 0 {
		fmt.Println("Message published with token successfully")
	}

	// Output:
	// Message published with token successfully
}

// snippet.revoke_token
// Example_revokeToken demonstrates revoking an access token
func Example_revokeToken() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"
	config.SecretKey = "demo-secret"

	// snippet.hide
	config = setPubnubExamplePAMConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	// First, grant a token to revoke
	channelPermissions := map[string]pubnub.ChannelPermissions{
		"my-channel": {
			Read:  true,
			Write: true,
		},
	}
	grantResp, _, _ := pn.GrantToken().
		TTL(60).
		Channels(channelPermissions).
		Execute()
	token := grantResp.Data.Token
	// snippet.show

	// Revoke the token
	_, status, err := pn.RevokeToken().
		Token(token).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("Token revoked successfully")
	}

	// Output:
	// Token revoked successfully
}

// snippet.end
