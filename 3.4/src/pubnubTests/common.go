package pubnubTests

import(
	"fmt"
	"encoding/xml"
)

var publishSuccessMessage = "1,\"Sent\""


type EmptyStruct struct {
}

// CustomStruct to test the custom structure encryption and decryption
// The variables "foo" and "bar" give an empty value when serialized, used "Foo" and "Bar" instead 
type CustomStruct struct {
    Foo string
    Bar []int
}

// CustomStruct to test the custom structure encryption and decryption
// The variables "foo" and "bar" give an empty value when serialized, used "Foo" and "Bar" instead 
type CustomSingleElementStruct struct {
    Foo string
}

// CustomStruct to test the custom structure encryption and decryption
// The variables "foo" and "bar" give an empty value when serialized, used "Foo" and "Bar" instead 
type CustomComplexMessage struct {
    VersionId 		float32   
    TimeToken 		int64
    OperationName 	string
    Channels 		[]string
    DemoMessage 	PubnubDemoMessage
    SampleXml 		[]byte
}

type PubnubDemoMessage struct {
	DefaultMessage string
}

func InitComplexMessage() CustomComplexMessage{
	pubnubDemoMessage := PubnubDemoMessage{
		DefaultMessage:  "~!@#$%^&*()_+ `1234567890-= qwertyuiop[]\\ {}| asdfghjkl;' :\" zxcvbnm,./ <>? ",
	}
	
	xmlDoc := &Person{Id: 13, FirstName: "John", LastName: "Doe", Age: 42}
	xmlDoc.Comment = " Need more details. "
	xmlDoc.Address = Address{"Hanga Roa", "Easter Island"}
	
	//_, err := xml.MarshalIndent(xmlDoc, "  ", "    ")
	output, err := xml.MarshalIndent(xmlDoc, "  ", "    ")
	if err != nil {
	    fmt.Printf("error: %v\n", err)
	    return CustomComplexMessage{}
	}
	customComplexMessage := CustomComplexMessage{
	    VersionId		: 3.4,   
	    TimeToken 		: 13601488652764619,
	    OperationName	: "Publish",
		Channels		: []string{"ch1"},
		DemoMessage 	: pubnubDemoMessage,
		//SampleXml		: xmlDoc,
		SampleXml		: output,
	}
	return customComplexMessage
}

type Address struct {
    City, State string
}

type Person struct {
    XMLName   xml.Name `xml:"person"`
    Id        int      `xml:"id,attr"`
    FirstName string   `xml:"name>first"`
    LastName  string   `xml:"name>last"`
    Age       int      `xml:"age"`
    Height    float32  `xml:"height,omitempty"`
    Married   bool
    Address
    Comment string `xml:",comment"`
}

func PrintTestMessage(message string){
	fmt.Println(" ")
	fmt.Println(message)
	fmt.Println(" ")
}
