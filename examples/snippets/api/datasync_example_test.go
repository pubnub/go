// snippet.includes
// Replace with your package name (usually "main")
package pubnub_samples_test

import (
	"fmt"
	"time"

	pubnub "github.com/pubnub/go/v9"
)

// snippet.end

/*
IMPORTANT NOTE FOR COPYING EXAMPLES:

Throughout this file, you'll see code between "snippet.hide" and "snippet.show" comments.
These sections are used for CI/CD testing and should be SKIPPED if you're copying examples.

When copying examples to your own code:
- Use your own publish/subscribe/secret keys instead of the "demo" keys
- Ensure Access Manager is enabled in your PubNub Admin Portal, DataSync requires it on every request
*/

// ==================== Users ====================

// snippet.create_user_basic_usage
// Example_createUserBasicUsage demonstrates creating a DataSync user
func Example_createUserBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("user-alice").Execute()
	// snippet.show

	_, status, err := pn.DataSync.CreateUser().
		ID("user-alice").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Alice", "type": "shopper"}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("User created, status: %d\n", status.StatusCode)

	// Output:
	// User created, status: 200
}

// snippet.get_user_basic_usage
// Example_getUserBasicUsage demonstrates reading a single DataSync user
func Example_getUserBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("user-alice").Execute()
	pn.DataSync.CreateUser().
		ID("user-alice").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Alice", "type": "shopper"}).
		Execute()
	// snippet.show

	res, status, err := pn.DataSync.GetUser().
		ID("user-alice").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Got user: %s, status: %d\n", res.Data.ID, status.StatusCode)

	// Output:
	// Got user: user-alice, status: 200
}

// snippet.get_users_basic_usage
// Example_getUsersBasicUsage demonstrates listing DataSync users
func Example_getUsersBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("user-alice").Execute()
	pn.DataSync.CreateUser().
		ID("user-alice").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Alice", "type": "shopper"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.GetUsers().
		Limit(20).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Users listed, status: %d\n", status.StatusCode)

	// Output:
	// Users listed, status: 200
}

// snippet.get_users_filter_fast
// Example_getUsersFilterFast demonstrates filtering users with FilterFast
func Example_getUsersFilterFast() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("user-alice").Execute()
	pn.DataSync.CreateUser().
		ID("user-alice").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Alice", "type": "shopper"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.GetUsers().
		FilterFast("type == \"shopper\"").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Filtered users retrieved, status: %d\n", status.StatusCode)

	// Output:
	// Filtered users retrieved, status: 200
}

// snippet.get_users_filter
// Example_getUsersFilter demonstrates filtering users with Filter
func Example_getUsersFilter() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("user-alice").Execute()
	pn.DataSync.CreateUser().
		ID("user-alice").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Alice", "type": "shopper"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.GetUsers().
		Filter("name LIKE \"*Alice*\"").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Filtered users retrieved, status: %d\n", status.StatusCode)

	// Output:
	// Filtered users retrieved, status: 200
}

// snippet.get_users_pagination
// Example_getUsersPagination demonstrates paging through the user list
func Example_getUsersPagination() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("user-alice").Execute()
	pn.DataSync.CreateUser().
		ID("user-alice").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Alice", "type": "shopper"}).
		Execute()
	// snippet.show

	res, status, err := pn.DataSync.GetUsers().
		Limit(20).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Users page retrieved, status: %d\n", status.StatusCode)

	if res.Meta != nil && res.Meta.HasNext {
		_, _, err = pn.DataSync.GetUsers().
			Cursor(res.Meta.NextCursor).
			Limit(20).
			Execute()

		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}

	// Output:
	// Users page retrieved, status: 200
}

// snippet.set_user_basic_usage
// Example_setUserBasicUsage demonstrates replacing a DataSync user in full
func Example_setUserBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("user-alice").Execute()
	pn.DataSync.CreateUser().
		ID("user-alice").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Alice", "type": "shopper"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.SetUser().
		ID("user-alice").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Alice B.", "type": "shopper"}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("User replaced, status: %d\n", status.StatusCode)

	// Output:
	// User replaced, status: 200
}

// snippet.update_user_basic_usage
// Example_updateUserBasicUsage demonstrates applying a JSON Patch update to a DataSync user
func Example_updateUserBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("user-alice").Execute()
	pn.DataSync.CreateUser().
		ID("user-alice").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Alice", "type": "shopper"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.UpdateUser().
		ID("user-alice").
		Replace("/status", "active").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("User updated, status: %d\n", status.StatusCode)

	// Output:
	// User updated, status: 200
}

// snippet.remove_user_basic_usage
// Example_removeUserBasicUsage demonstrates deleting a DataSync user
func Example_removeUserBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("user-alice").Execute()
	pn.DataSync.CreateUser().
		ID("user-alice").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Alice", "type": "shopper"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.RemoveUser().
		ID("user-alice").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("User removed, status: %d\n", status.StatusCode)

	// Output:
	// User removed, status: 200
}

// ==================== Channels ====================

// snippet.create_channel_basic_usage
// Example_createChannelBasicUsage demonstrates creating a DataSync channel
func Example_createChannelBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveChannel().ID("channel-summer-sale").Execute()
	// snippet.show

	_, status, err := pn.DataSync.CreateChannel().
		ID("channel-summer-sale").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Summer Sale", "type": "promotion"}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Channel created, status: %d\n", status.StatusCode)

	// Output:
	// Channel created, status: 200
}

