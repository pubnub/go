package utils

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/anovikov1984/go-vcr/cassette"
)

var logMu sync.Mutex

func NewPubnubMatcher(skipFields []string) cassette.Matcher {
	return &PubnubMatcher{
		skipFields: skipFields,
	}
}

func NewPubnubSubscribeMatcher(skipFields []string) cassette.Matcher {
	return &PubnubMatcher{
		skipFields: skipFields,
	}
}

// Matcher for non-subscribe requests
type PubnubMatcher struct {
	cassette.Matcher

	isSubscribe bool
	skipFields  []string
}

func (m *PubnubMatcher) Match(interactions []*cassette.Interaction,
	r *http.Request) (*cassette.Interaction, error) {

interactionsLoop:
	for _, i := range interactions {
		if r.Method != i.Request.Method {
			continue
		}

		expectedURL, err := url.Parse(i.URL)
		if err != nil {
			continue
		}

		if expectedURL.Host != r.URL.Host {
			continue
		}

		if !m.matchPath(expectedURL.Path, r.URL.Path) {
			continue
		}
		eQuery := expectedURL.Query()
		aQuery := r.URL.Query()

		for fKey, _ := range eQuery {
			if hasKey(fKey, m.skipFields) {
				continue
			}

			if aQuery[fKey] == nil || eQuery.Get(fKey) != aQuery.Get(fKey) {
				continue interactionsLoop
			}
		}

		return i, nil
	}

	return nil, errorInteractionNotFound(interactions)
}

func (m *PubnubMatcher) matchPath(expected, actual string) bool {
	if isSubscribeRe.MatchString(expected) && isSubscribeRe.MatchString(actual) {
		return urlsMatch(expected, actual)
	} else {
		return expected == actual
	}
}

func errorInteractionNotFound(
	interactions []*cassette.Interaction) error {

	var urlsBuffer bytes.Buffer

	for _, i := range interactions {
		urlsBuffer.WriteString(i.URL)
		urlsBuffer.WriteString("\n")
	}

	return errors.New(fmt.Sprintf(
		"Interaction not found in:\n%s",
		urlsBuffer.String()))
}

var isSubscribeRe = regexp.MustCompile("^/subscribe/.*$")
var subscribePathRe = regexp.MustCompile("^((?:http|https)://[^/]+)?(/subscribe/[^/]+/)([^/]+)(/[^?]+)?.+$")

func urlsMatch(expected, actual string) bool {
	eAllMatches := subscribePathRe.FindAllStringSubmatch(expected, -4)
	aAllMatches := subscribePathRe.FindAllStringSubmatch(actual, -4)

	fmt.Println(eAllMatches[0][1], aAllMatches[0][1])

	if len(eAllMatches) > 0 && len(aAllMatches) > 0 {
		eMatches := eAllMatches[0][2:]
		aMatches := aAllMatches[0][2:]

		if eMatches[0] != aMatches[0] {
			return false
		}

		eChannels := strings.Split(eMatches[1], ",")
		aChannels := strings.Split(aMatches[1], ",")

		if !AssertStringSliceElementsEqual(eChannels, aChannels) {
			fmt.Println("chanels are NOT equal", eChannels, aChannels)
			return false
		} else {
			fmt.Println("chanels ARE equal", eChannels, aChannels)
		}

		if eMatches[2] != aMatches[2] {
			return false
		}

		eUrl, err := url.Parse(expected)
		if err != nil {
			panic(err.Error)
		}

		aUrl, err := url.Parse(actual)
		if err != nil {
			panic(err.Error)
		}

		eQuery := eUrl.Query()
		aQuery := aUrl.Query()

		for fKey, _ := range eQuery {
			if fKey == "channel-group" {
				if _, ok := aQuery["channel-group"]; ok {
					eCgs := eQuery.Get(fKey)
					aCgs := aQuery.Get(fKey)
					eChannels := strings.Split(eCgs, ",")
					aChannels := strings.Split(aCgs, ",")
					return AssertStringSliceElementsEqual(eChannels, aChannels)
				} else {
					return false
				}
			} else {
				if aQuery[fKey] == nil || eQuery.Get(fKey) != aQuery.Get(fKey) {
					return false
				}
			}
		}

		return true
	} else {
		return false
	}
}
