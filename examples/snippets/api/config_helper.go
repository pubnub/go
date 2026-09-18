package pubnub_samples_test

import (
	"os"

	pubnub "github.com/pubnub/go/v9"
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
