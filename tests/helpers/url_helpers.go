package helpers

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Custom URL matcher for outgoing pubnub server requests
func UrlsEqual(expectedString, actualString string,
	ignoreKeys, mixedKeys []string) (bool, error) {
	expected, err := url.Parse(expectedString)
	if err != nil {
		return false, err
	}

	actual, err := url.Parse(actualString)
	if err != nil {
		return false, err
	}

	if expected.Scheme != actual.Scheme {
		return false, nil
	}

	if expected.Host != actual.Host {
		return false, nil
	}

	if !PathsEqual(expected.Path, actual.Path, []int{}) {
		return false, nil
	}

	eQuery := expected.Query()
	aQuery := actual.Query()

	if !QueriesEqual(&eQuery, &aQuery, ignoreKeys, mixedKeys) {
		return false, nil
	}

	return true, nil
}

// PathsEqual mixedPositions - a position of items which can contain unsorted items like
// multiple unsorted channels. If no such positions expected use an empty slice.
// Like in arrays, the first position is 0.
func PathsEqual(expectedString, actualString string,
	mixedPositions []int) bool {

	if expectedString == actualString {
		return true
	}

	expected := strings.Split(expectedString, "/")
	actual := strings.Split(actualString, "/")

	if len(actual) != len(expected) {
		return false
	}

	for k, v := range expected {
		if !isValueInSlice(k, mixedPositions) {
			if v != actual[k] {
				return false
			}
		} else {
			expectedItems := strings.Split(v, ",")
			actualItems := strings.Split(actual[k], ",")

			if len(expectedItems) != len(actualItems) {
				return false
			}

			for _, v := range expectedItems {
				if !isValueInSlice(v, expectedItems) {
					return false
				}
			}
		}
	}

	return true
}

func QueriesEqual(expectedString, actualString *url.Values,
	ignoreKeys []string, mixedKeys []string) bool {

	if expectedString.Encode() == actualString.Encode() {
		return true
	}

	for k, aVal := range *actualString {
		if isValueInSlice(k, ignoreKeys) {
			continue
		}

		if eVal, ok := (*expectedString)[k]; ok {
			if isValueInSlice(k, mixedKeys) {
				eParts := strings.Split(eVal[0], ",")
				aParts := strings.Split(aVal[0], ",")

				if len(aParts) != len(eParts) {
					return false
				}

				for _, e := range eParts {
					if !isValueInSlice(e, aParts) {
						return false
					}
				}
			} else {
				if aVal[0] != eVal[0] {
					return false
				}
			}
		} else {
			return false
		}
	}

	for k, _ := range *expectedString {
		if val := actualString.Get(k); val == "" {
			return false
		}
	}

	return true
}

func isValueInSlice(item interface{}, slice interface{}) bool {
	if s, ok := slice.([]string); ok {
		for _, v := range s {
			if item == v {
				return true
			}
		}
	} else if s, ok := slice.([]int); ok {
		for _, v := range s {
			if item == v {
				return true
			}
		}
	}

	return false
}

// Assertion wrappers for tests
func AssertPathsEqual(t *testing.T, expectedString, actualString string,
	itemsPositions []int) {
	match := PathsEqual(expectedString, actualString, itemsPositions)

	assert.True(t, match, "Paths are not equal:\nExpected: %s\nActual:   %s\n",
		expectedString, actualString)
}

func AssertQueriesEqual(t *testing.T, expectedString, actualString *url.Values,
	ignoreKeys, mixedKeys []string) {
	match := QueriesEqual(expectedString, actualString, ignoreKeys, mixedKeys)

	assert.True(t, match, "Queries are not equal:\nExpected: %s\nActual: %s\n",
		expectedString, actualString)
}
