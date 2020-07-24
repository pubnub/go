package pubnub

import (
	"fmt"
	"net/url"
	"testing"

	h "github.com/pubnub/go/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestHistoryDeleteRequestAllParams(t *testing.T) {
	assert := assert.New(t)

	opts := &historyDeleteOpts{
		Channel:  "ch",
		SetStart: true,
		SetEnd:   true,
		Start:    int64(123),
		End:      int64(456),
		pubnub:   pubnub,
	}

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

	opts := &historyDeleteOpts{
		Channel:  "ch",
		SetStart: true,
		SetEnd:   true,
		Start:    int64(123),
		End:      int64(456),
		pubnub:   pubnub,
	}

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
	opts := &historyDeleteOpts{
		Channel:  "ch",
		SetStart: true,
		SetEnd:   true,
		Start:    int64(123),
		End:      int64(456),
		pubnub:   pn,
	}

	assert.Equal("pubnub/validation: pubnub: Delete messages: Missing Subscribe Key", opts.validate().Error())
}

func TestHistoryDeleteOptsValidateSec(t *testing.T) {
	assert := assert.New(t)
	pn := NewPubNub(NewDemoConfig())
	pn.Config.SecretKey = ""
	opts := &historyDeleteOpts{
		Channel:  "ch",
		SetStart: true,
		SetEnd:   true,
		Start:    int64(123),
		End:      int64(456),
		pubnub:   pn,
	}

	assert.Equal("pubnub/validation: pubnub: Delete messages: Missing Secret Key", opts.validate().Error())
}
