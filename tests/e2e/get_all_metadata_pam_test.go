package e2e

import (
	"testing"

	pubnub "github.com/pubnub/go/v10"
	"github.com/stretchr/testify/assert"
)

func TestGetAllMetadataPAMWithoutToken(t *testing.T) {
	client := newPAMClient(t)
	assertGetAllMetadataDenied(t, client)
}

func TestGetAllMetadataPAMWithCategoryGrant(t *testing.T) {
	admin := newPAMAdmin(t)
	client := newPAMClient(t)

	res, status, err := admin.GrantToken().
		TTL(10).
		Categories(pubnub.PNGrantCategories{
			Channels: pubnub.CategoryPermissions{Get: true},
			UUIDs:    pubnub.CategoryPermissions{Get: true},
		}).
		Execute()
	if err != nil || res == nil || res.Data.Token == "" {
		t.Fatalf("GrantToken categories: err=%v status=%d body=%s", err, status.StatusCode, status.OriginalResponse)
	}

	parsed, parseErr := pubnub.ParseToken(res.Data.Token)
	assert.Nil(t, parseErr)
	assert.True(t, parsed.Categories.Channels.Get)
	assert.True(t, parsed.Categories.UUIDs.Get)

	client.SetToken(res.Data.Token)
	assertGetAllMetadataAllowed(t, client)
}

func assertGetAllMetadataDenied(t *testing.T, client *pubnub.PubNub) {
	t.Helper()

	_, uuidStatus, uuidErr := client.GetAllUUIDMetadata().Limit(1).Execute()
	assertServerCode(t, uuidErr, uuidStatus, 403)

	_, channelStatus, channelErr := client.GetAllChannelMetadata().Limit(1).Execute()
	assertServerCode(t, channelErr, channelStatus, 403)
}

func assertGetAllMetadataAllowed(t *testing.T, client *pubnub.PubNub) {
	t.Helper()

	uuidRes, uuidStatus, uuidErr := client.GetAllUUIDMetadata().Limit(1).Execute()
	if uuidErr != nil {
		t.Fatalf("GetAllUUIDMetadata: err=%v status=%d body=%s", uuidErr, uuidStatus.StatusCode, serverBody(uuidErr))
	}
	assert.Equal(t, 200, uuidStatus.StatusCode)
	assert.NotNil(t, uuidRes)
	assert.NotNil(t, uuidRes.Data)

	channelRes, channelStatus, channelErr := client.GetAllChannelMetadata().Limit(1).Execute()
	if channelErr != nil {
		t.Fatalf("GetAllChannelMetadata: err=%v status=%d body=%s", channelErr, channelStatus.StatusCode, serverBody(channelErr))
	}
	assert.Equal(t, 200, channelStatus.StatusCode)
	assert.NotNil(t, channelRes)
	assert.NotNil(t, channelRes.Data)
}

func serverBody(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func newPAMAdmin(t *testing.T) *pubnub.PubNub {
	t.Helper()
	requirePAMKeys(t)
	pn := pubnub.NewPubNub(pamConfigCopy())
	t.Cleanup(pn.Destroy)
	return pn
}

func newPAMClient(t *testing.T) *pubnub.PubNub {
	t.Helper()
	requirePAMKeys(t)
	cfg := pamConfigCopy()
	cfg.SecretKey = ""
	pn := pubnub.NewPubNub(cfg)
	t.Cleanup(pn.Destroy)
	return pn
}

func requirePAMKeys(t *testing.T) {
	t.Helper()
	if pamConfig.PublishKey == "" || pamConfig.SubscribeKey == "" || pamConfig.SecretKey == "" {
		t.Skip("PAM_PUBLISH_KEY / PAM_SUBSCRIBE_KEY / PAM_SECRET_KEY not set")
	}
}
