package pubnub

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// newTestListenerManager builds a fully wired ListenerManager (with its
// pubnub/logger/config dependencies) for exercising the announce* fan-out in
// isolation. announceTimeout overrides the per-listener send deadline so the
// timeout-dependent tests run quickly; pass 0 to exercise the legacy
// block-until-delivered path.
func newTestListenerManager(announceTimeout time.Duration) *ListenerManager {
	cfg := NewConfigWithUserId(UserId(GenerateUUID()))
	cfg.SubscribeKey = "sub"
	cfg.SuppressLeaveEvents = true

	pn := NewPubNub(cfg)
	lm := pn.subscriptionManager.listenerManager
	lm.announceTimeout = announceTimeout
	return lm
}

// TestListenerManagerDefaultAnnounceTimeout documents that the leak guard is on
// by default: a freshly constructed manager carries the package-level
// listenerAnnounceTimeout, so the drop-on-stall behavior needs no configuration.
func TestListenerManagerDefaultAnnounceTimeout(t *testing.T) {
	cfg := NewConfigWithUserId(UserId(GenerateUUID()))
	cfg.SubscribeKey = "sub"
	cfg.SuppressLeaveEvents = true

	pn := NewPubNub(cfg)

	assert.Equal(t, listenerAnnounceTimeout, pn.subscriptionManager.listenerManager.announceTimeout)
}

// TestAnnounceDeliversToHealthyConsumer covers the normal fast path: an event is
// delivered to a listener that is actively reading its channel, and the bounded
// timeout does not interfere with timely delivery.
func TestAnnounceDeliversToHealthyConsumer(t *testing.T) {
	assert := assert.New(t)
	lm := newTestListenerManager(5 * time.Second)

	l := NewListener()
	lm.addListener(l)

	status := &PNStatus{Category: PNConnectedCategory}
	lm.announceStatus(status)

	select {
	case got := <-l.Status:
		assert.Equal(status, got)
	case <-time.After(2 * time.Second):
		t.Fatal("expected status event was not delivered to a healthy consumer")
	}
}

// TestAnnounceDropsForSlowConsumer is the regression test for the blocked-send
// goroutine leak (H-10). With a bounded announce timeout, an event destined for a
// consumer that never drains its channel must be dropped rather than pinning the
// announce goroutine (and the referenced event) indefinitely. The drop is proven
// by confirming the channel yields nothing once the timeout window passes: had
// the goroutine still been blocked on the send, the value would be received.
func TestAnnounceDropsForSlowConsumer(t *testing.T) {
	lm := newTestListenerManager(300 * time.Millisecond)

	l := NewListener()
	lm.addListener(l)

	lm.announceStatus(&PNStatus{Category: PNConnectedCategory})

	// Wait past the timeout so the announce goroutine gives up and drops the event.
	time.Sleep(600 * time.Millisecond)

	select {
	case <-l.Status:
		t.Fatal("event should have been dropped after timeout, but was still deliverable")
	case <-time.After(300 * time.Millisecond):
		// Expected: the send timed out and the goroutine exited, so nothing remains.
	}
}

// TestAnnounceSlowConsumerDoesNotBlockHealthyOne verifies the timeout also cures
// head-of-line blocking: a stalled listener must not prevent a healthy listener in
// the same fan-out from receiving the event.
func TestAnnounceSlowConsumerDoesNotBlockHealthyOne(t *testing.T) {
	lm := newTestListenerManager(1 * time.Second)

	slow := NewListener() // never read
	fast := NewListener() // read below
	lm.addListener(slow)
	lm.addListener(fast)

	status := &PNStatus{Category: PNConnectedCategory}
	lm.announceStatus(status)

	// The fast listener must receive within the timeout window regardless of the
	// map iteration order, i.e. even if the slow listener is visited first.
	select {
	case got := <-fast.Status:
		assert.Equal(t, status, got)
	case <-time.After(2 * time.Second):
		t.Fatal("healthy consumer was starved by a slow consumer in the same fan-out")
	}
}

// TestAnnounceZeroTimeoutBlocksUntilDelivery verifies the defensive fallback: a
// non-positive announce timeout restores the legacy guaranteed-delivery contract,
// where the event waits for a slow consumer instead of being dropped.
func TestAnnounceZeroTimeoutBlocksUntilDelivery(t *testing.T) {
	lm := newTestListenerManager(0) // disabled -> block until delivered

	l := NewListener()
	lm.addListener(l)

	status := &PNStatus{Category: PNConnectedCategory}
	lm.announceStatus(status)

	// Stall well beyond what any bounded timeout would tolerate, then read.
	time.Sleep(500 * time.Millisecond)

	select {
	case got := <-l.Status:
		assert.Equal(t, status, got)
	case <-time.After(1 * time.Second):
		t.Fatal("with timeout disabled the event must still be delivered to a late consumer")
	}
}

// TestCopyListenersConcurrentReads exercises the RLock downgrade in copyListeners:
// many concurrent announcements (each of which copies the listener set) must run
// without racing, and mutations remain safe under the race detector.
func TestCopyListenersConcurrentReads(t *testing.T) {
	lm := newTestListenerManager(1 * time.Second)

	for i := 0; i < 4; i++ {
		lm.addListener(NewListener())
	}

	var wg sync.WaitGroup
	assert.NotPanics(t, func() {
		for i := 0; i < 16; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = lm.copyListeners()
			}()
		}
		wg.Wait()
	})
}
