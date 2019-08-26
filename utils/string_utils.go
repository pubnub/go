package utils

import (
	"encoding/json"
	//"log"

	//"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"

	pnerr "github.com/pubnub/go/pnerr"
	uuid "github.com/satori/go.uuid"
)

// EnumArrayToStringArray converts a string enum to an array
func EnumArrayToStringArray(include string) []string {
	f := strings.Fields(include)
	j := strings.Join(f, ",")
	t := strings.Trim(j, "[]")
	return strings.Fields(t)
}

// JoinChannels encodes and joins channels
func JoinChannels(channels []string) []byte {
	if len(channels) == 0 {
		return []byte(",")
	}

	var encodedChannels []string

	for _, value := range channels {
		encodedChannels = append(encodedChannels, URLEncode(value))
	}

	return []byte(strings.Join(encodedChannels, ","))
}

// encodeJSONAsPathComponent properly encodes serialized JSON
// for placement within a URI path
func EncodeJSONAsPathComponent(jsonBytes string) string {
	u := &url.URL{Path: jsonBytes}
	encodedPath := u.String()

	// Go 1.8 inserts a ./ per RFC 3986 §4.2. Previous versions
	// will be unaffected by this under the assumption that jsonBytes
	// represents valid JSON
	return strings.TrimLeft(encodedPath, "./")
}

func Serialize(msg interface{}) ([]byte, error) {
	jsonSerialized, errJSONMarshal := json.Marshal(msg)
	if errJSONMarshal != nil {
		return []byte{}, errJSONMarshal
	}
	return jsonSerialized, nil
}

func SerializeAndEncrypt(msg interface{}, cipherKey string, serialize bool) (string, error) {
	var encrypted string
	if serialize {
		jsonSerialized, errJSONMarshal := json.Marshal(msg)
		if errJSONMarshal != nil {
			return "", errJSONMarshal
		}
		encrypted = EncryptString(cipherKey, string(jsonSerialized))
	} else {
		if serializedMsg, ok := msg.(string); ok {
			encrypted = EncryptString(cipherKey, serializedMsg)
		} else {
			return "", pnerr.NewBuildRequestError("Message is not JSON serialized.")
		}
	}

	return encrypted, nil
}

func SerializeEncryptAndSerialize(msg interface{}, cipherKey string, serialize bool) (string, error) {
	var encrypted string

	if serialize {
		jsonSerialized, errJSONMarshal := json.Marshal(msg)
		if errJSONMarshal != nil {
			return "", errJSONMarshal
		}
		encrypted = EncryptString(cipherKey, string(jsonSerialized))
	} else {
		if serializedMsg, ok := msg.(string); ok {
			encrypted = EncryptString(cipherKey, serializedMsg)
		} else {
			return "", pnerr.NewBuildRequestError("Message is not JSON serialized.")
		}
	}
	jsonSerialized, errJSONMarshal := json.Marshal(encrypted)
	if errJSONMarshal != nil {
		return "", errJSONMarshal
	}
	return string(jsonSerialized), nil
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
func UUID() string {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	return id.String()
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
			stringifiedQuery += fmt.Sprintf("%s=%s", v, PamEncode((*params)[v][0]))
		} else {
			stringifiedQuery += fmt.Sprintf("%s=%s&", v, PamEncode((*params)[v][0]))
		}

		i++
	}

	return stringifiedQuery
}

func PamEncode(value string) string {
	// *!'()[]~
	stringifiedParam := URLEncode(value)

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
func URLEncode(s string) string {
	v := url.QueryEscape(s)

	var replacer = strings.NewReplacer(
		"+", "%20",
	)

	return replacer.Replace(v)
}
