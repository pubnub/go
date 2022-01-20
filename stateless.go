package pubnub

import (
	"encoding/json"
)

func (pn *PubNub) handshake(
	ctx Context,
	channels []string,
	channelGroups []string,
) (chan subscribeEnvelope, chan error) {

	opts := &subscribeOpts{
		pubnub:           pn,
		Channels:         channels,
		ChannelGroups:    channelGroups,
		Timetoken:        0,
		Heartbeat:        pn.Config.PresenceTimeout,
		FilterExpression: pn.Config.FilterExpression,
		ctx:              ctx,
		QueryParam:       make(map[string]string),
	}

	resChan := make(chan subscribeEnvelope)
	errChan := make(chan error)

	go func() {
		res, _, err := executeRequest(opts)
		if err != nil {
			errChan <- err
		} else {
			var envelope subscribeEnvelope
			err = json.Unmarshal(res, &envelope)
			if err != nil {
				errChan <- err
			}
			resChan <- envelope
		}
		close(errChan)
		close(resChan)
	}()

	return resChan, errChan
}

func (pn *PubNub) receiveMessages(

	ctx Context,
	channels []string,
	channelGroups []string,
	timetoken int64,
	region string,
) (chan subscribeEnvelope, chan error) {
	resChan := make(chan subscribeEnvelope)
	errChan := make(chan error)

	go func() {
		defer close(errChan)
		defer close(resChan)
		opts := &subscribeOpts{
			pubnub:           pn,
			Channels:         channels,
			ChannelGroups:    channelGroups,
			Timetoken:        timetoken,
			Region:           region,
			Heartbeat:        pn.Config.PresenceTimeout,
			FilterExpression: pn.Config.FilterExpression,
			ctx:              ctx,
			QueryParam:       make(map[string]string),
		}

		pn.Config.Log.Println("before req")

		res, _, err := executeRequest(opts)
		pn.Config.Log.Println("before send")
		var envelope subscribeEnvelope
		if err == nil {
			err = json.Unmarshal(res, &envelope)
		}
		if err != nil {
			select {
			case <-ctx.Done():
			case errChan <- err:
			}
		} else {
			select {
			case <-ctx.Done():
			case resChan <- envelope:
			}
		}
		pn.Config.Log.Println("After everything")
	}()

	return resChan, errChan
}

func (pn *PubNub) iAmHere(
	ctx Context,
	channels []string,
	channelGroups []string,
) (chan interface{}, chan error) {
	resChan := make(chan interface{})
	errChan := make(chan error)

	go func() {
		c := newHeartbeatBuilderWithContext(pn, ctx)
		c.Channels(channels)
		c.ChannelGroups(channelGroups)

		res, _, err := c.Execute()

		if err != nil {
			select {
			case <-ctx.Done():
			case errChan <- err:
			}
		} else {
			select {
			case <-ctx.Done():
			case resChan <- res:
			}
		}
		close(resChan)
		close(errChan)

	}()

	return resChan, errChan
}

func (pn *PubNub) iAmAway(
	ctx Context,
	channels []string,
	channelGroups []string,
) (chan interface{}, chan error) {
	resChan := make(chan interface{})
	errChan := make(chan error)

	go func() {
		c := newLeaveBuilderWithContext(pn, ctx)
		c.Channels(channels)
		c.ChannelGroups(channelGroups)

		_, err := c.Execute()
		if err != nil {
			errChan <- err
		} else {
			resChan <- ""
		}
		close(resChan)
		close(errChan)

	}()

	return resChan, errChan
}

func (pn *PubNub) setPresenceState(
	ctx Context,
	channels []string,
	channelGroups []string,
	state map[string]interface{},
) (chan *SetStateResponse, chan error) {
	errChan := make(chan error)
	resChan := make(chan *SetStateResponse)

	go func() {
		c := newSetStateBuilderWithContext(pn, ctx)
		c.Channels(channels)
		c.ChannelGroups(channelGroups)
		c.State(state)

		res, _, err := c.Execute()
		if err != nil {
			errChan <- err
		} else {
			resChan <- res
		}
		close(errChan)
		close(resChan)

	}()

	return resChan, errChan
}
