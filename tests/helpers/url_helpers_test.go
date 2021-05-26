package helpers

import (
	"net/url"
	"testing"

	"github.com/pubnub/go/v5/utils"
	"github.com/stretchr/testify/assert"
)

func TestUrlsEqual(t *testing.T) {
	expected := "https://ps.pndns.com/publish/pub_key/sub_key/0/ch/0/\"hey\"?pnsdk=123&uuid=firstUUID"
	actual := "https://ps.pndns.com/publish/pub_key/sub_key/0/ch/0/\"hey\"?pnsdk=123&uuid=anotherUUID"

	match, err := UrlsEqual(expected, actual, []string{"uuid"}, []string{})

	if err != nil {
		assert.Fail(t, err.Error())
	}

	assert.True(t, match)
}

func TestUrlsNotEqual(t *testing.T) {
	expected := "https://ps.pndns.com/publish/pub_key/sub_key/0/ch/0/\"hey\"?pnsdk=123&uuid=firstUUID"
	actual := "https://ps.pndns.com/publish/pub_key/sub_key/0/ch/0/\"hey\"?pnsdk=123&uuid=anotherUUID"

	match, err := UrlsEqual(expected, actual, []string{}, []string{})

	if err != nil {
		assert.Fail(t, err.Error())
	}

	assert.False(t, match)
}

func TestSimplePathsEqual(t *testing.T) {
	expected := "one/two/three"
	actual := "one/two/three"

	assert.True(t, PathsEqual(expected, actual, []int{}))
}

func TestComplexPathsEqual(t *testing.T) {
	expected := "one/foo,bar,blah/three"
	actual := "one/foo,bar,blah/three"

	assert.True(t, PathsEqual(expected, actual, []int{1}))
}

func TestMixedPathsEqual(t *testing.T) {
	expected := "one/bar,foo,blah/three"
	actual := "one/foo,bar,blah/three"

	assert.True(t, PathsEqual(expected, actual, []int{1}))
}

func TestQueriesSameSizeEqual(t *testing.T) {
	expected := &url.Values{}
	expected.Set("channel", "my_ch")
	expected.Set("uuid", utils.UUID())

	actual := &url.Values{}
	actual.Set("channel", "my_ch")
	actual.Set("uuid", utils.UUID())

	assert.True(t, QueriesEqual(expected, actual, []string{"uuid"}, []string{}))
}

func TestQueriesDifferentSizeEqual(t *testing.T) {
	expected := &url.Values{}
	expected.Set("channel", "my_ch")
	expected.Set("uuid", utils.UUID())

	actual := &url.Values{}
	actual.Set("channel", "my_ch")
	actual.Set("group", "my_gr")
	actual.Set("uuid", utils.UUID())

	assert.True(t, QueriesEqual(expected, actual, []string{"uuid", "group"}, []string{}))
}

func TestQueriesDifferentSizeNotEqual(t *testing.T) {
	expected := &url.Values{}
	expected.Set("channel", "my_ch")
	expected.Set("group", "my_gr")
	expected.Set("uuid", utils.UUID())

	actual := &url.Values{}
	actual.Set("channel", "my_ch")
	actual.Set("uuid", utils.UUID())

	assert.False(t, QueriesEqual(expected, actual, []string{}, []string{}))
}

func TestQueriesSameSizeNotEqual(t *testing.T) {
	expected := &url.Values{}
	expected.Set("channel", "my_ch")
	expected.Set("uuid", utils.UUID())

	actual := &url.Values{}
	actual.Set("channel", "my_ch")
	actual.Set("uuid", utils.UUID())

	assert.False(t, QueriesEqual(expected, actual, []string{}, []string{}))
}

func TestMixedQueriesEqual(t *testing.T) {
	expected := &url.Values{}
	expected.Set("channel", "ch1,ch2,ch3")
	expected.Set("uuid", utils.UUID())

	actual := &url.Values{}
	actual.Set("channel", "ch2,ch1,ch3")
	actual.Set("uuid", utils.UUID())

	assert.True(t, QueriesEqual(expected, actual, []string{"uuid"},
		[]string{"channel"}))
}

func TestMixedQueriesNotEqual(t *testing.T) {
	expected := &url.Values{}
	expected.Set("heartbeat", "300")
	expected.Set("hey", "123")

	actual := &url.Values{}
	actual.Set("heartbeat", "300")
	actual.Set("pnsdk", utils.UUID())
	actual.Set("uuid", utils.UUID())
	assert.False(t, QueriesEqual(expected, actual, []string{"pnsdk", "uuid", "tt"},
		[]string{}))
}