// snippet.get_channel_basic_usage
// Example_getChannelBasicUsage demonstrates reading a single DataSync channel
func Example_getChannelBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveChannel().ID("channel-summer-sale").Execute()
	pn.DataSync.CreateChannel().
		ID("channel-summer-sale").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Summer Sale", "type": "promotion"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.GetChannel().
		ID("channel-summer-sale").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Channel retrieved, status: %d\n", status.StatusCode)

	// Output:
	// Channel retrieved, status: 200
}

// snippet.get_channels_basic_usage
// Example_getChannelsBasicUsage demonstrates listing DataSync channels
func Example_getChannelsBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveChannel().ID("channel-summer-sale").Execute()
	pn.DataSync.CreateChannel().
		ID("channel-summer-sale").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Summer Sale", "type": "promotion"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.GetChannels().
		Limit(20).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Channels listed, status: %d\n", status.StatusCode)

	// Output:
	// Channels listed, status: 200
}

// snippet.get_channels_filter_fast
// Example_getChannelsFilterFast demonstrates filtering channels with FilterFast
func Example_getChannelsFilterFast() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveChannel().ID("channel-summer-sale").Execute()
	pn.DataSync.CreateChannel().
		ID("channel-summer-sale").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Summer Sale", "type": "promotion"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.GetChannels().
		FilterFast("type == \"promotion\"").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Filtered channels retrieved, status: %d\n", status.StatusCode)

	// Output:
	// Filtered channels retrieved, status: 200
}

// snippet.get_channels_filter
// Example_getChannelsFilter demonstrates filtering channels with Filter
func Example_getChannelsFilter() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveChannel().ID("channel-summer-sale").Execute()
	pn.DataSync.CreateChannel().
		ID("channel-summer-sale").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Summer Sale", "type": "promotion"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.GetChannels().
		Filter("name LIKE \"*Sale*\"").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Filtered channels retrieved, status: %d\n", status.StatusCode)

	// Output:
	// Filtered channels retrieved, status: 200
}

// snippet.get_channels_pagination
// Example_getChannelsPagination demonstrates paging through the channel list
func Example_getChannelsPagination() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveChannel().ID("channel-summer-sale").Execute()
	pn.DataSync.CreateChannel().
		ID("channel-summer-sale").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Summer Sale", "type": "promotion"}).
		Execute()
	// snippet.show

	res, status, err := pn.DataSync.GetChannels().
		Limit(20).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Channels page retrieved, status: %d\n", status.StatusCode)

	if res.Meta != nil && res.Meta.HasNext {
		_, _, err = pn.DataSync.GetChannels().
			Cursor(res.Meta.NextCursor).
			Limit(20).
			Execute()

		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}

	// Output:
	// Channels page retrieved, status: 200
}

// snippet.set_channel_basic_usage
// Example_setChannelBasicUsage demonstrates replacing a DataSync channel in full
func Example_setChannelBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveChannel().ID("channel-summer-sale").Execute()
	pn.DataSync.CreateChannel().
		ID("channel-summer-sale").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Summer Sale", "type": "promotion"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.SetChannel().
		ID("channel-summer-sale").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Summer Sale 2026", "type": "promotion"}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Channel replaced, status: %d\n", status.StatusCode)

	// Output:
	// Channel replaced, status: 200
}

// snippet.update_channel_basic_usage
// Example_updateChannelBasicUsage demonstrates applying a JSON Patch update to a DataSync channel
func Example_updateChannelBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveChannel().ID("channel-summer-sale").Execute()
	pn.DataSync.CreateChannel().
		ID("channel-summer-sale").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Summer Sale", "type": "promotion"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.UpdateChannel().
		ID("channel-summer-sale").
		Replace("/status", "active").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Channel updated, status: %d\n", status.StatusCode)

	// Output:
	// Channel updated, status: 200
}

// snippet.remove_channel_basic_usage
// Example_removeChannelBasicUsage demonstrates deleting a DataSync channel
func Example_removeChannelBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveChannel().ID("channel-summer-sale").Execute()
	pn.DataSync.CreateChannel().
		ID("channel-summer-sale").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Summer Sale", "type": "promotion"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.RemoveChannel().
		ID("channel-summer-sale").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Channel removed, status: %d\n", status.StatusCode)

	// Output:
	// Channel removed, status: 200
}

// ==================== Memberships ====================

// snippet.create_membership_basic_usage
// Example_createMembershipBasicUsage demonstrates linking a user to a channel
func Example_createMembershipBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("user-alice").Execute()
	pn.DataSync.CreateUser().
		ID("user-alice").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Alice", "type": "shopper"}).
		Execute()
	pn.DataSync.RemoveChannel().ID("channel-summer-sale").Execute()
	pn.DataSync.CreateChannel().
		ID("channel-summer-sale").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Summer Sale", "type": "promotion"}).
		Execute()
	pn.DataSync.RemoveMembership().ID("membership-alice-summer-sale").Execute()
	// snippet.show

	_, status, err := pn.DataSync.CreateMembership().
		ID("membership-alice-summer-sale").
		ChannelID("channel-summer-sale").
		UserID("user-alice").
		RelationshipClassVersion(1).
		Payload(map[string]interface{}{"role": "viewer"}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Membership created, status: %d\n", status.StatusCode)

	// Output:
	// Membership created, status: 200
}

