package contract

import (
	"context"
	"fmt"
	"reflect"

	pubnub "github.com/pubnub/go/v6"
)

func theTTL(ctx context.Context, arg1 int) error {
	state := getAccessState(ctx)
	state.TTL = arg1
	return nil
}

func grantPermissionREAD(ctx context.Context) error {
	state := getAccessState(ctx)

	switch v := state.CurrentPermissions.(type) {
	case *pubnub.ChannelPermissions:
		v.Read = true
	case *pubnub.GroupPermissions:
		v.Read = true
	default:
		return fmt.Errorf("Not expected type %s", reflect.TypeOf(v).String())
	}
	return nil
}

func grantPermissionDELETE(ctx context.Context) error {
	state := getAccessState(ctx)

	switch v := state.CurrentPermissions.(type) {
	case *pubnub.ChannelPermissions:
		v.Delete = true
	case *pubnub.UUIDPermissions:
		v.Delete = true
	default:
		return fmt.Errorf("Not expected type %s", reflect.TypeOf(v).String())
	}
	return nil
}

func grantPermissionGET(ctx context.Context) error {
	state := getAccessState(ctx)

	switch v := state.CurrentPermissions.(type) {
	case *pubnub.ChannelPermissions:
		v.Get = true
	case *pubnub.UUIDPermissions:
		v.Get = true
	default:
		return fmt.Errorf("Not expected type %s", reflect.TypeOf(v).String())
	}
	return nil
}

func grantPermissionJOIN(ctx context.Context) error {
	state := getAccessState(ctx)

	switch v := state.CurrentPermissions.(type) {
	case *pubnub.ChannelPermissions:
		v.Join = true
	default:
		return fmt.Errorf("Not expected type %s", reflect.TypeOf(v).String())
	}
	return nil
}

func grantPermissionMANAGE(ctx context.Context) error {
	state := getAccessState(ctx)

	switch v := state.CurrentPermissions.(type) {
	case *pubnub.ChannelPermissions:
		v.Manage = true
	case *pubnub.GroupPermissions:
		v.Manage = true
	default:
		return fmt.Errorf("Not expected type %s", reflect.TypeOf(v).String())
	}
	return nil
}

func grantPermissionUPDATE(ctx context.Context) error {
	state := getAccessState(ctx)

	switch v := state.CurrentPermissions.(type) {
	case *pubnub.ChannelPermissions:
		v.Update = true
	case *pubnub.UUIDPermissions:
		v.Update = true
	default:
		return fmt.Errorf("Not expected type %s", reflect.TypeOf(v).String())
	}
	return nil
}

func grantPermissionWRITE(ctx context.Context) error {
	state := getAccessState(ctx)

	switch v := state.CurrentPermissions.(type) {
	case *pubnub.ChannelPermissions:
		v.Write = true
	default:
		return fmt.Errorf("Not expected type %s", reflect.TypeOf(v).String())
	}
	return nil
}

func theCHANNELResourceAccessPermissions(ctx context.Context, channel string) error {
	state := getAccessState(ctx)

	permissions := pubnub.ChannelPermissions{}
	state.ChannelPermissions[channel] = &permissions
	state.CurrentPermissions = &permissions

	return nil
}

func theCHANNELPatternAccessPermissions(ctx context.Context, pattern string) error {
	state := getAccessState(ctx)

	permissions := pubnub.ChannelPermissions{}
	state.ChannelPatternPermissions[pattern] = &permissions
	state.CurrentPermissions = &permissions

	return nil
}

func theCHANNEL_GROUPResourceAccessPermissions(ctx context.Context, id string) error {
	state := getAccessState(ctx)

	permissions := pubnub.GroupPermissions{}
	state.ChannelGroupPermissions[id] = &permissions
	state.CurrentPermissions = &permissions

	return nil
}

func theCHANNEL_GROUPPatternAccessPermissions(ctx context.Context, pattern string) error {
	state := getAccessState(ctx)

	permissions := pubnub.GroupPermissions{}
	state.ChannelGroupPatternPermissions[pattern] = &permissions
	state.CurrentPermissions = &permissions

	return nil
}

func theUUIDResourceAccessPermissions(ctx context.Context, id string) error {
	state := getAccessState(ctx)

	permissions := pubnub.UUIDPermissions{}
	state.UUIDPermissions[id] = &permissions
	state.CurrentPermissions = &permissions

	return nil
}

func theUUIDPatternAccessPermissions(ctx context.Context, pattern string) error {
	state := getAccessState(ctx)

	permissions := pubnub.UUIDPermissions{}
	state.UUIDPatternPermissions[pattern] = &permissions
	state.CurrentPermissions = &permissions

	return nil
}

func theAuthorizedUUID(ctx context.Context, uuid string) error {
	state := getAccessState(ctx)
	state.AuthorizedUUID = uuid

	return nil
}

const tokenWithEverything = "qEF2AkF0GmEI03xDdHRsGDxDcmVzpURjaGFuoWljaGFubmVsLTEY70NncnChb2NoYW5uZWxfZ3JvdXAtMQVDdXNyoENzcGOgRHV1aWShZnV1aWQtMRhoQ3BhdKVEY2hhbqFtXmNoYW5uZWwtXFMqJBjvQ2dycKF0XjpjaGFubmVsX2dyb3VwLVxTKiQFQ3VzcqBDc3BjoER1dWlkoWpedXVpZC1cUyokGGhEbWV0YaBEdXVpZHR0ZXN0LWF1dGhvcml6ZWQtdXVpZENzaWdYIPpU-vCe9rkpYs87YUrFNWkyNq8CVvmKwEjVinnDrJJc"

func iHaveAKnownTokenWithEverything(ctx context.Context) error {
	state := getAccessState(ctx)
	state.TokenString = tokenWithEverything
	return nil
}

func denyResourcePermissionGET(ctx context.Context) error {
	state := getAccessState(ctx)

	switch v := state.CurrentPermissions.(type) {
	case *pubnub.ChannelPermissions:
		v.Get = false
	case *pubnub.UUIDPermissions:
		v.Get = false
	default:
		return fmt.Errorf("Not expected type %s", reflect.TypeOf(v).String())
	}
	return nil
}

func aToken(ctx context.Context) error {
	state := getAccessState(ctx)
	state.TokenString = tokenWithEverything
	return nil
}

func aValidTokenWithPermissionsToPublishWithChannelChannel(ctx context.Context) error {
	state := getAccessState(ctx)
	state.TokenString = tokenWithEverything
	return nil
}

func anExpiredTokenWithPermissionsToPublishWithChannelChannel(ctx context.Context) error {
	state := getAccessState(ctx)
	state.TokenString = tokenWithEverything
	return nil
}

func theTokenString(ctx context.Context, token string) error {
	state := getAccessState(ctx)
	state.TokenString = token
	return nil
}

