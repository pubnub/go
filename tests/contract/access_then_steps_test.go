package contract

import (
	"context"
	"fmt"

	pubnub "github.com/pubnub/go/v5"
)

func theTokenContainsTheTTL(ctx context.Context, expectedTTL int) error {
	state := getAccessState(ctx)
	if expectedTTL != state.ParsedToken.TTL {
		return fmt.Errorf("Expected %d but found %d", expectedTTL, state.ParsedToken.TTL)
	}

	return nil
}

func theTokenDoesNotContainAnAuthorizedUuid(ctx context.Context) error {
	state := getAccessState(ctx)
	if state.ParsedToken.AuthorizedUUID != "" {
		return fmt.Errorf("Expected empty AuthorizedUUID but found %s", state.ParsedToken.AuthorizedUUID)
	}

	return nil
}

func theTokenHasCHANNELResourceAccessPermissions(ctx context.Context, channel string) error {
	state := getAccessState(ctx)
	permissions, ok := state.ParsedToken.Resources.Channels[channel]
	if !ok {
		return fmt.Errorf("Expected channel %s in ParsedToken", channel)
	}

	state.ResourcePermissions = permissions

	return nil
}

func resourceHasPermission(ctx context.Context, perm string) error {
	state := getAccessState(ctx)
	resourcePermissions := state.ResourcePermissions
	switch v := resourcePermissions.(type) {
	case pubnub.ChannelPermissions:
		answer, err := channelHasPermission(v, perm)
		if err != nil {
			return err
		} else if !answer {
			return fmt.Errorf("Resource expected to have %s", perm)
		}
	case pubnub.GroupPermissions:
		answer, err := channelGroupHasPermission(v, perm)
		if err != nil {
			return err
		} else if !answer {
			return fmt.Errorf("Resource expected to have %s", perm)
		}

	case pubnub.UUIDPermissions:
		answer, err := uuidHasPermission(v, perm)
		if err != nil {
			return err
		} else if !answer {
			return fmt.Errorf("Resource expected to have %s", perm)
		}

	}

	return nil
}

func channelHasPermission(permissions pubnub.ChannelPermissions, perm string) (bool, error) {
	switch perm {
	case "READ":
		return permissions.Read, nil
	case "DELETE":
		return permissions.Delete, nil
	case "JOIN":
		return permissions.Join, nil
	case "MANAGE":
		return permissions.Manage, nil
	case "UPDATE":
		return permissions.Update, nil
	case "GET":
		return permissions.Get, nil
	case "WRITE":
		return permissions.Write, nil
	default:
		return permissions.Read, fmt.Errorf("Unsupported permissions %s", perm)
	}
}

func channelGroupHasPermission(permissions pubnub.GroupPermissions, perm string) (bool, error) {
	switch perm {
	case "READ":
		return permissions.Read, nil
	case "MANAGE":
		return permissions.Manage, nil
	default:
		return permissions.Read, fmt.Errorf("Unsupported permissions %s", perm)
	}
}

func uuidHasPermission(permissions pubnub.UUIDPermissions, perm string) (bool, error) {
	switch perm {
	case "DELETE":
		return permissions.Delete, nil
	case "UPDATE":
		return permissions.Update, nil
	case "GET":
		return permissions.Get, nil
	default:
		return permissions.Get, fmt.Errorf("Unsupported permissions %s", perm)
	}
}

func theTokenContainsTheAuthorizedUUID(ctx context.Context, uuid string) error {
	state := getAccessState(ctx)
	if state.ParsedToken.AuthorizedUUID != uuid {
		return fmt.Errorf("Expected %s but found %s", uuid, state.ParsedToken.AuthorizedUUID)
	}

	return nil
}

func theTokenHasCHANNELPatternAccessPermissions(ctx context.Context, pattern string) error {
	state := getAccessState(ctx)
	permissions, ok := state.ParsedToken.Patterns.Channels[pattern]
	if !ok {
		return fmt.Errorf("Expected channel %s in ParsedToken", pattern)
	}

	state.ResourcePermissions = permissions

	return nil
}

func theTokenHasCHANNEL_GROUPResourceAccessPermissions(ctx context.Context, id string) error {
	state := getAccessState(ctx)
	permissions, ok := state.ParsedToken.Resources.ChannelGroups[id]
	if !ok {
		return fmt.Errorf("Expected group %s in ParsedToken", id)
	}

	state.ResourcePermissions = permissions

	return nil
}

func theTokenHasCHANNEL_GROUPPatternAccessPermissions(ctx context.Context, pattern string) error {
	state := getAccessState(ctx)
	permissions, ok := state.ParsedToken.Patterns.ChannelGroups[pattern]
	if !ok {
		return fmt.Errorf("Expected group %s in ParsedToken", pattern)
	}

	state.ResourcePermissions = permissions

	return nil
}

func theTokenHasUUIDPatternAccessPermissions(ctx context.Context, pattern string) error {
	state := getAccessState(ctx)
	permissions, ok := state.ParsedToken.Patterns.UUIDs[pattern]
	if !ok {
		return fmt.Errorf("Expected uuid %s in ParsedToken", pattern)
	}

	state.ResourcePermissions = permissions

	return nil
}

func theTokenHasUUIDResourceAccessPermissions(ctx context.Context, id string) error {
	state := getAccessState(ctx)
	permissions, ok := state.ParsedToken.Resources.UUIDs[id]
	if !ok {
		return fmt.Errorf("Expected uuid %s in ParsedToken", id)
	}

	state.ResourcePermissions = permissions

	return nil
}
