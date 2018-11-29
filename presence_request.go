package pubnub

import (
	"strings"
)

type presenceBuilder struct {
	opts *presenceOpts
}

type presenceOpts struct {
	pubnub *PubNub

	Channels      []string
	ChannelGroups []string
	Connected     bool
	ctx           Context
}

func newPresenceBuilder(pubnub *PubNub) *presenceBuilder {
	builder := presenceBuilder{
		opts: &presenceOpts{
			pubnub: pubnub,
		},
	}

	return &builder
}

func newPresenceBuilderWithContext(pubnub *PubNub, context Context) *presenceBuilder {
	builder := presenceBuilder{
		opts: &presenceOpts{
			pubnub: pubnub,
			ctx:    context,
		},
	}

	return &builder
}

// Channels sets the Channels for the Presence request.
func (b *presenceBuilder) Channels(ch []string) *presenceBuilder {
	b.opts.Channels = ch

	return b
}

// ChannelGroups sets the ChannelGroups for the Presence request.
func (b *presenceBuilder) ChannelGroups(cg []string) *presenceBuilder {
	b.opts.ChannelGroups = cg

	return b
}

// Channels sets the Channels for the Presence request.
func (b *presenceBuilder) Connected(connected bool) *presenceBuilder {
	b.opts.Connected = connected

	return b
}

// func (o *presenceOpts) buildPath() (string, error) {
// 	channels := string(utils.JoinChannels(o.Channels))

// 	return fmt.Sprintf(heartbeatPath,
// 		o.pubnub.Config.SubscribeKey,
// 		channels), nil
// }

// func (o *presenceOpts) buildQuery() (*url.Values, error) {
// 	q := defaultQuery(o.pubnub.Config.UUID, o.pubnub.telemetryManager)

// 	q.Set("heartbeat", strconv.Itoa(o.pubnub.Config.PresenceTimeout))

// 	if len(o.ChannelGroups) > 0 {
// 		q.Set("channel-group", strings.Join(o.ChannelGroups, ","))
// 	}

// 	if o.State != nil {
// 		state, err := utils.ValueAsString(o.State)
// 		if err != nil {
// 			return &url.Values{}, err
// 		}

// 		if string(state) != "{}" {
// 			q.Set("state", string(state))
// 		}
// 	}

// 	return q, nil
// }

// func (o *presenceOpts) validate() error {
// 	if o.config().SubscribeKey == "" {
// 		return newValidationError(o, StrMissingSubKey)
// 	}

// 	if len(o.Channels) == 0 && len(o.ChannelGroups) == 0 {
// 		return newValidationError(o, "Missing Channel or Channel Group")
// 	}

// 	return nil
// }

func (b *presenceBuilder) Execute() {
	for _, ch := range b.opts.Channels {
		if strings.Contains(ch, "-pnpres") {
			ch = strings.Replace(ch, "-pnpres", "", -1)
		}
		b.opts.pubnub.heartbeatManager.heartbeatChannels[ch] = newSubscriptionItem(ch)
	}
	for _, cg := range b.opts.ChannelGroups {
		if strings.Contains(cg, "-pnpres") {
			cg = strings.Replace(cg, "-pnpres", "", -1)
		}
		b.opts.pubnub.heartbeatManager.heartbeatGroups[cg] = newSubscriptionItem(cg)
	}

	b.opts.pubnub.heartbeatManager.startHeartbeatTimer()
}

// func (o *presenceOpts) config() Config {
// 	return *o.pubnub.Config
// }

// func (o *presenceOpts) client() *http.Client {
// 	return o.pubnub.GetClient()
// }

// func (o *presenceOpts) context() Context {
// 	return o.ctx
// }

// func (o *presenceOpts) jobQueue() chan *JobQItem {
// 	return o.pubnub.jobQueue
// }

// func (o *presenceOpts) buildBody() ([]byte, error) {
// 	return []byte{}, nil
// }

// func (o *presenceOpts) httpMethod() string {
// 	return "GET"
// }

// func (o *presenceOpts) isAuthRequired() bool {
// 	return true
// }

// func (o *presenceOpts) requestTimeout() int {
// 	return o.pubnub.Config.NonSubscribeRequestTimeout
// }

// func (o *presenceOpts) connectTimeout() int {
// 	return o.pubnub.Config.ConnectTimeout
// }

// func (o *presenceOpts) operationType() OperationType {
// 	return PNHeartBeatOperation
// }

// func (o *presenceOpts) telemetryManager() *TelemetryManager {
// 	return o.pubnub.telemetryManager
// }

//starthb
//Stop hb
//perform hb loop

//execute

//read channel from sub
//when the hb runs on subscribe, stop this.

//store channels and channel groups in statemanager
