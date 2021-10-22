package contract

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/pubnub/go/v6/pnerr"
)

func theErrorContains(ctx context.Context, substr string) error {
	state := getCommonState(ctx)

	if !strings.Contains(state.err.Error(), substr) {
		return fmt.Errorf("Expecting error containing '%s' but found '%s'", substr, state.err.Error())
	}

	return nil
}

func theErrorStatusCodeIs(ctx context.Context, errorCode int) error {
	state := getCommonState(ctx)
	switch v := state.err.(type) {
	default:
		return fmt.Errorf("Expecting *pnerr.ServerError but found type %T", v)
	case *pnerr.ServerError:
		if v.StatusCode != errorCode {
			return fmt.Errorf("Expecting %d but found %d", errorCode, v.StatusCode)
		}
	}

	return nil
}

func anErrorIsReturned(ctx context.Context) error {
	state := getCommonState(ctx)
	err := state.err

	fmt.Println(err)

	if err == nil {
		return errors.New("Expecting error, but found nil")
	}

	return nil
}
