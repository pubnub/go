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
	config = setPubnubExampleConfigData(config)  // <- Skip this line (for testing only)
	defer pn.RemoveUUIDMetadata().Execute()      // <- Skip this line (cleanup for tests)
	// snippet.show

When copying examples to your own code:
- Use your own publish/subscribe keys instead of the "demo" keys
- Remove any statements that are between snippet.hide and snippet.show (they're only for testing purposes)
*/

// snippet.set_uuid_metadata
// Example_setUUIDMetadata demonstrates setting user metadata (UUID metadata)
func Example_setUUIDMetadata() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo" // Replace with your subscribe key
	config.PublishKey = "demo"   // Replace with your publish key

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	// Cleanup: Remove metadata after test
	defer pn.RemoveUUIDMetadata().UUID("user-123").Execute()
	// snippet.show

	// Set user metadata with profile information
	response, status, err := pn.SetUUIDMetadata().
		UUID("user-123").                                // User ID
		Name("John Doe").                                // Display name
		Email("john.doe@example.com").                   // Email address
		ProfileURL("https://example.com/profiles/john"). // Profile URL
		ExternalID("ext-123").                           // External system ID
		Custom(map[string]interface{}{                   // Custom metadata
			"role":     "admin",
			"language": "en",
		}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Printf("User metadata set for UUID: %s\n", response.Data.ID)
	}

	// Output:
	// User metadata set for UUID: user-123
}

// snippet.get_uuid_metadata
// Example_getUUIDMetadata demonstrates retrieving user metadata
func Example_getUUIDMetadata() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	defer pn.RemoveUUIDMetadata().UUID("user-456").Execute()
	// snippet.show

	// First, set user metadata
	pn.SetUUIDMetadata().
		UUID("user-456").
		Name("Jane Smith").
		Email("jane@example.com").
		Execute()

	// Then retrieve the user metadata
	response, status, err := pn.GetUUIDMetadata().
		UUID("user-456"). // User ID to retrieve
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Printf("User: %s\n", response.Data.Name)
		fmt.Printf("Email: %s\n", response.Data.Email)
	}

	// Output:
	// User: Jane Smith
	// Email: jane@example.com
}

// snippet.get_uuid_metadata_with_includes
// Example_getUUIDMetadataWithIncludes demonstrates retrieving user metadata with all include options
func Example_getUUIDMetadataWithIncludes() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	defer pn.RemoveUUIDMetadata().UUID("user-with-all").Execute()
	// snippet.show

	// Create user metadata with all fields
	pn.SetUUIDMetadata().
		UUID("user-with-all").
		Name("Complete User").
		Custom(map[string]interface{}{
			"department": "engineering",
		}).
		Status("active").
		Type("employee").
		Execute()

	// Get user metadata with all include options
	response, status, err := pn.GetUUIDMetadata().
		UUID("user-with-all").
		Include([]pubnub.PNUUIDMetadataInclude{
			pubnub.PNUUIDMetadataIncludeCustom, // Include custom fields
			pubnub.PNUUIDMetadataIncludeStatus, // Include status field
			pubnub.PNUUIDMetadataIncludeType,   // Include type field
		}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Printf("User: %s\n", response.Data.Name)
		fmt.Printf("Status: %s\n", response.Data.Status)
		fmt.Printf("Type: %s\n", response.Data.Type)
		if dept, ok := response.Data.Custom["department"].(string); ok {
			fmt.Printf("Department: %s\n", dept)
		}
	}

	// Output:
	// User: Complete User
	// Status: active
	// Type: employee
	// Department: engineering
}

// snippet.get_all_uuid_metadata
// Example_getAllUUIDMetadata demonstrates listing all user metadata
func Example_getAllUUIDMetadata() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	defer pn.RemoveUUIDMetadata().UUID("user-list-1").Execute()
	defer pn.RemoveUUIDMetadata().UUID("user-list-2").Execute()
	// snippet.show

	// Create some user metadata first
	pn.SetUUIDMetadata().UUID("user-list-1").Name("Alice").Execute()
	pn.SetUUIDMetadata().UUID("user-list-2").Name("Bob").Execute()

	// Get all user metadata with pagination
	_, status, err := pn.GetAllUUIDMetadata().
		Limit(10).   // Number of results per page
		Count(true). // Include total count
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("Retrieved user metadata list successfully")
	}

	// Output:
	// Retrieved user metadata list successfully
}

