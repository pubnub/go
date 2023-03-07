package e2e

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	pubnub "github.com/pubnub/go/v7"
	"github.com/pubnub/go/v7/utils"
	"github.com/stretchr/testify/assert"
)

func TestFileUpload(t *testing.T) {
	FileUploadCommon(t, false, "", "file_upload_test.txt", "file_upload_test_output.txt")
}

func TestFileUploadWithCipher(t *testing.T) {
	FileUploadCommon(t, true, "", "file_upload_test.txt", "file_upload_test_output.txt")
}

func TestFileUploadWithCustomCipher(t *testing.T) {
	FileUploadCommon(t, true, "enigma2", "file_upload_test.txt", "file_upload_test_output.txt")
}

type FileData struct {
	id, url, name, message string
	messageType            pubnub.MessageType
	spaceId                pubnub.SpaceId
}

func FileUploadCommon(t *testing.T, useCipher bool, customCipher string, filepathInput, filepathOutput string) {
	assert := assert.New(t)
	config := pamConfigCopy()
	config.Log = log.Default()
	pn := pubnub.NewPubNub(config)
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}
	cipherKey := ""
	if useCipher {
		if customCipher != "" {
			cipherKey = customCipher
			pn.Config.CipherKey = "enigma"
		} else {
			pn.Config.CipherKey = "enigma"
			cipherKey = pn.Config.CipherKey
		}
	}

	file, err := os.Open(filepathInput)

	defer file.Close()
	if err != nil {
		fmt.Println("File open error : ", err)
		assert.Fail("File open error")
	}

	fileDataChannel := make(chan FileData)
	r := GenRandom()
	rno := r.Intn(99999)
	id := ""
	retURL := ""
	ch := fmt.Sprintf("test_file_upload_channel_%d", rno)
	name := fmt.Sprintf("test_file_upload_name_%d.txt", rno)
	message := fmt.Sprintf("test file %s", name)
	expectedMessageType := pubnub.MessageType("This_is_messageType")
	expectedSpaceId := pubnub.SpaceId("This_is_spaceId")

	listener := pubnub.NewListener()
	exitListener := make(chan bool)
	go func() {
	ExitLabel:
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					break
				default:
				}

			case file := <-listener.File:
				if enableDebuggingInTests {

					fmt.Println(" --- File: ")
					fmt.Printf("%v\n", file)
					fmt.Printf("file.File.PNMessage.Text: %s\n", file.File.PNMessage.Text)
					fmt.Printf("file.File.PNFile.Name: %s\n", file.File.PNFile.Name)
					fmt.Printf("file.File.PNFile.ID: %s\n", file.File.PNFile.ID)
					fmt.Printf("file.File.PNFile.URL: %s\n", file.File.PNFile.URL)
					fmt.Printf("file.Channel: %s\n", file.Channel)
					fmt.Printf("file.Timetoken: %d\n", file.Timetoken)
					fmt.Printf("file.SubscribedChannel: %s\n", file.SubscribedChannel)
					fmt.Printf("file.Publisher: %s\n", file.Publisher)
				}
				fileDataChannel <- FileData{file.File.PNFile.ID, file.File.PNFile.URL, file.File.PNFile.Name, file.File.PNMessage.Text, file.MessageType, file.SpaceId}

			case <-exitListener:
				break ExitLabel

			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().Channels([]string{ch}).Execute()
	// Sleep a bit, to give client some time to subscribe on channels firs.
	time.Sleep(100 * time.Millisecond)

	resSendFile, statusSendFile, err := pn.SendFile().
		Channel(ch).
		Message(message).
		CipherKey(cipherKey).
		Name(name).
		File(file).
		MessageType(expectedMessageType).
		SpaceId(expectedSpaceId).
		Execute()
	assert.Equal(200, statusSendFile.StatusCode)
	if enableDebuggingInTests {
		fmt.Println("statusSendFile.AdditionalData:", statusSendFile.AdditionalData)
	}

	if err != nil {
		close(fileDataChannel)
		t.Error(err)
		assert.Fail("resSendFile nil")
		return
	}

	id = resSendFile.Data.ID

	timer := time.NewTimer(5 * time.Second)
	var fileData FileData
	select {
	case fileData = <-fileDataChannel:
	case <-timer.C:
		assert.Fail("Timeout when waiting on file event")
		return
	}
	retURL = fileData.url

	if enableDebuggingInTests {
		fmt.Println("resSendFile.Data.ID ==>", resSendFile.Data.ID)
	}

	resGetFile, statusGetFile, errGetFile := pn.GetFileURL().Channel(ch).ID(id).Name(name).Execute()
	if enableDebuggingInTests {
		fmt.Println(statusGetFile)
	}
	assert.Equal(200, statusGetFile.StatusCode)
	assert.Nil(errGetFile)

	if resGetFile == nil {
		assert.Fail("resGetFile nil")
		return
	}

	location := resGetFile.URL

	secure := ""
	if pn.Config.Secure {
		secure = "s"
	}

	i1 := strings.Index(retURL, "?")
	i2 := strings.Index(location, "?")
	retURL = retURL[:i1]
	location = location[:i2]

	path := fmt.Sprintf("v1/files/%s/channels/%s/files/%s/%s", pn.Config.SubscribeKey, ch, id, name)
	locationTest := fmt.Sprintf("http%s://%s/%s", secure, pn.Config.Origin, path)
	if enableDebuggingInTests {
		fmt.Println("location:", location)
	}
	assert.Contains(location, locationTest)

	assert.Equal(name, fileData.name)
	assert.Equal(id, fileData.id)
	assert.Equal(message, fileData.message)
	assert.Equal(retURL, location)
	assert.Equal(expectedMessageType, fileData.messageType)
	assert.Equal(expectedSpaceId, fileData.spaceId)

	out, errDL := os.Create(filepathOutput)
	defer out.Close()
	if errDL != nil {
		if enableDebuggingInTests {
			fmt.Println(errDL)
		}
	} else {

		resDLFile, statusDLFile, errDLFile := pn.DownloadFile().Channel(ch).CipherKey(cipherKey).ID(id).Name(name).Execute()
		assert.Nil(errDLFile)
		if enableDebuggingInTests {
			fmt.Println("statusDLFile.StatusCode ===>", statusDLFile.StatusCode)
		}
		if resDLFile != nil {
			_, err := io.Copy(out, resDLFile.File)

			if err != nil {
				fmt.Println(err)
			} else {
				fileText, _ := ioutil.ReadFile(filepathInput)
				fileTextOut, _ := ioutil.ReadFile(filepathOutput)
				assert.Equal(fileText, fileTextOut)
			}
		}
	}

	fetchCall := func() error {
		ret1, _, _ := pn.FetchWithContext(backgroundContext).
			Channels([]string{ch}).
			Count(25).
			IncludeMessageType(true).
			IncludeUUID(true).
			Reverse(true).
			Execute()
		chMessages := ret1.Messages[ch]
		bFoundInFetch := false
		//{"status": 200, "error": false, "error_message": "", "channels": {"test_file_upload_channnel_86621":[{"message": {"message": {"text": "test file test_file_upload_name_86621"}, "file": {"name": "test_file_upload_name_86621", "id": "4c5644d4-4e18-48a1-924f-932252acea74"}}, "timetoken": "15935884935043935"}]}}
		for i := 0; i < len(chMessages); i++ {

			m := chMessages[i].Message
			file := chMessages[i].File
			if enableDebuggingInTests {
				fmt.Println("pubnub.PNFileDetails", file.ID)
				fmt.Println("pubnub.PNFileDetails", file.Name)
			}
			if msg, ok := m.(pubnub.PNPublishMessage); !ok {
				if enableDebuggingInTests {
					fmt.Println("!pubnub.PNPublishMessage")
				}
			} else {
				if enableDebuggingInTests {
					fmt.Println("pubnub.PNPublishMessage", msg.Text)
				}
				if msg.Text == message && file.ID == id && file.Name == name && chMessages[i].MessageType == "pn_file" && chMessages[i].UUID == pn.Config.UUID {
					bFoundInFetch = true
					break
				}
			}

		}
		if !bFoundInFetch {
			return errors.New("bFoundInFetch is false")
		}
		return nil
	}

	checkFor(assert, time.Second*3, time.Millisecond*500, fetchCall)

	resListFile, statusListFile, errListFile := pn.ListFiles().Channel(ch).Execute()
	assert.Nil(errListFile)
	assert.Equal(200, statusListFile.StatusCode)

	if resListFile != nil {
		bFound := false
		for _, m := range resListFile.Data {
			if enableDebuggingInTests {
				fmt.Println("file =====> ", m.ID, m.Created, m.Name, m.Size)
			}
			if m.ID == id && m.Name == name {
				bFound = true
			}
		}
		assert.True(resListFile.Count > 0)
		assert.True(bFound)
	}

	_, statusDelFile, errDelFile := pn.DeleteFile().Channel(ch).ID(id).Name(name).Execute()
	assert.Nil(errDelFile)
	assert.Equal(200, statusDelFile.StatusCode)

	// _, statusGetFile2, _ := pn.DownloadFile().Channel(ch).ID(id).Name(name).Execute()
	// assert.Equal(404, statusGetFile2.StatusCode)

}

