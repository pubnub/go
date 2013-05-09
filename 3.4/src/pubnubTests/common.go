// Package pubnubMessaging has the unit tests of package pubnubMessaging.
// common file has the reused methods across the varoius unit test files  
package pubnubTests

import(
    "fmt"
    "encoding/xml"
    "bytes"
    "strings"
)

// publishSuccessMessage: the reponse that is received when a message is 
// successfully published on a pubnub channel
var publishSuccessMessage = "1,\"Sent\""

// EmptyStruct provided the empty struct to test the encryption
type EmptyStruct struct {
}

// CustomStruct to test the custom structure encryption and decryption
// The variables "foo" and "bar" give an empty value when serialized, used "Foo" and "Bar" instead 
type CustomStruct struct {
    Foo string
    Bar []int
}

// CustomSingleElementStruct: Used to test the custom structure encryption and decryption
// The variables "foo" and "bar" give an empty value when serialized, used "Foo" and "Bar" instead 
type CustomSingleElementStruct struct {
    Foo string
}

// CustomComplexMessage is used to test the custom structure encryption and decryption
// The variables "foo" and "bar" give an empty value when serialized, used "Foo" and "Bar" instead 
type CustomComplexMessage struct {
    VersionId         float32        `json:",string"`
    TimeToken         int64        `json:",string"`
    OperationName     string
    Channels         []string
    DemoMessage     PubnubDemoMessage `json:",string"`
    SampleXml         string         `json:",string"`
}

// 
type PubnubDemoMessage struct {
    DefaultMessage string `json:",string"`
}

func InitComplexMessage() CustomComplexMessage{
    pubnubDemoMessage := PubnubDemoMessage{
        DefaultMessage:  "~!@#$%^&*()_+ `1234567890-= qwertyuiop[]\\ {}| asdfghjkl;' :\" zxcvbnm,./ <>? ",
    }
    
    xmlDoc := &Data{Name:"Doe", Age:42 }
    
    //_, err := xml.MarshalIndent(xmlDoc, "  ", "    ")
    //output, err := xml.MarshalIndent(xmlDoc, "  ", "    ")
    output := new(bytes.Buffer) 
    enc := xml.NewEncoder(output)
    
    err := enc.Encode(xmlDoc)
    if err != nil {
        fmt.Printf("error: %v\n", err)
        return CustomComplexMessage{}
    }
    //fmt.Printf("xmlDoc: %v\n", xmlDoc)    
    customComplexMessage := CustomComplexMessage{
        VersionId        : 3.4,   
        TimeToken         : 13601488652764619,
        OperationName    : "Publish",
        Channels        : []string{"ch1", "ch 2"},
        DemoMessage     : pubnubDemoMessage,
        //SampleXml        : xmlDoc,
        SampleXml        : output.String(),
    }
    return customComplexMessage
}

// Represents a <data> element
type Data struct {
    XMLName xml.Name `xml:"data"`
    //Entry   []Entry  `xml:"entry"`
    Name string `xml:"name"`
    Age  int    `xml:"age"`
}

// Represents an <entry> element
type Entry struct {
    Name string `xml:"name"`
    Age  int    `xml:"age"`
}

func PrintTestMessage(message string){
    fmt.Println(" ")
    fmt.Println(message)
    fmt.Println(" ")
}

func ReplaceEncodedChars(str string) string{
    str = strings.Replace(str, "\\u003c", "<", -1)
    str = strings.Replace(str, "\\u003e", ">", -1)
    str = strings.Replace(str, "\\u0026", "&", -1)
    return str
}