// snippet.get_all_uuid_metadata_with_includes
// Example_getAllUUIDMetadataWithIncludes demonstrates listing all users with all include options
func Example_getAllUUIDMetadataWithIncludes() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	defer pn.RemoveUUIDMetadata().UUID("user-all-list-1").Execute()
	// snippet.show

	// Create user with all fields
	pn.SetUUIDMetadata().
		UUID("user-all-list-1").
		Name("Full User").
		Custom(map[string]interface{}{"role": "admin"}).
		Status("active").
		Type("staff").
		Execute()

	// Get all users with all include options
	response, status, err := pn.GetAllUUIDMetadata().
		Include([]pubnub.PNUUIDMetadataInclude{
			pubnub.PNUUIDMetadataIncludeCustom, // Include custom fields
			pubnub.PNUUIDMetadataIncludeStatus, // Include status field
			pubnub.PNUUIDMetadataIncludeType,   // Include type field
		}).
		Limit(10).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 && len(response.Data) > 0 {
		fmt.Println("Retrieved users with all metadata fields")
	}

	// Output:
	// Retrieved users with all metadata fields
}

// snippet.remove_uuid_metadata
// Example_removeUUIDMetadata demonstrates removing user metadata
func Example_removeUUIDMetadata() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	// Setup: Create metadata first
	pn.SetUUIDMetadata().UUID("user-to-remove").Name("Temp User").Execute()
	// snippet.show

	// Remove user metadata
	_, status, err := pn.RemoveUUIDMetadata().
		UUID("user-to-remove"). // User ID to remove
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("User metadata removed successfully")
	}

	// Output:
	// User metadata removed successfully
}

// snippet.set_channel_metadata
// Example_setChannelMetadata demonstrates setting channel metadata
func Example_setChannelMetadata() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	// Cleanup: Remove metadata after test
	defer pn.RemoveChannelMetadata().Channel("support-channel-123").Execute()
	// snippet.show

	// Set channel metadata with descriptive information
	response, status, err := pn.SetChannelMetadata().
		Channel("support-channel-123").            // Channel ID
		Name("Customer Support").                  // Display name
		Description("24/7 customer support chat"). // Description
		Custom(map[string]interface{}{             // Custom metadata
			"department": "support",
			"priority":   "high",
		}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Printf("Channel metadata set for: %s\n", response.Data.ID)
	}

	// Output:
	// Channel metadata set for: support-channel-123
}

// snippet.get_channel_metadata
// Example_getChannelMetadata demonstrates retrieving channel metadata
func Example_getChannelMetadata() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	defer pn.RemoveChannelMetadata().Channel("sales-channel-456").Execute()
	// snippet.show

	// First, set channel metadata
	pn.SetChannelMetadata().
		Channel("sales-channel-456").
		Name("Sales Team").
		Description("Sales team discussions").
		Execute()

	// Then retrieve the channel metadata
	response, status, err := pn.GetChannelMetadata().
		Channel("sales-channel-456"). // Channel ID to retrieve
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Printf("Channel: %s\n", response.Data.Name)
		fmt.Printf("Description: %s\n", response.Data.Description)
	}

	// Output:
	// Channel: Sales Team
	// Description: Sales team discussions
}

