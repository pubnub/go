package pubnub

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type customStruct struct {
	Foo string
	Bar []int
}

func TestParseCipherInterfaceCipherWithCipher(t *testing.T) {
	assert := assert.New(t)
	s := "Wi24KS4pcTzvyuGOHubiXg=="

	pn := NewPubNub(NewDemoConfig())
	pn.Config.CipherKey = "enigma"

	intf, err := parseCipherInterface(s, pn.Config)

	assert.Nil(err)
	assert.Equal("yay!", intf.(string))

}

func TestParseCipherInterfacePlainWithCipher(t *testing.T) {
	assert := assert.New(t)
	s := "yay!"

	pn := NewPubNub(NewDemoConfig())
	pn.Config.CipherKey = "enigma"

	intf, _ := parseCipherInterface(s, pn.Config)

	assert.Equal("yay!", intf.(string))
}

func TestParseCipherInterfaceCipherWithDiffCipher(t *testing.T) {
	assert := assert.New(t)
	s := "Wi24KS4pcTzvyuGOHubiXg=="

	pn := NewPubNub(NewDemoConfig())
	pn.Config.CipherKey = "test"

	intf, _ := parseCipherInterface(s, pn.Config)

	assert.Equal("Wi24KS4pcTzvyuGOHubiXg==", intf.(string))

}

func TestParseCipherInterfacePlainWithDiffCipher(t *testing.T) {
	assert := assert.New(t)
	s := "yay!"

	pn := NewPubNub(NewDemoConfig())
	pn.Config.CipherKey = "test"

	intf, _ := parseCipherInterface(s, pn.Config)

	assert.Equal("yay!", intf.(string))
}

func TestParseCipherInterfaceCipherWithoutCipher(t *testing.T) {
	assert := assert.New(t)
	s := "Wi24KS4pcTzvyuGOHubiXg=="

	pn := NewPubNub(NewDemoConfig())
	pn.Config.CipherKey = ""

	intf, _ := parseCipherInterface(s, pn.Config)

	assert.Equal("Wi24KS4pcTzvyuGOHubiXg==", intf.(string))
}

func TestParseCipherInterfacePlainWithoutCipher(t *testing.T) {
	assert := assert.New(t)
	s := "yay!"

	pn := NewPubNub(NewDemoConfig())
	pn.Config.CipherKey = ""

	intf, _ := parseCipherInterface(s, pn.Config)

	assert.Equal("yay!", intf.(string))
}

func TestParseCipherInterfacePlainWithCipherStruct(t *testing.T) {
	assert := assert.New(t)
	s := customStruct{
		Foo: "hi!",
		Bar: []int{1, 2, 3, 4, 5},
	}

	pn := NewPubNub(NewDemoConfig())
	pn.Config.CipherKey = "enigma"

	intf, err := parseCipherInterface(s, pn.Config)

	assert.Nil(err)
	if msg, ok := intf.(customStruct); !ok {
		assert.Fail(fmt.Sprintf("not map %s", reflect.TypeOf(intf).Kind()))
	} else {
		assert.Equal("hi!", msg.Foo)
		assert.Equal(2, msg.Bar[1])
	}

}

func TestParseCipherInterfacePlainWithoutCipherStruct(t *testing.T) {
	assert := assert.New(t)
	s := customStruct{
		Foo: "hi!",
		Bar: []int{1, 2, 3, 4, 5},
	}

	pn := NewPubNub(NewDemoConfig())
	pn.Config.CipherKey = ""

	intf, err := parseCipherInterface(s, pn.Config)

	assert.Nil(err)
	if msg, ok := intf.(customStruct); !ok {
		assert.Fail(fmt.Sprintf("not map %s", reflect.TypeOf(intf).Kind()))
	} else {
		assert.Equal("hi!", msg.Foo)
		assert.Equal(2, msg.Bar[1])
	}

}

func TestParseCipherInterfacePlainWithCipherMapPNOther(t *testing.T) {
	assert := assert.New(t)
	s1 := map[string]interface{}{
		"id":        1,
		"not_other": "1234567",
	}
	s := map[string]interface{}{
		"id":        1,
		"not_other": "12345",
		"pn_other":  s1,
	}

	pn := NewPubNub(NewDemoConfig())
	pn.Config.CipherKey = "enigma"

	intf, _ := parseCipherInterface(s, pn.Config)

	msg := intf.(map[string]interface{})
	assert.Equal("12345", msg["not_other"])
	if msgOther, ok := msg["pn_other"].(map[string]interface{}); !ok {
		assert.Fail("!map[string]interface{}")
	} else {
		assert.Equal("1234567", msgOther["not_other"])
	}
}

func TestParseCipherInterfacePlainWithoutCipherMapPNOther(t *testing.T) {
	assert := assert.New(t)
	s1 := map[string]interface{}{
		"id":        1,
		"not_other": "1234567",
	}
	s := map[string]interface{}{
		"id":        1,
		"not_other": "12345",
		"pn_other":  s1,
	}

	pn := NewPubNub(NewDemoConfig())
	pn.Config.CipherKey = ""

	intf, _ := parseCipherInterface(s, pn.Config)

	msg := intf.(map[string]interface{})
	assert.Equal("12345", msg["not_other"])
	if msgOther, ok := msg["pn_other"].(map[string]interface{}); !ok {
		assert.Fail("!map[string]interface{}")
	} else {
		assert.Equal("1234567", msgOther["not_other"])
	}
}

