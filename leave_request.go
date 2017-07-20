package pubnub

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pubnub/go/utils"
)

const LEAVE_PATH = "/v2/presence/sub-key/%s/channel/%s/leave"

func LeaveRequest(pn *PubNub, opts *LeaveOpts) error {
	opts.pubnub = pn
	_, err := executeRequest(opts)
	if err != nil {
		return err
	}

	return nil
}

type LeaveOpts struct {
	Channels      []string
	ChannelGroups []string

	pubnub *PubNub
	ctx    Context
}

func (o *LeaveOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *LeaveOpts) httpMethod() string {
	return "GET"
}

func (o *LeaveOpts) buildPath() (string, error) {
	channels, err := utils.ChannelsAsString(o.Channels)

	if err != nil {
		return "", err
	}

	if string(channels) == "" {
		channels = []byte(",")
	}

	return fmt.Sprintf(LEAVE_PATH,
		o.pubnub.Config.SubscribeKey,
		channels), nil
}

func (o *LeaveOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.pubnub.Config.Uuid)

	if len(o.ChannelGroups) > 0 {
		channelGroup, _ := utils.ChannelsAsString(o.ChannelGroups)
		q.Set("channel-group", string(channelGroup))
	}

	return q, nil
}

func (o *LeaveOpts) client() *http.Client {
	return o.pubnub.GetClient()
}

func (o *LeaveOpts) config() Config {
	return *o.pubnub.Config
}

func (o *LeaveOpts) context() Context {
	return o.ctx
}

func (o *LeaveOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return ErrMissingSubKey
	}

	if len(o.Channels) == 0 {
		return ErrMissingChannel
	}

	return nil
}