// snippet.get_channel_metadata_with_includes
// Example_getChannelMetadataWithIncludes demonstrates retrieving channel metadata with all include options
func Example_getChannelMetadataWithIncludes() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	defer pn.RemoveChannelMetadata().Channel("channel-with-all").Execute()
	// snippet.show

	// Create channel metadata with all fields
	pn.SetChannelMetadata().
		Channel("channel-with-all").
		Name("Complete Channel").
		Description("Full featured channel").
		Custom(map[string]interface{}{
			"category": "general",
		}).
		Status("active").
		Type("public").
		Execute()

	// Get channel metadata with all include options
	response, status, err := pn.GetChannelMetadata().
		Channel("channel-with-all").
		Include([]pubnub.PNChannelMetadataInclude{
			pubnub.PNChannelMetadataIncludeCustom, // Include custom fields
			pubnub.PNChannelMetadataIncludeStatus, // Include status field
			pubnub.PNChannelMetadataIncludeType,   // Include type field
		}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Printf("Channel: %s\n", response.Data.Name)
		fmt.Printf("Status: %s\n", response.Data.Status)
		fmt.Printf("Type: %s\n", response.Data.Type)
		if cat, ok := response.Data.Custom["category"].(string); ok {
			fmt.Printf("Category: %s\n", cat)
		}
	}

	// Output:
	// Channel: Complete Channel
	// Status: active
	// Type: public
	// Category: general
}

// snippet.get_all_channel_metadata
// Example_getAllChannelMetadata demonstrates listing all channel metadata
func Example_getAllChannelMetadata() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	defer pn.RemoveChannelMetadata().Channel("channel-list-1").Execute()
	defer pn.RemoveChannelMetadata().Channel("channel-list-2").Execute()
	// snippet.show

	// Create some channel metadata first
	pn.SetChannelMetadata().Channel("channel-list-1").Name("General").Execute()
	pn.SetChannelMetadata().Channel("channel-list-2").Name("Support").Execute()

	// Get all channel metadata with pagination
	_, status, err := pn.GetAllChannelMetadata().
		Limit(10).   // Number of results per page
		Count(true). // Include total count
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("Retrieved channel metadata list successfully")
	}

	// Output:
	// Retrieved channel metadata list successfully
}

// snippet.get_all_channel_metadata_with_includes
// Example_getAllChannelMetadataWithIncludes demonstrates listing all channels with all include options
func Example_getAllChannelMetadataWithIncludes() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	defer pn.RemoveChannelMetadata().Channel("channel-all-list-1").Execute()
	// snippet.show

	// Create channel with all fields
	pn.SetChannelMetadata().
		Channel("channel-all-list-1").
		Name("Full Channel").
		Custom(map[string]interface{}{"priority": "high"}).
		Status("active").
		Type("team").
		Execute()

	// Get all channels with all include options
	response, status, err := pn.GetAllChannelMetadata().
		Include([]pubnub.PNChannelMetadataInclude{
			pubnub.PNChannelMetadataIncludeCustom, // Include custom fields
			pubnub.PNChannelMetadataIncludeStatus, // Include status field
			pubnub.PNChannelMetadataIncludeType,   // Include type field
		}).
		Limit(10).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 && len(response.Data) > 0 {
		fmt.Println("Retrieved channels with all metadata fields")
	}

	// Output:
	// Retrieved channels with all metadata fields
}

// snippet.remove_channel_metadata
// Example_removeChannelMetadata demonstrates removing channel metadata
func Example_removeChannelMetadata() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	// Setup: Create metadata first
	pn.SetChannelMetadata().Channel("temp-channel").Name("Temporary Channel").Execute()
	// snippet.show

	// Remove channel metadata
	_, status, err := pn.RemoveChannelMetadata().
		Channel("temp-channel"). // Channel ID to remove
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("Channel metadata removed successfully")
	}

	// Output:
	// Channel metadata removed successfully
}

// snippet.set_memberships
// Example_setMemberships demonstrates adding a user to channels (memberships)
func Example_setMemberships() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	// Setup: Create user and channels first
	pn.SetUUIDMetadata().UUID("member-user-789").Name("Member User").Execute()
	pn.SetChannelMetadata().Channel("channel-member-1").Name("Channel 1").Execute()
	pn.SetChannelMetadata().Channel("channel-member-2").Name("Channel 2").Execute()
	defer pn.RemoveUUIDMetadata().UUID("member-user-789").Execute()
	defer pn.RemoveChannelMetadata().Channel("channel-member-1").Execute()
	defer pn.RemoveChannelMetadata().Channel("channel-member-2").Execute()
	// snippet.show

	// Add user to multiple channels with custom metadata
	response, status, err := pn.SetMemberships().
		UUID("member-user-789").       // User ID
		Set([]pubnub.PNMembershipsSet{ // Channels to join
			{
				Channel: pubnub.PNMembershipsChannel{ID: "channel-member-1"},
				Custom: map[string]interface{}{
					"role": "moderator",
				},
			},
			{
				Channel: pubnub.PNMembershipsChannel{ID: "channel-member-2"},
				Custom: map[string]interface{}{
					"role": "member",
				},
			},
		}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Printf("User added to %d channel(s)\n", len(response.Data))
	}

	// Output:
	// User added to 2 channel(s)
}

