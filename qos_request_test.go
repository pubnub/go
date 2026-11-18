package pubnub

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const qosMockBody = `[1,"Sent","1"]`

type qosSeen struct {
	qos         string
	hasDeadline bool
	remaining   time.Duration
}

// qosCaptureTransport records the deadline http.Client applied before RoundTrip.
// It does not open a listening server.
type qosCaptureTransport struct {
	mu          sync.Mutex
	seen        []qosSeen
	entered     chan struct{}
	enteredOnce sync.Once
	release     chan struct{}
	waitCancel  bool
}

func (t *qosCaptureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	deadline, ok := req.Context().Deadline()
	remaining := time.Duration(0)
	if ok {
		remaining = time.Until(deadline)
	}
	t.mu.Lock()
	t.seen = append(t.seen, qosSeen{qos: req.URL.Query().Get("qos"), hasDeadline: ok, remaining: remaining})
	t.mu.Unlock()

	if t.entered != nil {
		t.enteredOnce.Do(func() { close(t.entered) })
	}
	if t.waitCancel {
		<-req.Context().Done()
		return nil, req.Context().Err()
	}
	if t.release != nil && req.URL.Query().Get("qos") == "1" {
		select {
		case <-t.release:
		case <-req.Context().Done():
			return nil, req.Context().Err()
		}
	}

	return &http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(qosMockBody)),
		Request:    req,
	}, nil
}

func (t *qosCaptureTransport) snapshot() []qosSeen {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]qosSeen, len(t.seen))
	copy(out, t.seen)
	return out
}

func pubnubWithPinnedTimeout(timeout time.Duration, tr http.RoundTripper) (*PubNub, *http.Client) {
	pn := NewPubNub(NewDemoConfig())
	client := &http.Client{Transport: tr, Timeout: timeout}
	pn.SetClient(client)
	return pn, client
}

func assertRemainingNear(t *testing.T, got, want time.Duration) {
	t.Helper()
	const slack = 8 * time.Second
	if got > want+time.Second || got < want-slack {
		t.Fatalf("deadline remaining = %s, want near %s", got, want)
	}
}

func TestClientForQoSTimeoutFloor(t *testing.T) {
	transport := &http.Transport{}
	redirect := func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }

	short := &http.Client{Transport: transport, Timeout: 10 * time.Second, CheckRedirect: redirect}
	raised := clientForQoS(short)
	require.NotSame(t, short, raised)
	assert.Equal(t, qosMinRequestTimeout, raised.Timeout)
	assert.Same(t, transport, raised.Transport)
	assert.NotNil(t, raised.CheckRedirect)
	assert.Equal(t, 10*time.Second, short.Timeout)

	exact := &http.Client{Transport: transport, Timeout: qosMinRequestTimeout}
	assert.Same(t, exact, clientForQoS(exact))

	longer := &http.Client{Transport: transport, Timeout: 45 * time.Second}
	assert.Same(t, longer, clientForQoS(longer))
	assert.Equal(t, 45*time.Second, longer.Timeout)

	unlimited := &http.Client{Transport: transport, Timeout: 0}
	assert.Same(t, unlimited, clientForQoS(unlimited))

	assert.Equal(t, qosMinRequestTimeout, clientForQoS(nil).Timeout)
}

func TestPublishQoSQuery(t *testing.T) {
	pn := NewPubNub(NewDemoConfig())

	enabled, err := newPublishBuilder(pn).Channel("ch").Message("hey").QoS(true).opts.buildQuery()
	require.NoError(t, err)
	assert.Equal(t, "1", enabled.Get("qos"))

	disabled, err := newPublishBuilder(pn).Channel("ch").Message("hey").QoS(false).opts.buildQuery()
	require.NoError(t, err)
	assert.Empty(t, disabled.Get("qos"))

	omitted, err := newPublishBuilder(pn).Channel("ch").Message("hey").opts.buildQuery()
	require.NoError(t, err)
	assert.Empty(t, omitted.Get("qos"))

	post, err := newPublishBuilder(pn).Channel("ch").Message("hey").UsePost(true).QoS(true).opts.buildQuery()
	require.NoError(t, err)
	assert.Equal(t, "1", post.Get("qos"))

	overridden, err := newPublishBuilder(pn).Channel("ch").Message("hey").
		QueryParam(map[string]string{"qos": "0"}).QoS(true).opts.buildQuery()
	require.NoError(t, err)
	assert.Equal(t, "1", overridden.Get("qos"))
}

