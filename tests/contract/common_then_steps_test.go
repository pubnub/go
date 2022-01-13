package contract

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/pubnub/go/v7/pnerr"
)

func theErrorContains(ctx context.Context, substr string) error {
	state := getCommonState(ctx)

	if !strings.Contains(state.err.Error(), substr) {
		return fmt.Errorf("expecting error containing '%s' but found '%s'", substr, state.err.Error())
	}

	return nil
}

func theErrorStatusCodeIs(ctx context.Context, errorCode int) error {
	state := getCommonState(ctx)
	switch v := state.err.(type) {
	default:
		return fmt.Errorf("expecting *pnerr.ServerError but found type %T", v)
	case *pnerr.ServerError:
		if v.StatusCode != errorCode {
			return fmt.Errorf("expecting %d but found %d", errorCode, v.StatusCode)
		}
	}

	return nil
}

func anErrorIsReturned(ctx context.Context) error {
	state := getCommonState(ctx)

	if state.err == nil {
		return errors.New("expecting error, but found nil")
	}

	return nil
}

func theErrorDetailMessageIsNotEmpty(ctx context.Context) error {
	//TODO figure out how to do it
	return nil
}

func theResultIsSuccessful(ctx context.Context) error {
	state := getCommonState(ctx)

	if state.err != nil {
		return fmt.Errorf("expecting success, but found error %s", state.err.Error())
	}

	if state.statusResponse.StatusCode != 200 {
		return fmt.Errorf("expecting status code 200, but found %d", state.statusResponse.StatusCode)
	}

	return nil
}