// snippet.get_membership_basic_usage
// Example_getMembershipBasicUsage demonstrates reading a single membership
func Example_getMembershipBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("user-alice").Execute()
	pn.DataSync.CreateUser().
		ID("user-alice").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Alice", "type": "shopper"}).
		Execute()
	pn.DataSync.RemoveChannel().ID("channel-summer-sale").Execute()
	pn.DataSync.CreateChannel().
		ID("channel-summer-sale").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Summer Sale", "type": "promotion"}).
		Execute()
	pn.DataSync.RemoveMembership().ID("membership-alice-summer-sale").Execute()
	pn.DataSync.CreateMembership().
		ID("membership-alice-summer-sale").
		ChannelID("channel-summer-sale").
		UserID("user-alice").
		RelationshipClassVersion(1).
		Payload(map[string]interface{}{"role": "viewer"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.GetMembership().
		ID("membership-alice-summer-sale").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Membership retrieved, status: %d\n", status.StatusCode)

	// Output:
	// Membership retrieved, status: 200
}

// snippet.get_memberships_basic_usage
// Example_getMembershipsBasicUsage demonstrates listing memberships
func Example_getMembershipsBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("user-alice").Execute()
	pn.DataSync.CreateUser().
		ID("user-alice").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Alice", "type": "shopper"}).
		Execute()
	pn.DataSync.RemoveChannel().ID("channel-summer-sale").Execute()
	pn.DataSync.CreateChannel().
		ID("channel-summer-sale").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Summer Sale", "type": "promotion"}).
		Execute()
	pn.DataSync.RemoveMembership().ID("membership-alice-summer-sale").Execute()
	pn.DataSync.CreateMembership().
		ID("membership-alice-summer-sale").
		ChannelID("channel-summer-sale").
		UserID("user-alice").
		RelationshipClassVersion(1).
		Payload(map[string]interface{}{"role": "viewer"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.GetMemberships().
		UserID("user-alice").
		Limit(20).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Memberships listed, status: %d\n", status.StatusCode)

	// Output:
	// Memberships listed, status: 200
}

// snippet.get_memberships_by_channel_id
// Example_getMembershipsByChannelID demonstrates listing a channel's members
func Example_getMembershipsByChannelID() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("user-alice").Execute()
	pn.DataSync.CreateUser().
		ID("user-alice").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Alice", "type": "shopper"}).
		Execute()
	pn.DataSync.RemoveChannel().ID("channel-summer-sale").Execute()
	pn.DataSync.CreateChannel().
		ID("channel-summer-sale").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Summer Sale", "type": "promotion"}).
		Execute()
	pn.DataSync.RemoveMembership().ID("membership-alice-summer-sale").Execute()
	pn.DataSync.CreateMembership().
		ID("membership-alice-summer-sale").
		ChannelID("channel-summer-sale").
		UserID("user-alice").
		RelationshipClassVersion(1).
		Payload(map[string]interface{}{"role": "viewer"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.GetMemberships().
		ChannelID("channel-summer-sale").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Channel members listed, status: %d\n", status.StatusCode)

	// Output:
	// Channel members listed, status: 200
}

// snippet.get_memberships_filter_fast
// Example_getMembershipsFilterFast demonstrates filtering memberships with FilterFast
func Example_getMembershipsFilterFast() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("user-alice").Execute()
	pn.DataSync.CreateUser().
		ID("user-alice").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Alice", "type": "shopper"}).
		Execute()
	pn.DataSync.RemoveChannel().ID("channel-summer-sale").Execute()
	pn.DataSync.CreateChannel().
		ID("channel-summer-sale").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Summer Sale", "type": "promotion"}).
		Execute()
	pn.DataSync.RemoveMembership().ID("membership-alice-summer-sale").Execute()
	pn.DataSync.CreateMembership().
		ID("membership-alice-summer-sale").
		ChannelID("channel-summer-sale").
		UserID("user-alice").
		RelationshipClassVersion(1).
		Payload(map[string]interface{}{"role": "viewer"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.GetMemberships().
		FilterFast("role == \"viewer\"").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Filtered memberships retrieved, status: %d\n", status.StatusCode)

	// Output:
	// Filtered memberships retrieved, status: 200
}

// snippet.get_memberships_filter
// Example_getMembershipsFilter demonstrates filtering memberships with Filter
func Example_getMembershipsFilter() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("user-alice").Execute()
	pn.DataSync.CreateUser().
		ID("user-alice").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Alice", "type": "shopper"}).
		Execute()
	pn.DataSync.RemoveChannel().ID("channel-summer-sale").Execute()
	pn.DataSync.CreateChannel().
		ID("channel-summer-sale").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Summer Sale", "type": "promotion"}).
		Execute()
	pn.DataSync.RemoveMembership().ID("membership-alice-summer-sale").Execute()
	pn.DataSync.CreateMembership().
		ID("membership-alice-summer-sale").
		ChannelID("channel-summer-sale").
		UserID("user-alice").
		RelationshipClassVersion(1).
		Payload(map[string]interface{}{"role": "viewer"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.GetMemberships().
		Filter("role LIKE \"*view*\"").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Filtered memberships retrieved, status: %d\n", status.StatusCode)

	// Output:
	// Filtered memberships retrieved, status: 200
}

// snippet.get_memberships_pagination
// Example_getMembershipsPagination demonstrates paging through the membership list
func Example_getMembershipsPagination() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("user-alice").Execute()
	pn.DataSync.CreateUser().
		ID("user-alice").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Alice", "type": "shopper"}).
		Execute()
	pn.DataSync.RemoveChannel().ID("channel-summer-sale").Execute()
	pn.DataSync.CreateChannel().
		ID("channel-summer-sale").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Summer Sale", "type": "promotion"}).
		Execute()
	pn.DataSync.RemoveMembership().ID("membership-alice-summer-sale").Execute()
	pn.DataSync.CreateMembership().
		ID("membership-alice-summer-sale").
		ChannelID("channel-summer-sale").
		UserID("user-alice").
		RelationshipClassVersion(1).
		Payload(map[string]interface{}{"role": "viewer"}).
		Execute()
	// snippet.show

	res, status, err := pn.DataSync.GetMemberships().
		Limit(20).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Memberships page retrieved, status: %d\n", status.StatusCode)

	if res.Meta != nil && res.Meta.HasNext {
		_, _, err = pn.DataSync.GetMemberships().
			Cursor(res.Meta.NextCursor).
			Limit(20).
			Execute()

		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}

	// Output:
	// Memberships page retrieved, status: 200
}

// snippet.set_membership_basic_usage
// Example_setMembershipBasicUsage demonstrates replacing a membership in full
func Example_setMembershipBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("user-alice").Execute()
	pn.DataSync.CreateUser().
		ID("user-alice").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Alice", "type": "shopper"}).
		Execute()
	pn.DataSync.RemoveChannel().ID("channel-summer-sale").Execute()
	pn.DataSync.CreateChannel().
		ID("channel-summer-sale").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Summer Sale", "type": "promotion"}).
		Execute()
	pn.DataSync.RemoveMembership().ID("membership-alice-summer-sale").Execute()
	pn.DataSync.CreateMembership().
		ID("membership-alice-summer-sale").
		ChannelID("channel-summer-sale").
		UserID("user-alice").
		RelationshipClassVersion(1).
		Payload(map[string]interface{}{"role": "viewer"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.SetMembership().
		ID("membership-alice-summer-sale").
		RelationshipClassVersion(1).
		Payload(map[string]interface{}{"role": "moderator"}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Membership replaced, status: %d\n", status.StatusCode)

	// Output:
	// Membership replaced, status: 200
}

// snippet.update_membership_basic_usage
// Example_updateMembershipBasicUsage demonstrates applying a JSON Patch update to a membership
func Example_updateMembershipBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("user-alice").Execute()
	pn.DataSync.CreateUser().
		ID("user-alice").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Alice", "type": "shopper"}).
		Execute()
	pn.DataSync.RemoveChannel().ID("channel-summer-sale").Execute()
	pn.DataSync.CreateChannel().
		ID("channel-summer-sale").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Summer Sale", "type": "promotion"}).
		Execute()
	pn.DataSync.RemoveMembership().ID("membership-alice-summer-sale").Execute()
	pn.DataSync.CreateMembership().
		ID("membership-alice-summer-sale").
		ChannelID("channel-summer-sale").
		UserID("user-alice").
		RelationshipClassVersion(1).
		Payload(map[string]interface{}{"role": "viewer"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.UpdateMembership().
		ID("membership-alice-summer-sale").
		Replace("/status", "active").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Membership updated, status: %d\n", status.StatusCode)

	// Output:
	// Membership updated, status: 200
}

// snippet.remove_membership_basic_usage
// Example_removeMembershipBasicUsage demonstrates deleting a membership
func Example_removeMembershipBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("user-alice").Execute()
	pn.DataSync.CreateUser().
		ID("user-alice").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Alice", "type": "shopper"}).
		Execute()
	pn.DataSync.RemoveChannel().ID("channel-summer-sale").Execute()
	pn.DataSync.CreateChannel().
		ID("channel-summer-sale").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Summer Sale", "type": "promotion"}).
		Execute()
	pn.DataSync.RemoveMembership().ID("membership-alice-summer-sale").Execute()
	pn.DataSync.CreateMembership().
		ID("membership-alice-summer-sale").
		ChannelID("channel-summer-sale").
		UserID("user-alice").
		RelationshipClassVersion(1).
		Payload(map[string]interface{}{"role": "viewer"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.RemoveMembership().
		ID("membership-alice-summer-sale").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Membership removed, status: %d\n", status.StatusCode)

	// Output:
	// Membership removed, status: 200
}

// ==================== Entities ====================

// snippet.create_entity_basic_usage
// Example_createEntityBasicUsage demonstrates creating a custom DataSync entity
func Example_createEntityBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveEntity().ID("product-sneaker-42").Execute()
	// snippet.show

	_, status, err := pn.DataSync.CreateEntity().
		ID("product-sneaker-42").
		EntityClass("product").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Retro Sneaker", "price": 89.99, "stock": 12}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Entity created, status: %d\n", status.StatusCode)

	// Output:
	// Entity created, status: 200
}

// snippet.get_entity_basic_usage
// Example_getEntityBasicUsage demonstrates reading a single entity
func Example_getEntityBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveEntity().ID("product-sneaker-42").Execute()
	pn.DataSync.CreateEntity().
		ID("product-sneaker-42").
		EntityClass("product").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Retro Sneaker", "price": 89.99, "stock": 12}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.GetEntity().
		ID("product-sneaker-42").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Entity retrieved, status: %d\n", status.StatusCode)

	// Output:
	// Entity retrieved, status: 200
}

// snippet.get_entities_basic_usage
// Example_getEntitiesBasicUsage demonstrates listing entities of a class
func Example_getEntitiesBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveEntity().ID("product-sneaker-42").Execute()
	pn.DataSync.CreateEntity().
		ID("product-sneaker-42").
		EntityClass("product").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Retro Sneaker", "price": 89.99, "stock": 12}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.GetEntities().
		EntityClass("product").
		Limit(20).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Entities listed, status: %d\n", status.StatusCode)

	// Output:
	// Entities listed, status: 200
}

// snippet.get_entities_filter_fast
// Example_getEntitiesFilterFast demonstrates filtering entities with FilterFast
func Example_getEntitiesFilterFast() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveEntity().ID("product-sneaker-42").Execute()
	pn.DataSync.CreateEntity().
		ID("product-sneaker-42").
		EntityClass("product").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Retro Sneaker", "price": 89.99, "stock": 12}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.GetEntities().
		EntityClass("product").
		FilterFast("price < 100 && stock > 0").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Filtered entities retrieved, status: %d\n", status.StatusCode)

	// Output:
	// Filtered entities retrieved, status: 200
}

// snippet.get_entities_filter
// Example_getEntitiesFilter demonstrates filtering entities with Filter
func Example_getEntitiesFilter() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveEntity().ID("product-sneaker-42").Execute()
	pn.DataSync.CreateEntity().
		ID("product-sneaker-42").
		EntityClass("product").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Retro Sneaker", "price": 89.99, "stock": 12}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.GetEntities().
		EntityClass("product").
		Filter("name LIKE \"*Sneaker*\" || price > 500").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Filtered entities retrieved, status: %d\n", status.StatusCode)

	// Output:
	// Filtered entities retrieved, status: 200
}

// snippet.get_entities_pagination
// Example_getEntitiesPagination demonstrates paging through the entity list
func Example_getEntitiesPagination() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveEntity().ID("product-sneaker-42").Execute()
	pn.DataSync.CreateEntity().
		ID("product-sneaker-42").
		EntityClass("product").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Retro Sneaker", "price": 89.99, "stock": 12}).
		Execute()
	// snippet.show

	res, status, err := pn.DataSync.GetEntities().
		EntityClass("product").
		Limit(20).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Entities page retrieved, status: %d\n", status.StatusCode)

	if res.Meta != nil && res.Meta.HasNext {
		_, _, err = pn.DataSync.GetEntities().
			EntityClass("product").
			Cursor(res.Meta.NextCursor).
			Limit(20).
			Execute()

		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}

	// Output:
	// Entities page retrieved, status: 200
}

// snippet.set_entity_basic_usage
// Example_setEntityBasicUsage demonstrates replacing an entity in full
func Example_setEntityBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveEntity().ID("product-sneaker-42").Execute()
	pn.DataSync.CreateEntity().
		ID("product-sneaker-42").
		EntityClass("product").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Retro Sneaker", "price": 89.99, "stock": 12}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.SetEntity().
		ID("product-sneaker-42").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Retro Sneaker", "price": 79.99, "stock": 12}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Entity replaced, status: %d\n", status.StatusCode)

	// Output:
	// Entity replaced, status: 200
}