func TestSignalQoSQuery(t *testing.T) {
	pn := NewPubNub(NewDemoConfig())

	enabled, err := newSignalBuilder(pn).Channel("ch").Message("hey").QoS(true).opts.buildQuery()
	require.NoError(t, err)
	assert.Equal(t, "1", enabled.Get("qos"))

	disabled, err := newSignalBuilder(pn).Channel("ch").Message("hey").QoS(false).opts.buildQuery()
	require.NoError(t, err)
	assert.Empty(t, disabled.Get("qos"))

	omitted, err := newSignalBuilder(pn).Channel("ch").Message("hey").opts.buildQuery()
	require.NoError(t, err)
	assert.Empty(t, omitted.Get("qos"))
}

func TestQoSExecuteRaisesShortClientTimeout(t *testing.T) {
	tests := []struct {
		name string
		call func(*PubNub) error
	}{
		{
			name: "publish",
			call: func(pn *PubNub) error {
				_, _, err := pn.Publish().Channel("ch").Message("hey").QoS(true).Execute()
				return err
			},
		},
		{
			name: "publish post",
			call: func(pn *PubNub) error {
				_, _, err := pn.Publish().Channel("ch").Message("hey").UsePost(true).QoS(true).Execute()
				return err
			},
		},
		{
			name: "signal",
			call: func(pn *PubNub) error {
				_, _, err := pn.Signal().Channel("ch").Message("hey").QoS(true).Execute()
				return err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transport := &qosCaptureTransport{}
			cfg := NewDemoConfig()
			cfg.NonSubscribeRequestTimeout = 60
			pn := NewPubNub(cfg)
			pinned := &http.Client{Transport: transport, Timeout: 10 * time.Second}
			pn.SetClient(pinned)

			require.NoError(t, test.call(pn))

			seen := transport.snapshot()
			require.Len(t, seen, 1)
			assert.Equal(t, "1", seen[0].qos)
			assert.True(t, seen[0].hasDeadline)
			assertRemainingNear(t, seen[0].remaining, qosMinRequestTimeout)
			assert.Equal(t, 10*time.Second, pn.GetClient().Timeout)
			assert.Same(t, pinned, pn.GetClient())
			assert.Same(t, transport, pn.GetClient().Transport)
		})
	}
}

func TestQoSExecuteKeepsLongerAndUnlimitedClientTimeouts(t *testing.T) {
	transport := &qosCaptureTransport{}
	pn, pinned := pubnubWithPinnedTimeout(45*time.Second, transport)

	_, _, err := pn.Publish().Channel("ch").Message("hey").QoS(true).Execute()
	require.NoError(t, err)

	seen := transport.snapshot()
	require.Len(t, seen, 1)
	assert.True(t, seen[0].hasDeadline)
	assertRemainingNear(t, seen[0].remaining, 45*time.Second)
	assert.Equal(t, 45*time.Second, pinned.Timeout)
	assert.Same(t, pinned, pn.GetClient())

	unlimitedTransport := &qosCaptureTransport{}
	unlimitedPN, unlimited := pubnubWithPinnedTimeout(0, unlimitedTransport)
	_, _, err = unlimitedPN.Signal().Channel("ch").Message("hey").QoS(true).Execute()
	require.NoError(t, err)
	seen = unlimitedTransport.snapshot()
	require.Len(t, seen, 1)
	assert.Equal(t, "1", seen[0].qos)
	assert.False(t, seen[0].hasDeadline)
	assert.Equal(t, time.Duration(0), unlimited.Timeout)
}

