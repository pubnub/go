package pubnub

import (
	"fmt"
	"net/url"
	"testing"

	h "github.com/pubnub/go/v7/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestHistoryDeleteRequestAllParams(t *testing.T) {
	assert := assert.New(t)

	opts := optsWithDefaultTestValues(pubnub)

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v3/history/sub-key/sub_key/channel/%s", opts.Channel),
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("start", "123")
	expected.Set("end", "456")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestHistoryDeleteRequestQueryParams(t *testing.T) {
	assert := assert.New(t)

	opts := optsWithDefaultTestValues(pubnub)

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts.QueryParam = queryParam

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("start", "123")
	expected.Set("end", "456")
	expected.Set("q1", "v1")
	expected.Set("q2", "v2")

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestNewHistoryDeleteBuilder(t *testing.T) {
	assert := assert.New(t)

	o := newHistoryDeleteBuilder(pubnub)
	o.Channel("ch")

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v3/history/sub-key/sub_key/channel/%s", o.opts.Channel),
		u.EscapedPath(), []int{})

	_, err1 := o.opts.buildQuery()
	assert.Nil(err1)

}

func TestNewHistoryDeleteBuilderContext(t *testing.T) {
	assert := assert.New(t)

	o := newHistoryDeleteBuilderWithContext(pubnub, backgroundContext)
	o.Channel("ch")

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v3/history/sub-key/sub_key/channel/%s", o.opts.Channel),
		u.EscapedPath(), []int{})

	_, err1 := o.opts.buildQuery()
	assert.Nil(err1)

}

func TestHistoryDeleteOptsValidateSub(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := optsWithDefaultTestValues(pn)

	assert.Equal("pubnub/validation: pubnub: Delete messages: Missing Subscribe Key", opts.validate().Error())
}

func TestHistoryDeleteOptsValidateSec(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SecretKey = ""
	opts := optsWithDefaultTestValues(pn)

	assert.Equal("pubnub/validation: pubnub: Delete messages: Missing Secret Key", opts.validate().Error())
}

func optsWithDefaultTestValues(pn *PubNub) *historyDeleteOpts {
	opts := newHistoryDeleteOpts(pn, pn.ctx)
	opts.Channel = "ch"
	opts.SetStart = true
	opts.SetEnd = true
	opts.Start = int64(123)
	opts.End = int64(456)
	return opts
}
