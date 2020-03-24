package utils

import (
	"bytes"
	"encoding/xml"
)

// Data represents a <data> element.
type data struct {
	XMLName xml.Name `xml:"data"`
	//Entry   []Entry  `xml:"entry"`
	Name string `xml:"name"`
	Age  int    `xml:"age"`
}

// PubnubDemoMessage is a struct to test a non-alphanumeric message
type pubnubDemoMessage struct {
	DefaultMessage string `json:",string"`
}

// CustomComplexMessage is used to test the custom structure encryption and decryption.
// The variables "foo" and "bar" as used in the other languages are not
// accepted by golang and give an empty value when serialized, used "Foo"
// and "Bar" instead.
type customComplexMessage struct {
	VersionID     float32 `json:",string"`
	TimeToken     int64   `json:",string"`
	OperationName string
	Channels      []string
	DemoMessage   pubnubDemoMessage `json:",string"`
	SampleXML     string            `json:",string"`
}

// InitComplexMessage initializes a complex structure of the
// type CustomComplexMessage which includes a xml, struct of type PubnubDemoMessage,
// strings, float and integer.
func initComplexMessage() customComplexMessage {
	pubnubDemoMessage := pubnubDemoMessage{
		DefaultMessage: "~!@#$%^&*()_+ `1234567890-= qwertyuiop[]\\ {}| asdfghjkl;' :\" zxcvbnm,./ <>? ",
	}

	xmlDoc := &data{Name: "Doe", Age: 42}

	output := new(bytes.Buffer)
	enc := xml.NewEncoder(output)

	err := enc.Encode(xmlDoc)
	if err != nil {
		//fmt.Printf("error: %v\n", err)
		return customComplexMessage{}
	}

	customComplexMessage := customComplexMessage{
		VersionID:     3.4,
		TimeToken:     13601488652764619,
		OperationName: "Publish",
		Channels:      []string{"ch1", "ch 2"},
		DemoMessage:   pubnubDemoMessage,
		SampleXML:     output.String(),
	}
	return customComplexMessage
}