// snippet.get_memberships
// Example_getMemberships demonstrates retrieving channels a user belongs to
func Example_getMemberships() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	defer pn.RemoveUUIDMetadata().UUID("member-user-999").Execute()
	defer pn.RemoveChannelMetadata().Channel("channel-get-1").Execute()
	// snippet.show

	// Create user and channel
	pn.SetUUIDMetadata().UUID("member-user-999").Name("John").Execute()
	pn.SetChannelMetadata().Channel("channel-get-1").Name("Engineering").Execute()

	// Add user to the channel
	pn.SetMemberships().
		UUID("member-user-999").
		Set([]pubnub.PNMembershipsSet{
			{
				Channel: pubnub.PNMembershipsChannel{ID: "channel-get-1"},
				Custom: map[string]interface{}{
					"role": "developer",
				},
			},
		}).
		Execute()

	// Get all channels the user is a member of
	response, status, err := pn.GetMemberships().
		UUID("member-user-999"). // User ID
		Include([]pubnub.PNMembershipsInclude{
			pubnub.PNMembershipsIncludeCustom,  // Include custom data
			pubnub.PNMembershipsIncludeChannel, // Include channel details
		}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 && len(response.Data) > 0 {
		membership := response.Data[0]
		fmt.Printf("User is member of: %s\n", membership.Channel.Name)
		if role, ok := membership.Custom["role"].(string); ok {
			fmt.Printf("Role: %s\n", role)
		}
	}

	// Output:
	// User is member of: Engineering
	// Role: developer
}

// snippet.get_memberships_with_all_includes
// Example_getMembershipsWithAllIncludes demonstrates getting memberships with all include options
func Example_getMembershipsWithAllIncludes() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	defer pn.RemoveUUIDMetadata().UUID("user-full-membership").Execute()
	defer pn.RemoveChannelMetadata().Channel("channel-full-membership").Execute()
	// snippet.show

	// Create user and channel with all fields
	pn.SetUUIDMetadata().
		UUID("user-full-membership").
		Name("Full Member").
		Status("active").
		Type("user").
		Execute()

	pn.SetChannelMetadata().
		Channel("channel-full-membership").
		Name("Full Channel").
		Custom(map[string]interface{}{"category": "tech"}).
		Status("active").
		Type("public").
		Execute()

	// Add user to channel with custom membership data
	pn.SetMemberships().
		UUID("user-full-membership").
		Set([]pubnub.PNMembershipsSet{
			{
				Channel: pubnub.PNMembershipsChannel{ID: "channel-full-membership"},
				Custom: map[string]interface{}{
					"role": "moderator",
				},
				Status: "active",
				Type:   "membership",
			},
		}).
		Execute()

	// Get memberships with all include options
	response, status, err := pn.GetMemberships().
		UUID("user-full-membership").
		Include([]pubnub.PNMembershipsInclude{
			pubnub.PNMembershipsIncludeCustom,        // Include membership custom data
			pubnub.PNMembershipsIncludeChannel,       // Include channel details
			pubnub.PNMembershipsIncludeChannelCustom, // Include channel custom data
			pubnub.PNMembershipsIncludeChannelStatus, // Include channel status
			pubnub.PNMembershipsIncludeChannelType,   // Include channel type
			pubnub.PNMembershipsIncludeStatus,        // Include membership status
			pubnub.PNMembershipsIncludeType,          // Include membership type
		}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 && len(response.Data) > 0 {
		membership := response.Data[0]
		fmt.Printf("Channel: %s\n", membership.Channel.Name)
		fmt.Printf("Channel Status: %s\n", membership.Channel.Status)
		fmt.Printf("Membership Status: %s\n", membership.Status)
		if role, ok := membership.Custom["role"].(string); ok {
			fmt.Printf("Role: %s\n", role)
		}
	}

	// Output:
	// Channel: Full Channel
	// Channel Status: active
	// Membership Status: active
	// Role: moderator
}

