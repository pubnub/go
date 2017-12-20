// +build !go1.7

package pubnub

import (
	"net/http"
)

func setRequestContext(r *http.Request, ctx Context) *http.Request {
	r.Cancel = ctx.Done()

	return r
}
