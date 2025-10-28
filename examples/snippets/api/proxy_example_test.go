// snippet.includes
// Replace with your package name (usually "main")
package pubnub_samples_test

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	pubnub "github.com/pubnub/go/v7"
)

// snippet.end

/*
IMPORTANT NOTE FOR COPYING EXAMPLES:

Throughout this file, you'll see code between "snippet.hide" and "snippet.show" comments.
These sections are used for CI/CD testing and should be SKIPPED if you're copying examples.

Example of what to skip:
	// snippet.hide
	config = setPubnubExampleConfigData(config)  // <- Skip this line (for testing only)
	defer pn.DeleteChannelGroup().Execute()      // <- Skip this line (cleanup for tests)
	// snippet.show

When copying examples to your own code:
- Use your own publish/subscribe keys instead of the "demo" keys
- Remove any statements that are between snippet.hide and snippet.show (they're only for testing purposes)
*/

// snippet.proxy_subscribe
// proxySubscribe demonstrates configuring a proxy for subscribe requests
func proxySubscribe() {
	// Configure a proxy specifically for subscribe requests.
	var pn *pubnub.PubNub
	config := pubnub.NewConfigWithUserId(pubnub.UserId("myUniqueUserId"))
	config.UseHTTP2 = false

	pn = pubnub.NewPubNub(config)

	transport := &http.Transport{
		MaxIdleConnsPerHost: pn.Config.MaxIdleConnsPerHost,
		Dial: (&net.Dialer{
			Timeout:   time.Duration(pn.Config.ConnectTimeout) * time.Second,
			KeepAlive: 30 * time.Minute,
		}).Dial,
		ResponseHeaderTimeout: time.Duration(pn.Config.SubscribeRequestTimeout) * time.Second,
	}
	proxyURL, err := url.Parse(fmt.Sprintf("http://%s:%s@%s:%d", "proxyUser", "proxyPassword", "proxyServer", 8080))

	if err == nil {
		transport.Proxy = http.ProxyURL(proxyURL)
	} else {
		fmt.Printf("ERROR: creatSubHTTPClient: Proxy connection error: %s", err.Error())
	}
	c := pn.GetSubscribeClient()
	c.Transport = transport
	pn.SetSubscribeClient(c)
}

// snippet.proxy_non_subscribe
// proxyNonSubscribe demonstrates configuring a proxy for non-subscribe requests
func proxyNonSubscribe() {
	// Configure a proxy for non-subscribe requests (publish, history, etc.).
	var pn *pubnub.PubNub
	config := pubnub.NewConfigWithUserId(pubnub.UserId("myUniqueUserId"))
	config.UseHTTP2 = false

	pn = pubnub.NewPubNub(config)

	transport := &http.Transport{
		MaxIdleConnsPerHost: pn.Config.MaxIdleConnsPerHost,
		Dial: (&net.Dialer{
			Timeout:   time.Duration(pn.Config.ConnectTimeout) * time.Second,
			KeepAlive: 30 * time.Minute,
		}).Dial,
		ResponseHeaderTimeout: time.Duration(pn.Config.NonSubscribeRequestTimeout) * time.Second,
	}
	proxyURL, err := url.Parse(fmt.Sprintf("http://%s:%s@%s:%d", "proxyUser", "proxyPassword", "proxyServer", 8080))

	if err == nil {
		transport.Proxy = http.ProxyURL(proxyURL)
	} else {
		fmt.Printf("ERROR: createNonSubHTTPClient: Proxy connection error: %s", err.Error())
	}
	c := pn.GetClient()
	c.Transport = transport
	pn.SetClient(c)
}

// snippet.end