// snippet.remove_memberships
// Example_removeMemberships demonstrates removing a user from channels (removing memberships)
func Example_removeMemberships() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	// Setup: Create channels and add memberships
	pn.SetChannelMetadata().Channel("temp-channel-1").Name("Temp 1").Execute()
	pn.SetChannelMetadata().Channel("temp-channel-2").Name("Temp 2").Execute()
	pn.SetUUIDMetadata().UUID("remove-user-888").Name("Bob").Execute()
	defer pn.RemoveChannelMetadata().Channel("temp-channel-1").Execute()
	defer pn.RemoveChannelMetadata().Channel("temp-channel-2").Execute()
	defer pn.RemoveUUIDMetadata().UUID("remove-user-888").Execute()

	pn.SetMemberships().
		UUID("remove-user-888").
		Set([]pubnub.PNMembershipsSet{
			{Channel: pubnub.PNMembershipsChannel{ID: "temp-channel-1"}},
			{Channel: pubnub.PNMembershipsChannel{ID: "temp-channel-2"}},
		}).
		Execute()
	// snippet.show

	// Remove user from specific channels
	response, status, err := pn.RemoveMemberships().
		UUID("remove-user-888").             // User ID
		Remove([]pubnub.PNMembershipsRemove{ // Channels to remove from
			{Channel: pubnub.PNMembershipsChannel{ID: "temp-channel-1"}},
		}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Printf("Removed from %d channel(s)\n", len(response.Data))
	}

	// Output:
	// Removed from 1 channel(s)
}

// snippet.manage_memberships
// Example_manageMemberships demonstrates adding and removing channel memberships in a single call
func Example_manageMemberships() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	// Setup: Create channels
	pn.SetChannelMetadata().Channel("add-channel-999").Name("Add Channel").Execute()
	pn.SetChannelMetadata().Channel("remove-channel-999").Name("Remove Channel").Execute()
	pn.SetUUIDMetadata().UUID("manage-user-999").Name("Charlie").Execute()
	defer pn.RemoveChannelMetadata().Channel("add-channel-999").Execute()
	defer pn.RemoveChannelMetadata().Channel("remove-channel-999").Execute()
	defer pn.RemoveUUIDMetadata().UUID("manage-user-999").Execute()

	// Add initial membership to be removed
	pn.SetMemberships().
		UUID("manage-user-999").
		Set([]pubnub.PNMembershipsSet{
			{Channel: pubnub.PNMembershipsChannel{ID: "remove-channel-999"}},
		}).
		Execute()
	// snippet.show

	// Manage memberships: add to some channels and remove from others
	response, status, err := pn.ManageMemberships().
		UUID("manage-user-999").       // User ID
		Set([]pubnub.PNMembershipsSet{ // Channels to add
			{
				Channel: pubnub.PNMembershipsChannel{ID: "add-channel-999"},
				Custom: map[string]interface{}{
					"role": "participant",
				},
			},
		}).
		Remove([]pubnub.PNMembershipsRemove{ // Channels to remove
			{Channel: pubnub.PNMembershipsChannel{ID: "remove-channel-999"}},
		}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Printf("Managed memberships: %d channel(s)\n", len(response.Data))
	}

	// Output:
	// Managed memberships: 1 channel(s)
}