// snippet.update_entity_basic_usage
// Example_updateEntityBasicUsage demonstrates applying a JSON Patch update to an entity
func Example_updateEntityBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveEntity().ID("product-sneaker-42").Execute()
	pn.DataSync.CreateEntity().
		ID("product-sneaker-42").
		EntityClass("product").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Retro Sneaker", "price": 89.99, "stock": 12}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.UpdateEntity().
		ID("product-sneaker-42").
		Replace("/payload/stock", 8).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Entity updated, status: %d\n", status.StatusCode)

	// Output:
	// Entity updated, status: 200
}

// snippet.update_entity_multiple_operations
// Example_updateEntityMultipleOperations demonstrates combining patch operations with a concurrency guard
func Example_updateEntityMultipleOperations() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveEntity().ID("product-sneaker-42").Execute()
	pn.DataSync.CreateEntity().
		ID("product-sneaker-42").
		EntityClass("product").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Retro Sneaker", "price": 89.99, "stock": 12}).
		Execute()
	// snippet.show

	current, _, err := pn.DataSync.GetEntity().
		ID("product-sneaker-42").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	_, status, err := pn.DataSync.UpdateEntity().
		ID("product-sneaker-42").
		Add("/payload/tags", []string{"sale"}).
		Replace("/payload/stock", 8).
		IfMatchETag(current.Data.ETag).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Entity patched with guard, status: %d\n", status.StatusCode)

	// Output:
	// Entity patched with guard, status: 200
}

