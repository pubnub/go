package pubnub_samples_test

import (
	"os"

	pubnub "github.com/pubnub/go/v10"
)

func setPubnubExampleConfigData(config *pubnub.Config) *pubnub.Config {
	config.SetUserId("GO_SDK_EXAMPLE_USER")
	config.PublishKey = os.Getenv("PUBLISH_KEY")
	config.SubscribeKey = os.Getenv("SUBSCRIBE_KEY")

	return config
}

func setPubnubExampleConfigDataWithSecretKey(config *pubnub.Config) *pubnub.Config {
	config.SetUserId("GO_SDK_EXAMPLE_USER")
	config.PublishKey = os.Getenv("PUBLISH_KEY")
	config.SubscribeKey = os.Getenv("SUBSCRIBE_KEY")
	config.SecretKey = os.Getenv("SECRET_KEY")

	return config
}

func setPubnubExamplePAMConfigData(config *pubnub.Config) *pubnub.Config {
	config.SetUserId("GO_SDK_EXAMPLE_USER")
	config.PublishKey = os.Getenv("PAM_PUBLISH_KEY")
	config.SubscribeKey = os.Getenv("PAM_SUBSCRIBE_KEY")
	config.SecretKey = os.Getenv("PAM_SECRET_KEY")

	return config
}

func setPubnubExampleDataSyncConfigData(config *pubnub.Config) *pubnub.Config {
	config.SetUserId("GO_SDK_EXAMPLE_USER")
	config.PublishKey = os.Getenv("DS_PUBLISH_KEY")
	config.SubscribeKey = os.Getenv("DS_SUBSCRIBE_KEY")
	config.SecretKey = os.Getenv("DS_SECRET_KEY")

	return config
}

const (
	dsExampleUserID       = "user-alice"
	dsExampleChannelID    = "channel-summer-sale"
	dsExampleMembershipID = "membership-alice-summer-sale"
)

// resetDataSyncExampleUser deletes the shared example membership (if any) then the user.
func resetDataSyncExampleUser(pn *pubnub.PubNub) {
	_, _, _ = pn.DataSync.RemoveMembership().ID(dsExampleMembershipID).Execute()
	_, _, _ = pn.DataSync.RemoveUser().ID(dsExampleUserID).Execute()
}

// resetDataSyncExampleChannel deletes the shared example membership (if any) then the channel.
func resetDataSyncExampleChannel(pn *pubnub.PubNub) {
	_, _, _ = pn.DataSync.RemoveMembership().ID(dsExampleMembershipID).Execute()
	_, _, _ = pn.DataSync.RemoveChannel().ID(dsExampleChannelID).Execute()
}

func seedDataSyncExampleUser(pn *pubnub.PubNub) {
	resetDataSyncExampleUser(pn)
	_, _, _ = pn.DataSync.CreateUser().
		ID(dsExampleUserID).
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Alice", "type": "shopper"}).
		Execute()
}

func seedDataSyncExampleChannel(pn *pubnub.PubNub) {
	resetDataSyncExampleChannel(pn)
	_, _, _ = pn.DataSync.CreateChannel().
		ID(dsExampleChannelID).
		EntityClassVersion(1).
		Payload(map[string]interface{}{"name": "Summer Sale", "type": "promotion"}).
		Execute()
}

func seedDataSyncExampleMembership(pn *pubnub.PubNub) {
	seedDataSyncExampleUser(pn)
	seedDataSyncExampleChannel(pn)
	_, _, _ = pn.DataSync.CreateMembership().
		ID(dsExampleMembershipID).
		ChannelID(dsExampleChannelID).
		UserID(dsExampleUserID).
		RelationshipClassVersion(1).
		Status("active").
		Payload(map[string]interface{}{"role": "viewer"}).
		Execute()
}