func TestParseCipherInterfaceCipherWithoutCipherStringPNOther(t *testing.T) {
	assert := assert.New(t)
	s := map[string]interface{}{
		"id":        1,
		"not_other": "1234",
		"pn_other":  "Wi24KS4pcTzvyuGOHubiXg==",
	}

	pn := NewPubNub(NewDemoConfig())
	pn.Config.CipherKey = ""

	intf, _ := parseCipherInterface(s, pn.Config)

	msg := intf.(map[string]interface{})
	assert.Equal("1234", msg["not_other"])
	if msgOther, ok := msg["pn_other"].(string); !ok {
		assert.Fail("!string")
	} else {
		assert.Equal("Wi24KS4pcTzvyuGOHubiXg==", msgOther)
	}
}

func TestParseCipherInterfaceCipherWithCipherStringPNOther(t *testing.T) {
	assert := assert.New(t)
	s := map[string]interface{}{
		"id":        1,
		"not_other": "1234",
		"pn_other":  "Wi24KS4pcTzvyuGOHubiXg==",
	}

	pn := NewPubNub(NewDemoConfig())
	pn.Config.CipherKey = "enigma"

	intf, _ := parseCipherInterface(s, pn.Config)

	msg := intf.(map[string]interface{})
	assert.Equal("1234", msg["not_other"])
	if msgOther, ok := msg["pn_other"].(string); !ok {
		assert.Fail("!string")
	} else {
		assert.Equal("yay!", msgOther)
	}
}

func TestParseCipherInterfaceCipherWithoutCipherStruct(t *testing.T) {
	assert := assert.New(t)
	s := "BMhiHh363wsb7kNk7krTtDcey/O6ZcoKDTvVc4yDhZY="

	pn := NewPubNub(NewDemoConfig())
	pn.Config.CipherKey = "enigma"

	intf, _ := parseCipherInterface(s, pn.Config)
	msg := intf.(map[string]interface{})
	assert.Equal("hi!", msg["Foo"])

}

func TestParseCipherInterfaceCipherWithCipherStructPNOther(t *testing.T) {
	assert := assert.New(t)
	s := map[string]interface{}{
		"id":        1,
		"not_other": "1234",
		"pn_other":  "BMhiHh363wsb7kNk7krTtDcey/O6ZcoKDTvVc4yDhZY=",
	}

	pn := NewPubNub(NewDemoConfig())
	pn.Config.CipherKey = "enigma"

	intf, _ := parseCipherInterface(s, pn.Config)

	msg := intf.(map[string]interface{})
	assert.Equal("1234", msg["not_other"])
	if msgOther, ok := msg["pn_other"].(map[string]interface{}); !ok {
		assert.Fail("!map[string]interface{}")
	} else {
		assert.Equal("hi!", msgOther["Foo"])
	}
}

func TestParseCipherInterfaceCipherWithOtherCipherStructPNOther(t *testing.T) {
	assert := assert.New(t)
	s := map[string]interface{}{
		"id":        1,
		"not_other": "1234",
		"pn_other":  "BMhiHh363wsb7kNk7krTtDcey/O6ZcoKDTvVc4yDhZY=",
	}

	pn := NewPubNub(NewDemoConfig())
	pn.Config.CipherKey = "test"

	intf, _ := parseCipherInterface(s, pn.Config)

	msg := intf.(map[string]interface{})
	assert.Equal("1234", msg["not_other"])
	assert.Equal("BMhiHh363wsb7kNk7krTtDcey/O6ZcoKDTvVc4yDhZY=", msg["pn_other"])
}

func TestParseCipherInterfaceCipherWithCipherStructPNOtherDisable(t *testing.T) {
	assert := assert.New(t)
	s := map[string]interface{}{
		"id":        1,
		"not_other": "1234",
		"pn_other":  "BMhiHh363wsb7kNk7krTtDcey/O6ZcoKDTvVc4yDhZY=",
	}

	pn := NewPubNub(NewDemoConfig())
	pn.Config.DisablePNOtherProcessing = true
	pn.Config.CipherKey = "enigma"

	intf, _ := parseCipherInterface(s, pn.Config)

	msg := intf.(map[string]interface{})
	assert.Equal("1234", msg["not_other"])
	assert.Equal("BMhiHh363wsb7kNk7krTtDcey/O6ZcoKDTvVc4yDhZY=", msg["pn_other"])

}

func TestParseCipherInterfaceCipherWithIntSlice(t *testing.T) {
	assert := assert.New(t)
	s := []int{1, 2, 3, 4, 5}

	pn := NewPubNub(NewDemoConfig())
	pn.Config.DisablePNOtherProcessing = true
	pn.Config.CipherKey = ""

	intf, _ := parseCipherInterface(s, pn.Config)

	msg := intf.([]int)
	assert.Equal(1, msg[0])

}

func TestParseCipherInterfaceCipherWithoutCipherStruct2(t *testing.T) {
	assert := assert.New(t)
	s := "kpNj0VFN5kkWBjbgQuG5DPkZGcJCKXFqQlZtaM7SLq2gHziTK1JlzQD/fxquAlGIwvM91wAT8KbBwxmDV3PTcP7KtY9whmhT1hSA9r1+RT4="

	pn := NewPubNub(NewDemoConfig())
	pn.Config.CipherKey = "enigma"

	intf, _ := parseCipherInterface(s, pn.Config)
	msg := intf.(map[string]interface{})
	assert.Equal("12345", msg["not_other"])
	if msgOther, ok := msg["pn_other"].(map[string]interface{}); !ok {
		assert.Fail("!map[string]interface{}")
	} else {
		assert.Equal(float64(1), msgOther["id"])
	}

}