// snippet.set_channel_members
// Example_setChannelMembers demonstrates adding users to a channel (channel members)
func Example_setChannelMembers() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	// Setup: Create channel and users first
	pn.SetChannelMetadata().Channel("team-channel-555").Name("Team Channel").Execute()
	pn.SetUUIDMetadata().UUID("team-member-1").Name("Alice").Execute()
	pn.SetUUIDMetadata().UUID("team-member-2").Name("Bob").Execute()
	defer pn.RemoveChannelMetadata().Channel("team-channel-555").Execute()
	defer pn.RemoveUUIDMetadata().UUID("team-member-1").Execute()
	defer pn.RemoveUUIDMetadata().UUID("team-member-2").Execute()
	// snippet.show

	// Add multiple users to a channel with custom metadata
	response, status, err := pn.SetChannelMembers().
		Channel("team-channel-555").      // Channel ID
		Set([]pubnub.PNChannelMembersSet{ // Users to add
			{
				UUID: pubnub.PNChannelMembersUUID{ID: "team-member-1"},
				Custom: map[string]interface{}{
					"role": "admin",
				},
			},
			{
				UUID: pubnub.PNChannelMembersUUID{ID: "team-member-2"},
				Custom: map[string]interface{}{
					"role": "member",
				},
			},
		}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Printf("Added %d member(s) to channel\n", len(response.Data))
	}

	// Output:
	// Added 2 member(s) to channel
}

// snippet.get_channel_members
// Example_getChannelMembers demonstrates retrieving members of a channel
func Example_getChannelMembers() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	defer pn.RemoveChannelMetadata().Channel("project-channel-777").Execute()
	defer pn.RemoveUUIDMetadata().UUID("project-member-1").Execute()
	// snippet.show

	// Create channel and user
	pn.SetChannelMetadata().Channel("project-channel-777").Name("Alpha Project").Execute()
	pn.SetUUIDMetadata().UUID("project-member-1").Name("Sarah").Execute()

	// Add user to the channel
	pn.SetChannelMembers().
		Channel("project-channel-777").
		Set([]pubnub.PNChannelMembersSet{
			{
				UUID: pubnub.PNChannelMembersUUID{ID: "project-member-1"},
				Custom: map[string]interface{}{
					"role": "lead",
				},
			},
		}).
		Execute()

	// Get all members of the channel
	response, status, err := pn.GetChannelMembers().
		Channel("project-channel-777"). // Channel ID
		Include([]pubnub.PNChannelMembersInclude{
			pubnub.PNChannelMembersIncludeCustom, // Include custom data
			pubnub.PNChannelMembersIncludeUUID,   // Include user details
		}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 && len(response.Data) > 0 {
		member := response.Data[0]
		fmt.Printf("Member: %s\n", member.UUID.Name)
		if role, ok := member.Custom["role"].(string); ok {
			fmt.Printf("Role: %s\n", role)
		}
	}

	// Output:
	// Member: Sarah
	// Role: lead
}

// snippet.get_channel_members_with_all_includes
// Example_getChannelMembersWithAllIncludes demonstrates getting channel members with all include options
func Example_getChannelMembersWithAllIncludes() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	defer pn.RemoveChannelMetadata().Channel("channel-full-members").Execute()
	defer pn.RemoveUUIDMetadata().UUID("user-full-member").Execute()
	// snippet.show

	// Create channel and user with all fields
	pn.SetChannelMetadata().
		Channel("channel-full-members").
		Name("Complete Team").
		Status("active").
		Type("team").
		Execute()

	pn.SetUUIDMetadata().
		UUID("user-full-member").
		Name("Complete Member").
		Custom(map[string]interface{}{"department": "dev"}).
		Status("active").
		Type("employee").
		Execute()

	// Add user to channel with custom member data
	pn.SetChannelMembers().
		Channel("channel-full-members").
		Set([]pubnub.PNChannelMembersSet{
			{
				UUID: pubnub.PNChannelMembersUUID{ID: "user-full-member"},
				Custom: map[string]interface{}{
					"access": "admin",
				},
				Status: "active",
				Type:   "member",
			},
		}).
		Execute()

	// Get channel members with all include options
	response, status, err := pn.GetChannelMembers().
		Channel("channel-full-members").
		Include([]pubnub.PNChannelMembersInclude{
			pubnub.PNChannelMembersIncludeCustom,     // Include member custom data
			pubnub.PNChannelMembersIncludeUUID,       // Include user details
			pubnub.PNChannelMembersIncludeUUIDCustom, // Include user custom data
			pubnub.PNChannelMembersIncludeUUIDStatus, // Include user status
			pubnub.PNChannelMembersIncludeUUIDType,   // Include user type
			pubnub.PNChannelMembersIncludeStatus,     // Include member status
			pubnub.PNChannelMembersIncludeType,       // Include member type
		}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 && len(response.Data) > 0 {
		member := response.Data[0]
		fmt.Printf("Member: %s\n", member.UUID.Name)
		fmt.Printf("User Status: %s\n", member.UUID.Status)
		fmt.Printf("Member Status: %s\n", member.Status)
		if access, ok := member.Custom["access"].(string); ok {
			fmt.Printf("Access: %s\n", access)
		}
	}

	// Output:
	// Member: Complete Member
	// User Status: active
	// Member Status: active
	// Access: admin
}

