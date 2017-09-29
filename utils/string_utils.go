package utils

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"

	uuid "github.com/satori/go.uuid"
)

// TODO: return string
func JoinChannels(channels []string) []byte {
	if len(channels) == 0 {
		return []byte(",")
	}

	var encodedChannels []string

	for _, value := range channels {
		encodedChannels = append(encodedChannels, UrlEncode(value))
	}

	return []byte(strings.Join(encodedChannels, ","))
}

// PubNub - specific serializer
func ValueAsString(value interface{}) ([]byte, error) {
	switch t := value.(type) {
	case string:
		return []byte(fmt.Sprintf("\"%s\"", t)), nil
	default:
		val, err := json.Marshal(value)
		return val, err
	}
}

// Generate a random uuid string
func Uuid() string {
	return uuid.NewV4().String()
}

func sortQueries(params *url.Values) []string {
	sortedKeys := make([]string, len(*params))
	i := 0

	for k, _ := range *params {
		sortedKeys[i] = k
		i++
	}

	sort.Strings(sortedKeys)

	return sortedKeys
}

func PreparePamParams(params *url.Values) string {
	sortedKeys := sortQueries(params)
	stringifiedQuery := ""
	i := 0

	for _, v := range sortedKeys {
		if i == len(sortedKeys)-1 {
			stringifiedQuery += fmt.Sprintf("%s=%s", v, pamEncode((*params)[v][0]))
		} else {
			stringifiedQuery += fmt.Sprintf("%s=%s&", v, pamEncode((*params)[v][0]))
		}

		i++
	}

	return stringifiedQuery
}

func pamEncode(value string) string {
	// *!'()[]~
	stringifiedParam := UrlEncode(value)

	var replacer = strings.NewReplacer(
		"*", "%2A",
		"!", "%21",
		"'", "%27",
		"(", "%28",
		")", "%29",
		"[", "%5B",
		"]", "%5D",
		"~", "%7E")

	stringifiedParam = replacer.Replace(stringifiedParam)

	return stringifiedParam
}

func QueryToString(query *url.Values) string {
	stringifiedQuery := ""
	i := 0

	for k, v := range *query {
		if i == len(*query)-1 {
			stringifiedQuery += fmt.Sprintf("%s=%s", k, v[0])
		} else {
			stringifiedQuery += fmt.Sprintf("%s=%s&", k, v[0])
		}

		i++
	}

	return stringifiedQuery
}

// TODO: verify the helper is used where supposed to
func UrlEncode(s string) string {
	v := url.QueryEscape(s)

	var replacer = strings.NewReplacer(
		"+", "%20",
	)

	return replacer.Replace(v)
}
