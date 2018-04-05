package pubnub

import (
	"fmt"
	"net/url"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertSuccessFireGet(t *testing.T, expectedString string, message interface{}) {
	assert := assert.New(t)

	opts := &fireOpts{
		Channel: "ch",
		Message: message,
		pubnub:  pubnub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/publish/pub_key/sub_key/0/ch/0/%s", expectedString),
		path, []int{})

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)
}

func AssertSuccessFireQuery(t *testing.T, expectedString string, message interface{}) {
	opts := &fireOpts{
		Channel: "ch",
		Message: "hey",
		pubnub:  pubnub,
	}

	query, _ := opts.buildQuery()

	expected := &url.Values{}
	expected.Set("store", "0")
	expected.Set("norep", "true")

	h.AssertQueriesEqual(t, expected, query,
		[]string{"seqn", "pnsdk", "uuid", "store", "norep"}, []string{})

}

func TestFirePath(t *testing.T) {
	message := "test"
	AssertSuccessFireGet(t, "%22test%22", message)
}

func TestFireQuery(t *testing.T) {
	message := "test"
	AssertSuccessFireQuery(t, "%22test%22?store=0&norep=true&", message)
}
