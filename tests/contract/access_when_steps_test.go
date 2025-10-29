package contract

import (
	"context"

	pubnub "github.com/pubnub/go/v8"
)

func iGrantATokenSpecifyingThosePermissions(ctx context.Context) error {
	aState := getAccessState(ctx)
	cState := getCommonState(ctx)

	grantToken := cState.pubNub.GrantToken()

	if len(aState.ChannelPermissions) != 0 {
		channelPermissions := map[string]pubnub.ChannelPermissions{}

		for name, permissions := range aState.ChannelPermissions {
			channelPermissions[name] = *permissions
		}
		grantToken.Channels(channelPermissions)
	}

	if len(aState.ChannelPatternPermissions) != 0 {
		channelPermissions := map[string]pubnub.ChannelPermissions{}

		for name, permissions := range aState.ChannelPatternPermissions {
			channelPermissions[name] = *permissions
		}
		grantToken.ChannelsPattern(channelPermissions)
	}

	if len(aState.ChannelGroupPermissions) != 0 {
		groupPermissions := map[string]pubnub.GroupPermissions{}

		for name, permissions := range aState.ChannelGroupPermissions {
			groupPermissions[name] = *permissions
		}
		grantToken.ChannelGroups(groupPermissions)
	}

	if len(aState.ChannelGroupPatternPermissions) != 0 {
		groupPermissions := map[string]pubnub.GroupPermissions{}

		for name, permissions := range aState.ChannelGroupPatternPermissions {
			groupPermissions[name] = *permissions
		}
		grantToken.ChannelGroupsPattern(groupPermissions)
	}

	if len(aState.UUIDPermissions) != 0 {
		uuidPermissions := map[string]pubnub.UUIDPermissions{}

		for name, permissions := range aState.UUIDPermissions {
			uuidPermissions[name] = *permissions
		}
		grantToken.UUIDs(uuidPermissions)
	}

	if len(aState.UUIDPatternPermissions) != 0 {
		uuidPermissions := map[string]pubnub.UUIDPermissions{}

		for name, permissions := range aState.UUIDPatternPermissions {
			uuidPermissions[name] = *permissions
		}
		grantToken.UUIDsPattern(uuidPermissions)
	}

	grantToken.TTL(aState.TTL)

	if len(aState.AuthorizedUUID) != 0 {
		grantToken.AuthorizedUUID(aState.AuthorizedUUID)
	}

	res, _, err := grantToken.Execute()
	if err != nil {
		return err
	}
	aState.GrantTokenResult = *res
	aState.ParsedToken, err = pubnub.ParseToken(res.Data.Token)

	if err != nil {
		return nil
	}

	return nil
}

func iAttemptToGrantATokenSpecifyingThosePermissions(ctx context.Context) error {
	err := iGrantATokenSpecifyingThosePermissions(ctx)

	cState := getCommonState(ctx)
	cState.err = err

	return nil
}

func iParseTheToken(ctx context.Context) error {
	aState := getAccessState(ctx)
	parsedToken, err := pubnub.ParseToken(aState.TokenString)

	if err != nil {
		return nil
	}
	aState.ParsedToken = parsedToken

	return nil
}

func iRevokeAToken(ctx context.Context) error {
	aState := getAccessState(ctx)
	cState := getCommonState(ctx)

	_, status, err := cState.pubNub.RevokeToken().Token(aState.TokenString).Execute()

	cState.statusResponse = status
	cState.err = err

	return nil
}

func iPublishAMessageUsingThatAuthTokenWithChannelChannel(ctx context.Context, channel string) error {
	aState := getAccessState(ctx)
	cState := getCommonState(ctx)

	cState.pubNub.SetToken(aState.TokenString)
	_, status, err := cState.pubNub.Publish().Channel(channel).Message("Message").Execute()
	cState.statusResponse = status
	cState.err = err

	return nil
}