// snippet.remove_entity_basic_usage
// Example_removeEntityBasicUsage demonstrates deleting an entity
func Example_removeEntityBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveEntity().ID("product-sneaker-42").Execute()
	pn.DataSync.CreateEntity().
		ID("product-sneaker-42").
		EntityClass("product").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Retro Sneaker", "price": 89.99, "stock": 12}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.RemoveEntity().
		ID("product-sneaker-42").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Entity removed, status: %d\n", status.StatusCode)

	// Output:
	// Entity removed, status: 200
}

// ==================== Relationships ====================

// snippet.create_relationship_basic_usage
// Example_createRelationshipBasicUsage demonstrates linking two entities
func Example_createRelationshipBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("seller-bob").Execute()
	pn.DataSync.CreateUser().
		ID("seller-bob").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Bob"}).
		Execute()
	pn.DataSync.RemoveEntity().ID("product-sneaker-42").Execute()
	pn.DataSync.CreateEntity().
		ID("product-sneaker-42").
		EntityClass("product").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Retro Sneaker", "price": 89.99, "stock": 12}).
		Execute()
	pn.DataSync.RemoveRelationship().ID("rel-bob-owns-sneaker-42").Execute()
	// snippet.show

	_, status, err := pn.DataSync.CreateRelationship().
		ID("rel-bob-owns-sneaker-42").
		EntityAID("seller-bob").
		EntityBID("product-sneaker-42").
		RelationshipClass("ProductOwner").
		RelationshipClassVersion(1).
		Payload(map[string]interface{}{"since": "2026-07-13"}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Relationship created, status: %d\n", status.StatusCode)

	// Output:
	// Relationship created, status: 200
}

