package pubnub

import (
	"fmt"
	"testing"

	h "github.com/pubnub/go/v5/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func AssertGetFileURL(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	o := newGetFileURLBuilder(pn)
	if testContext {
		o = newGetFileURLBuilderWithContext(pn, backgroundContext)
	}

	channel := "chan"
	id := "fileid"
	name := "filename"
	o.Channel(channel)
	o.QueryParam(queryParam)
	o.ID(id)
	o.Name(name)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf(getFileURLPath, pn.Config.SubscribeKey, channel, id, name),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)

	if checkQueryParam {
		u, _ := o.opts.buildQuery()
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
	}

}

func TestGetFileURL(t *testing.T) {
	AssertGetFileURL(t, true, false)
}

func TestGetFileURLContext(t *testing.T) {
	AssertGetFileURL(t, true, true)
}

func TestGetFileURLResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	channel := "chan"
	id := "fileid"
	name := "filename"

	r, _, err := pn.GetFileURL().ID(id).Name(name).Channel(channel).Execute()
	assert.Contains(r.URL, fmt.Sprintf("%s/files/%s/%s", channel, id, name))

	assert.Nil(err)
}