// snippet.remove_channel_members
// Example_removeChannelMembers demonstrates removing users from a channel
func Example_removeChannelMembers() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	// Setup: Create channel and users, then add them as members
	pn.SetChannelMetadata().Channel("remove-members-channel").Name("Team").Execute()
	pn.SetUUIDMetadata().UUID("remove-member-1").Name("Alice").Execute()
	pn.SetUUIDMetadata().UUID("remove-member-2").Name("Bob").Execute()
	defer pn.RemoveChannelMetadata().Channel("remove-members-channel").Execute()
	defer pn.RemoveUUIDMetadata().UUID("remove-member-1").Execute()
	defer pn.RemoveUUIDMetadata().UUID("remove-member-2").Execute()

	pn.SetChannelMembers().
		Channel("remove-members-channel").
		Set([]pubnub.PNChannelMembersSet{
			{UUID: pubnub.PNChannelMembersUUID{ID: "remove-member-1"}},
			{UUID: pubnub.PNChannelMembersUUID{ID: "remove-member-2"}},
		}).
		Execute()
	// snippet.show

	// Remove specific users from the channel
	response, status, err := pn.RemoveChannelMembers().
		Channel("remove-members-channel").      // Channel ID
		Remove([]pubnub.PNChannelMembersRemove{ // Users to remove
			{UUID: pubnub.PNChannelMembersUUID{ID: "remove-member-1"}},
		}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Printf("Removed %d member(s) from channel\n", len(response.Data))
	}

	// Output:
	// Removed 1 member(s) from channel
}

// snippet.manage_channel_members
// Example_manageChannelMembers demonstrates adding and removing channel members in a single call
func Example_manageChannelMembers() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	// Setup: Create channel and users
	pn.SetChannelMetadata().Channel("manage-members-channel").Name("Project Team").Execute()
	pn.SetUUIDMetadata().UUID("new-member-111").Name("David").Execute()
	pn.SetUUIDMetadata().UUID("old-member-111").Name("Eve").Execute()
	defer pn.RemoveChannelMetadata().Channel("manage-members-channel").Execute()
	defer pn.RemoveUUIDMetadata().UUID("new-member-111").Execute()
	defer pn.RemoveUUIDMetadata().UUID("old-member-111").Execute()

	// Add initial member to be removed
	pn.SetChannelMembers().
		Channel("manage-members-channel").
		Set([]pubnub.PNChannelMembersSet{
			{UUID: pubnub.PNChannelMembersUUID{ID: "old-member-111"}},
		}).
		Execute()
	// snippet.show

	// Manage channel members: add new members and remove others
	response, status, err := pn.ManageChannelMembers().
		Channel("manage-members-channel"). // Channel ID
		Set([]pubnub.PNChannelMembersSet{  // Members to add
			{
				UUID: pubnub.PNChannelMembersUUID{ID: "new-member-111"},
				Custom: map[string]interface{}{
					"role": "developer",
				},
			},
		}).
		Remove([]pubnub.PNChannelMembersRemove{ // Members to remove
			{UUID: pubnub.PNChannelMembersUUID{ID: "old-member-111"}},
		}).
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Printf("Managed channel members: %d member(s)\n", len(response.Data))
	}

	// Output:
	// Managed channel members: 1 member(s)
}

// snippet.end