// snippet.get_relationship_basic_usage
// Example_getRelationshipBasicUsage demonstrates reading a single relationship
func Example_getRelationshipBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("seller-bob").Execute()
	pn.DataSync.CreateUser().
		ID("seller-bob").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Bob"}).
		Execute()
	pn.DataSync.RemoveEntity().ID("product-sneaker-42").Execute()
	pn.DataSync.CreateEntity().
		ID("product-sneaker-42").
		EntityClass("product").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Retro Sneaker", "price": 89.99, "stock": 12}).
		Execute()
	pn.DataSync.RemoveRelationship().ID("rel-bob-owns-sneaker-42").Execute()
	pn.DataSync.CreateRelationship().
		ID("rel-bob-owns-sneaker-42").
		EntityAID("seller-bob").
		EntityBID("product-sneaker-42").
		RelationshipClass("ProductOwner").
		RelationshipClassVersion(1).
		Payload(map[string]interface{}{"since": "2026-07-13"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.GetRelationship().
		ID("rel-bob-owns-sneaker-42").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Relationship retrieved, status: %d\n", status.StatusCode)

	// Output:
	// Relationship retrieved, status: 200
}

// snippet.get_relationships_basic_usage
// Example_getRelationshipsBasicUsage demonstrates listing relationships of a class
func Example_getRelationshipsBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("seller-bob").Execute()
	pn.DataSync.CreateUser().
		ID("seller-bob").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Bob"}).
		Execute()
	pn.DataSync.RemoveEntity().ID("product-sneaker-42").Execute()
	pn.DataSync.CreateEntity().
		ID("product-sneaker-42").
		EntityClass("product").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Retro Sneaker", "price": 89.99, "stock": 12}).
		Execute()
	pn.DataSync.RemoveRelationship().ID("rel-bob-owns-sneaker-42").Execute()
	pn.DataSync.CreateRelationship().
		ID("rel-bob-owns-sneaker-42").
		EntityAID("seller-bob").
		EntityBID("product-sneaker-42").
		RelationshipClass("ProductOwner").
		RelationshipClassVersion(1).
		Payload(map[string]interface{}{"since": "2026-07-13"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.GetRelationships().
		RelationshipClass("ProductOwner").
		Limit(20).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Relationships listed, status: %d\n", status.StatusCode)

	// Output:
	// Relationships listed, status: 200
}

// snippet.get_relationships_by_entity_b_id
// Example_getRelationshipsByEntityBID demonstrates listing an entity's incoming links
func Example_getRelationshipsByEntityBID() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("seller-bob").Execute()
	pn.DataSync.CreateUser().
		ID("seller-bob").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Bob"}).
		Execute()
	pn.DataSync.RemoveEntity().ID("product-sneaker-42").Execute()
	pn.DataSync.CreateEntity().
		ID("product-sneaker-42").
		EntityClass("product").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Retro Sneaker", "price": 89.99, "stock": 12}).
		Execute()
	pn.DataSync.RemoveRelationship().ID("rel-bob-owns-sneaker-42").Execute()
	pn.DataSync.CreateRelationship().
		ID("rel-bob-owns-sneaker-42").
		EntityAID("seller-bob").
		EntityBID("product-sneaker-42").
		RelationshipClass("ProductOwner").
		RelationshipClassVersion(1).
		Payload(map[string]interface{}{"since": "2026-07-13"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.GetRelationships().
		RelationshipClass("ProductOwner").
		EntityBID("product-sneaker-42").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Incoming relationships listed, status: %d\n", status.StatusCode)

	// Output:
	// Incoming relationships listed, status: 200
}

// snippet.get_relationships_filter_fast
// Example_getRelationshipsFilterFast demonstrates filtering relationships with FilterFast
func Example_getRelationshipsFilterFast() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("seller-bob").Execute()
	pn.DataSync.CreateUser().
		ID("seller-bob").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Bob"}).
		Execute()
	pn.DataSync.RemoveEntity().ID("product-sneaker-42").Execute()
	pn.DataSync.CreateEntity().
		ID("product-sneaker-42").
		EntityClass("product").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Retro Sneaker", "price": 89.99, "stock": 12}).
		Execute()
	pn.DataSync.RemoveRelationship().ID("rel-bob-owns-sneaker-42").Execute()
	pn.DataSync.CreateRelationship().
		ID("rel-bob-owns-sneaker-42").
		EntityAID("seller-bob").
		EntityBID("product-sneaker-42").
		RelationshipClass("ProductOwner").
		RelationshipClassVersion(1).
		Payload(map[string]interface{}{"since": "2026-07-13"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.GetRelationships().
		RelationshipClass("ProductOwner").
		FilterFast("since == \"2026-07-13\"").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Filtered relationships retrieved, status: %d\n", status.StatusCode)

	// Output:
	// Filtered relationships retrieved, status: 200
}

// snippet.get_relationships_filter
// Example_getRelationshipsFilter demonstrates filtering relationships with Filter
func Example_getRelationshipsFilter() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("seller-bob").Execute()
	pn.DataSync.CreateUser().
		ID("seller-bob").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Bob"}).
		Execute()
	pn.DataSync.RemoveEntity().ID("product-sneaker-42").Execute()
	pn.DataSync.CreateEntity().
		ID("product-sneaker-42").
		EntityClass("product").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Retro Sneaker", "price": 89.99, "stock": 12}).
		Execute()
	pn.DataSync.RemoveRelationship().ID("rel-bob-owns-sneaker-42").Execute()
	pn.DataSync.CreateRelationship().
		ID("rel-bob-owns-sneaker-42").
		EntityAID("seller-bob").
		EntityBID("product-sneaker-42").
		RelationshipClass("ProductOwner").
		RelationshipClassVersion(1).
		Payload(map[string]interface{}{"since": "2026-07-13"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.GetRelationships().
		RelationshipClass("ProductOwner").
		Filter("tier LIKE \"*gold*\"").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Filtered relationships retrieved, status: %d\n", status.StatusCode)

	// Output:
	// Filtered relationships retrieved, status: 200
}

// snippet.get_relationships_pagination
// Example_getRelationshipsPagination demonstrates paging through the relationship list
func Example_getRelationshipsPagination() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("seller-bob").Execute()
	pn.DataSync.CreateUser().
		ID("seller-bob").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Bob"}).
		Execute()
	pn.DataSync.RemoveEntity().ID("product-sneaker-42").Execute()
	pn.DataSync.CreateEntity().
		ID("product-sneaker-42").
		EntityClass("product").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Retro Sneaker", "price": 89.99, "stock": 12}).
		Execute()
	pn.DataSync.RemoveRelationship().ID("rel-bob-owns-sneaker-42").Execute()
	pn.DataSync.CreateRelationship().
		ID("rel-bob-owns-sneaker-42").
		EntityAID("seller-bob").
		EntityBID("product-sneaker-42").
		RelationshipClass("ProductOwner").
		RelationshipClassVersion(1).
		Payload(map[string]interface{}{"since": "2026-07-13"}).
		Execute()
	// snippet.show

	res, status, err := pn.DataSync.GetRelationships().
		RelationshipClass("ProductOwner").
		Limit(20).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Relationships page retrieved, status: %d\n", status.StatusCode)

	if res.Meta != nil && res.Meta.HasNext {
		_, _, err = pn.DataSync.GetRelationships().
			RelationshipClass("ProductOwner").
			Cursor(res.Meta.NextCursor).
			Limit(20).
			Execute()

		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}

	// Output:
	// Relationships page retrieved, status: 200
}

// snippet.set_relationship_basic_usage
// Example_setRelationshipBasicUsage demonstrates replacing a relationship in full
func Example_setRelationshipBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("seller-bob").Execute()
	pn.DataSync.CreateUser().
		ID("seller-bob").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Bob"}).
		Execute()
	pn.DataSync.RemoveEntity().ID("product-sneaker-42").Execute()
	pn.DataSync.CreateEntity().
		ID("product-sneaker-42").
		EntityClass("product").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Retro Sneaker", "price": 89.99, "stock": 12}).
		Execute()
	pn.DataSync.RemoveRelationship().ID("rel-bob-owns-sneaker-42").Execute()
	pn.DataSync.CreateRelationship().
		ID("rel-bob-owns-sneaker-42").
		EntityAID("seller-bob").
		EntityBID("product-sneaker-42").
		RelationshipClass("ProductOwner").
		RelationshipClassVersion(1).
		Payload(map[string]interface{}{"since": "2026-07-13"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.SetRelationship().
		ID("rel-bob-owns-sneaker-42").
		RelationshipClassVersion(1).
		Payload(map[string]interface{}{"since": "2026-07-13", "tier": "gold"}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Relationship replaced, status: %d\n", status.StatusCode)

	// Output:
	// Relationship replaced, status: 200
}

// snippet.update_relationship_basic_usage
// Example_updateRelationshipBasicUsage demonstrates applying a JSON Patch update to a relationship
func Example_updateRelationshipBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("seller-bob").Execute()
	pn.DataSync.CreateUser().
		ID("seller-bob").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Bob"}).
		Execute()
	pn.DataSync.RemoveEntity().ID("product-sneaker-42").Execute()
	pn.DataSync.CreateEntity().
		ID("product-sneaker-42").
		EntityClass("product").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Retro Sneaker", "price": 89.99, "stock": 12}).
		Execute()
	pn.DataSync.RemoveRelationship().ID("rel-bob-owns-sneaker-42").Execute()
	pn.DataSync.CreateRelationship().
		ID("rel-bob-owns-sneaker-42").
		EntityAID("seller-bob").
		EntityBID("product-sneaker-42").
		RelationshipClass("ProductOwner").
		RelationshipClassVersion(1).
		Payload(map[string]interface{}{"since": "2026-07-13", "tier": "gold"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.UpdateRelationship().
		ID("rel-bob-owns-sneaker-42").
		Replace("/payload/tier", "platinum").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Relationship updated, status: %d\n", status.StatusCode)

	// Output:
	// Relationship updated, status: 200
}

// snippet.remove_relationship_basic_usage
// Example_removeRelationshipBasicUsage demonstrates deleting a relationship
func Example_removeRelationshipBasicUsage() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	pn.DataSync.RemoveUser().ID("seller-bob").Execute()
	pn.DataSync.CreateUser().
		ID("seller-bob").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Bob"}).
		Execute()
	pn.DataSync.RemoveEntity().ID("product-sneaker-42").Execute()
	pn.DataSync.CreateEntity().
		ID("product-sneaker-42").
		EntityClass("product").
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Retro Sneaker", "price": 89.99, "stock": 12}).
		Execute()
	pn.DataSync.RemoveRelationship().ID("rel-bob-owns-sneaker-42").Execute()
	pn.DataSync.CreateRelationship().
		ID("rel-bob-owns-sneaker-42").
		EntityAID("seller-bob").
		EntityBID("product-sneaker-42").
		RelationshipClass("ProductOwner").
		RelationshipClassVersion(1).
		Payload(map[string]interface{}{"since": "2026-07-13"}).
		Execute()
	// snippet.show

	_, status, err := pn.DataSync.RemoveRelationship().
		ID("rel-bob-owns-sneaker-42").
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Relationship removed, status: %d\n", status.StatusCode)

	// Output:
	// Relationship removed, status: 200
}

// ==================== Real-time updates ====================

// snippet.data_sync_event_listener
// Example_dataSyncEventListener demonstrates receiving DataSync events with a listener
func Example_dataSyncEventListener() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	listener := pubnub.NewListener()
	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-listener.DataSyncEvent:
				if event.Event == pubnub.PNDataSyncEventDelete {
					fmt.Printf("DataSync %s deleted: %s\n", event.Type, event.ID)
					continue
				}
				fmt.Printf("DataSync %s %s on channel %s\n", event.Type, event.Event, event.Channel)
			case <-done:
				return
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().
		Channels([]string{"product-sneaker-42"}).
		Execute()

	fmt.Println("Subscribed to DataSync events")

	pn.UnsubscribeAll()
	close(done)

	// snippet.hide
	time.Sleep(100 * time.Millisecond)
	// snippet.show

	// Output:
	// Subscribed to DataSync events
}

// snippet.subscribe_to_projection_channels
// Example_subscribeToProjectionChannels demonstrates subscribing to both a base and a projection channel
func Example_subscribeToProjectionChannels() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	pn.Subscribe().
		Channels([]string{"product-sneaker-42", "__admin__product-sneaker-42"}).
		Execute()

	fmt.Println("Subscribed to base and admin projection channels")

	pn.UnsubscribeAll()

	// snippet.hide
	time.Sleep(100 * time.Millisecond)
	// snippet.show

	// Output:
	// Subscribed to base and admin projection channels
}

// ==================== Access Manager ====================

// snippet.grant_token_data_sync
// Example_grantTokenDataSync demonstrates granting DataSync user and resource permissions
func Example_grantTokenDataSync() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	_, status, err := pn.GrantToken().
		TTL(60).
		AuthorizedUUID("my-authorized-uuid").
		Users(map[string]pubnub.UUIDPermissions{
			"user-bob": {
				Get: true,
			},
		}).
		DataSync(pubnub.PNDataSyncTokenScopes{
			Entities: map[string]pubnub.DataSyncPermissions{
				"product-sneaker-42": {
					Get:    true,
					Create: true,
					Update: true,
					Delete: true,
				},
			},
			Memberships: map[string]pubnub.DataSyncPermissions{
				"membership-alice-summer-sale": {
					Get: true,
				},
			},
		}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Token granted, status: %d\n", status.StatusCode)

	// Output:
	// Token granted, status: 200
}

// snippet.grant_token_data_sync_projection
// Example_grantTokenDataSyncProjection demonstrates granting a token scoped to a single projection
func Example_grantTokenDataSyncProjection() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"     // Replace with your subscribe key
	config.PublishKey = "demo"       // Replace with your publish key
	config.SecretKey = "demo-secret" // Replace with your secret key (required for Access Manager)

	// snippet.hide
	config = setPubnubExampleDataSyncConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	_, status, err := pn.GrantToken().
		TTL(60).
		AuthorizedUUID("my-authorized-uuid").
		DataSync(pubnub.PNDataSyncTokenScopes{
			Entities: map[string]pubnub.DataSyncPermissions{
				"product-sneaker-42": {
					Get:    true,
					Create: true,
					Update: true,
				},
			},
		}).
		DataSyncProjections(pubnub.PNDataSyncProjections{
			Resources: pubnub.PNDataSyncProjectionScope{
				Entities: map[string]string{
					"product-sneaker-42": "admin",
				},
			},
		}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Token granted, status: %d\n", status.StatusCode)

	// Output:
	// Token granted, status: 200
}