func TestFileEncryptionDecryption(t *testing.T) {
	assert := assert.New(t)
	filepathInput := "file_upload_test.txt"
	filepathOutput := "file_upload_test_output.txt"
	filepathSampleOutput := "file_upload_sample_encrypted.txt"
	filepathOutputDec := "file_upload_dec_output.txt"

	out, _ := os.Create(filepathOutput)
	file, err := os.Open(filepathInput)
	if err != nil {
		panic(err)
	}
	utils.EncryptFile("enigma", []byte{133, 126, 158, 123, 43, 95, 96, 90, 215, 178, 17, 73, 166, 130, 79, 156}, out, file)
	fileText, _ := ioutil.ReadFile(filepathOutput)

	fileTextSample, _ := ioutil.ReadFile(filepathSampleOutput)
	assert.Equal(string(fileTextSample), string(fileText))

	outDec, _ := os.Open(filepathSampleOutput)
	fi, _ := outDec.Stat()
	contentLenEnc := fi.Size()
	defer outDec.Close()

	fileDec, _ := os.Create(filepathOutputDec)
	defer fileDec.Close()
	r, w := io.Pipe()
	utils.DecryptFile("enigma", contentLenEnc, outDec, w)
	io.Copy(fileDec, r)
	fileTextDec, _ := ioutil.ReadFile(filepathOutputDec)
	fileTextIn, _ := ioutil.ReadFile(filepathInput)
	assert.Equal(fileTextIn, fileTextDec)

}
