// +build go1.7

package pubnub

import (
	"fmt"
	"net/http"
)

func setRequestContext(r *http.Request, ctx Context) *http.Request {
	fmt.Println("new way")
	return r.WithContext(ctx)
}