func TestQoSPublishDoesNotChangeConcurrentRequestDeadline(t *testing.T) {
	release := make(chan struct{})
	var releaseOnce sync.Once
	closeRelease := func() { releaseOnce.Do(func() { close(release) }) }
	t.Cleanup(closeRelease)

	transport := &qosCaptureTransport{entered: make(chan struct{}), release: release}
	pn, pinned := pubnubWithPinnedTimeout(10*time.Second, transport)

	done := make(chan error, 1)
	go func() {
		_, _, err := pn.Publish().Channel("ch").Message("qos").QoS(true).Execute()
		done <- err
	}()

	select {
	case <-transport.entered:
	case <-time.After(5 * time.Second):
		t.Fatal("QoS publish did not reach the transport")
	}

	assert.Equal(t, 10*time.Second, pinned.Timeout)

	_, _, err := pn.Publish().Channel("ch").Message("plain").Execute()
	require.NoError(t, err)

	seen := transport.snapshot()
	require.Len(t, seen, 2)
	assert.Equal(t, "1", seen[0].qos)
	assertRemainingNear(t, seen[0].remaining, qosMinRequestTimeout)
	assert.Empty(t, seen[1].qos)
	assert.True(t, seen[1].hasDeadline)
	assertRemainingNear(t, seen[1].remaining, 10*time.Second)
	assert.Equal(t, 10*time.Second, pn.GetClient().Timeout)

	closeRelease()
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("QoS publish did not finish")
	}
	assert.Equal(t, 10*time.Second, pn.GetClient().Timeout)
}

func TestQoSContextCancelStopsLocalWait(t *testing.T) {
	tests := []struct {
		name string
		call func(*PubNub, context.Context) error
	}{
		{
			name: "publish",
			call: func(pn *PubNub, ctx context.Context) error {
				_, _, err := pn.PublishWithContext(ctx).Channel("ch").Message("hey").QoS(true).Execute()
				return err
			},
		},
		{
			name: "signal",
			call: func(pn *PubNub, ctx context.Context) error {
				_, _, err := pn.SignalWithContext(ctx).Channel("ch").Message("hey").QoS(true).Execute()
				return err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transport := &qosCaptureTransport{entered: make(chan struct{}), waitCancel: true}
			pn, pinned := pubnubWithPinnedTimeout(10*time.Second, transport)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			done := make(chan error, 1)
			go func() {
				done <- test.call(pn, ctx)
			}()

			select {
			case <-transport.entered:
			case <-time.After(5 * time.Second):
				t.Fatal("request did not reach the transport")
			}
			cancel()

			select {
			case err := <-done:
				require.Error(t, err)
				assert.Contains(t, err.Error(), "context canceled")
			case <-time.After(5 * time.Second):
				t.Fatal("canceled QoS call did not return")
			}
			assert.Len(t, transport.snapshot(), 1)
			assert.Equal(t, 10*time.Second, pinned.Timeout)
		})
	}
}

func TestQoSCallerDeadlineRemainsShorterThanFloor(t *testing.T) {
	transport := &qosCaptureTransport{waitCancel: true}
	pn, _ := pubnubWithPinnedTimeout(10*time.Second, transport)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	started := time.Now()
	_, _, err := pn.PublishWithContext(ctx).Channel("ch").Message("hey").QoS(true).Execute()
	elapsed := time.Since(started)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
	assert.Less(t, elapsed, 5*time.Second)

	seen := transport.snapshot()
	require.Len(t, seen, 1)
	assert.Equal(t, "1", seen[0].qos)
	assert.True(t, seen[0].hasDeadline)
	assert.Less(t, seen[0].remaining, 2*time.Second)
}
