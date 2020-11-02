package e2e

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	pubnub "github.com/pubnub/go"
	"github.com/pubnub/go/utils"
	"github.com/stretchr/testify/assert"
)

func processInterface(in interface{}) {
	switch v := in.(type) {
	case map[string]interface{}:
		for s, b := range v {
			fmt.Printf("%s: b=%v\n", s, b)
		}
	case map[string]*pubnub.PNFileMessageAndDetails:

		fmt.Printf("*pubnub.PNFileMessageAndDetails")
	case map[string]pubnub.PNFileMessageAndDetails:

		fmt.Printf("pubnub.PNFileMessageAndDetails")
	default:
		fmt.Println("unknown type")
	}
}

func TestFileUpload(t *testing.T) {
	FileUploadCommon(t, false, "", "file_upload_test.txt", "file_upload_test_output.txt")
}

func TestFileUploadWithCipher(t *testing.T) {
	FileUploadCommon(t, true, "", "file_upload_test.txt", "file_upload_test_output.txt")
}

func TestFileUploadWithCustomCipher(t *testing.T) {
	FileUploadCommon(t, true, "enigma2", "file_upload_test.txt", "file_upload_test_output.txt")
}

func FileUploadCommon(t *testing.T, useCipher bool, customCipher string, filepathInput, filepathOutput string) {
	assert := assert.New(t)

	pn := pubnub.NewPubNub(pamConfigCopy())
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

	var mut sync.RWMutex

	defer file.Close()
	if err != nil {
		fmt.Println("File open error : ", err)
		assert.Fail("File open error")
	} else {

		r := GenRandom()
		rno := r.Intn(99999)
		id := ""
		retURL := ""
		ch := fmt.Sprintf("test_file_upload_channel_%d", rno)
		name := fmt.Sprintf("test_file_upload_name_%d.txt", rno)
		message := fmt.Sprintf("test file %s", name)

		listener := pubnub.NewListener()
		exitListener := make(chan bool)
		messageMatch := false
		idMatch := false
		nameMatch := false
		urlMatch := false
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
						fmt.Println(fmt.Sprintf("%v", file))
						fmt.Println(fmt.Sprintf("file.File.PNMessage.Text: %s", file.File.PNMessage.Text))
						fmt.Println(fmt.Sprintf("file.File.PNFile.Name: %s", file.File.PNFile.Name))
						fmt.Println(fmt.Sprintf("file.File.PNFile.ID: %s", file.File.PNFile.ID))
						fmt.Println(fmt.Sprintf("file.File.PNFile.URL: %s", file.File.PNFile.URL))
						fmt.Println(fmt.Sprintf("file.Channel: %s", file.Channel))
						fmt.Println(fmt.Sprintf("file.Timetoken: %d", file.Timetoken))
						fmt.Println(fmt.Sprintf("file.SubscribedChannel: %s", file.SubscribedChannel))
						fmt.Println(fmt.Sprintf("file.Publisher: %s", file.Publisher))
					}
					mut.Lock()
					messageMatch = message == file.File.PNMessage.Text
					idMatch = id == file.File.PNFile.ID
					nameMatch = name == file.File.PNFile.Name
					retURL = file.File.PNFile.URL
					if enableDebuggingInTests {
						fmt.Println("messageMatch:", messageMatch, message, file.File.PNMessage.Text)
						fmt.Println("idMatch:", idMatch, id, file.File.PNFile.ID)
						fmt.Println("nameMatch:", nameMatch, name, file.File.PNFile.Name)
					}
					mut.Unlock()

				case <-exitListener:
					break ExitLabel

				}
			}
		}()

		pn.AddListener(listener)

		pn.Subscribe().Channels([]string{ch}).Execute()

		resSendFile, statusSendFile, _ := pn.SendFile().Channel(ch).Message(message).CipherKey(cipherKey).Name(name).File(file).Execute()
		assert.Equal(200, statusSendFile.StatusCode)
		if enableDebuggingInTests {
			fmt.Println("statusSendFile.AdditionalData:", statusSendFile.AdditionalData)
		}

		if resSendFile != nil {
			mut.Lock()
			id = resSendFile.Data.ID
			mut.Unlock()
			if enableDebuggingInTests {
				fmt.Println("resSendFile.Data.ID ==>", resSendFile.Data.ID)
			}
			assert.NotEqual(0, resSendFile.Timestamp)
			time.Sleep(2 * time.Second)

			resGetFile, statusGetFile, errGetFile := pn.GetFileURL().Channel(ch).ID(id).Name(name).Execute()
			if enableDebuggingInTests {
				fmt.Println(statusGetFile)
			}
			assert.Equal(200, statusGetFile.StatusCode)
			assert.Nil(errGetFile)

			if resGetFile != nil {
				location := resGetFile.URL

				secure := ""
				if pn.Config.Secure {
					secure = "s"
				}
				if enableDebuggingInTests {
					fmt.Println("urlMatch:", urlMatch, retURL, location)
				}
				mut.Lock()
				i1 := strings.Index(retURL, "?")
				i2 := strings.Index(location, "?")
				retURL = retURL[:i1]
				location = location[:i2]
				urlMatch = retURL == location

				if enableDebuggingInTests {
					fmt.Println("urlMatch:", urlMatch, retURL, location)
				}

				path := fmt.Sprintf("v1/files/%s/channels/%s/files/%s/%s", pn.Config.SubscribeKey, ch, id, name)
				locationTest := fmt.Sprintf("http%s://%s/%s", secure, pn.Config.Origin, path)
				if enableDebuggingInTests {
					fmt.Println("location:", location)
				}
				assert.Contains(location, locationTest)
				mut.Unlock()
			}

			mut.Lock()
			if enableDebuggingInTests {
				fmt.Println("2 messageMatch:", messageMatch)
				fmt.Println("2 idMatch:", idMatch)
				fmt.Println("2 nameMatch:", nameMatch)
				fmt.Println("2 urlMatch:", urlMatch)
			}
			assert.True(nameMatch && idMatch && messageMatch && urlMatch)
			mut.Unlock()

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
					if msg.Text == message && file.ID == id && file.Name == name && chMessages[i].MessageType == 4 && chMessages[i].UUID == pn.Config.UUID {
						bFoundInFetch = true
						break
					}
				}

			}
			assert.True(bFoundInFetch)

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

		} else {
			assert.Fail("resSendFile nil")
		}

	}
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

func unpadPKCS7(data []byte) ([]byte, error) {
	blocklen := 16
	if len(data)%blocklen != 0 || len(data) == 0 {
		return nil, fmt.Errorf("invalid data len %d", len(data))
	}
	padlen := int(data[len(data)-1])
	if padlen > blocklen || padlen == 0 {
		return nil, fmt.Errorf("padding is invalid")
	}
	// check padding
	pad := data[len(data)-padlen:]
	for i := 0; i < padlen; i++ {
		if pad[i] != byte(padlen) {
			return nil, fmt.Errorf("padding is invalid")
		}
	}

	return data[:len(data)-padlen], nil
}

func FileEncryptDecryptWithoutBase64Test(t *testing.T) {
	assert := assert.New(t)
	filepathInput := "file_upload_test_enc.txt"
	//filepathInput = "whoami_out.txt"
	// filepathInput = "file_dec_3642342.png"
	filepathInput = "tux.jpg"
	// filepathInput = "gopher.jpg"
	filepathInput = "Moth.jpg"
	// filepathInput = "file_upload_test.txt"
	filepathOutput := "file_upload_test_enc_out.txt"
	out, _ := os.Create(filepathOutput)
	file, err := os.Open(filepathInput)
	key := utils.EncryptCipherKey("enigma")
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	iv := make([]byte, aes.BlockSize)

	if _, err := rand.Read(iv); err != nil {
		panic(err)
	}
	// IV that was breaking TEST
	// iv = []byte{73, 187, 52, 211, 20, 183, 129, 64, 119, 12, 190, 93, 7, 15, 70, 7}
	iv = []byte{133, 126, 158, 123, 43, 95, 96, 90, 215, 178, 17, 73, 166, 130, 79, 156}
	fmt.Println("iv:=>>>", iv)
	n, e := out.Write(iv)
	fmt.Println(n)
	if e != nil {
		panic(e)
	}

	blockSize := 16
	bufferSize := 16
	p := make([]byte, bufferSize)
	mode := cipher.NewCBCEncrypter(block, iv)
	cryptoRan := false
	fii, _ := file.Stat()
	contentLenIn := fii.Size()
	var contentRead int64

	for {
		n2, err2 := io.ReadFull(file, p)
		contentRead += int64(n2)
		if err2 != nil {
			fmt.Println("\nerr2:", err2)
			if err2 == io.EOF {
				ciphertext := make([]byte, blockSize)
				copy(ciphertext[:n2], p[:n2])
				fmt.Println("Encrypt EOF EOF")
				fmt.Println("ct2 EOF ", ciphertext, string(ciphertext), p[:n2], p)
				out.Close()
				break
			}

			if err2 == io.ErrUnexpectedEOF {
				if !cryptoRan {
					text := make([]byte, blockSize)
					ciphertext := make([]byte, blockSize)
					copy(text[:n2], p[:n2])
					pad := bytes.Repeat([]byte{byte(blockSize - n2)}, blockSize-n2)
					copy(text[n2:], pad)
					// ciphertext = padWithPKCS7(text)
					// fmt.Println(n2, padErr)

					mode.CryptBlocks(ciphertext, text)
					fmt.Println("ct1", ciphertext)
					out.Write(ciphertext)
					out.Close()
				} else {
					text := make([]byte, blockSize)
					ciphertext := make([]byte, blockSize)
					copy(text[:n2], p[:n2])
					pad := bytes.Repeat([]byte{byte(blockSize - n2)}, blockSize-n2)

					copy(text[n2:], pad)
					fmt.Println("text[n2:]", text)
					// ciphertext = padWithPKCS7(text)

					mode.CryptBlocks(ciphertext, text)
					fmt.Println("ct2", ciphertext, string(ciphertext))

					out.Write(ciphertext)

					out.Close()

				}

			}
			fmt.Println("Exiting For")
			break
		}

		ciphertext := make([]byte, blockSize)
		cryptoRan = true
		if contentRead >= contentLenIn {
			// p, padErr := pkcs7Pad(p, blockSize)
			pad := bytes.Repeat([]byte{byte(blockSize - n2)}, blockSize-n2)
			copy(p[n2:], pad)

			fmt.Println("n2, padErr, ciphertext, p, blockSize", n2, pad, ciphertext, p, blockSize)
		}

		mode.CryptBlocks(ciphertext, p)

		fmt.Println("ct0 write", ciphertext, string(ciphertext))
		out.Write(ciphertext)

	}

	out.Close()
	filepathDecOutput := "file_upload_test_enc_dec.txt"
	filepathDecOutput = "tux_out.jpg"
	// filepathDecOutput = "gopher_out.jpg"
	filepathDecOutput = "Moth_out.jpg"
	// filepathDecOutput = "file_upload_test_out.txt"
	out, _ = os.Open(filepathOutput)

	blockSize = 16
	bufferSize = 16
	p = make([]byte, bufferSize)
	ivBuff := make([]byte, blockSize)
	emptyByteVar := make([]byte, blockSize)

	iv2 := make([]byte, blockSize)
	count := 0
	r, w := io.Pipe()
	done := make(chan bool)
	cryptoRan = false

	fi, _ := out.Stat()
	contentLenEnc := fi.Size()
	var contentDownloaded int64
	fmt.Println("contentLenEnc", contentLenEnc)

	go func() {
	ExitReadLabel:
		for {
			n2, err2 := io.ReadFull(out, p)
			fmt.Println("file contents ->", n2, len(p), p, contentDownloaded)
			if err2 != nil {
				fmt.Println("\nerr2:", err2)
				if err2 == io.EOF {
					ciphertext := make([]byte, blockSize)
					copy(ciphertext, p[:n2])
					ciphertext, _ = unpadPKCS7(ciphertext)
					w.Write(ciphertext)
					// fmt.Println("ct2 EOF ", ciphertext, string(ciphertext), p[:n2], p, diffLen)
					w.Close()
					break ExitReadLabel
				}

				if err2 == io.ErrUnexpectedEOF {
					if bytes.Equal(iv2, emptyByteVar) {
						copy(iv2, ivBuff[0:blockSize])
						fmt.Println("string IV err2:", string(iv2), n2)
						fmt.Println("Extracted IV err2:", iv2, len(iv2))
						mode = cipher.NewCBCDecrypter(block, iv2)
						done <- true
					}
					if !cryptoRan {
						text := make([]byte, blockSize)
						ciphertext := make([]byte, blockSize)
						copy(text, p[:n2])
						mode.CryptBlocks(ciphertext, text)
						ciphertext, _ = unpadPKCS7(ciphertext)
						w.Write(ciphertext)

						w.Close()
						break ExitReadLabel
					} else {
						// text := make([]byte, blockSize)
						ciphertext := make([]byte, blockSize)
						copy(ciphertext, p[:n2])
						ciphertext, _ = unpadPKCS7(ciphertext)
						w.Write(ciphertext)

						// w.Write(ciphertextCopy)

						w.Close()
						break ExitReadLabel
					}

				}
				fmt.Println("Exiting For")
				break ExitReadLabel
			} else {
				contentDownloaded += int64(n2)
				if count < blockSize/bufferSize {
					fmt.Println(string(p[:n2]))

					// If error is not nil then panics
					if err != nil {
						panic(err)
					}
					copy(ivBuff[bufferSize*count:], p)
					fmt.Println(string(ivBuff), n2)
				} else {

					if bytes.Equal(iv2, emptyByteVar) {
						copy(iv2, ivBuff[0:blockSize])
						fmt.Println("string IV:", string(iv2), n2)
						fmt.Println("Extracted IV:", iv2, len(iv2))
						mode = cipher.NewCBCDecrypter(block, iv2)
						done <- true
					}

					ciphertext := make([]byte, blockSize)

					text := make([]byte, blockSize)
					copy(text, p[:n2])
					diffLen := contentDownloaded - contentLenEnc
					fmt.Println("p===>", p, text, n2, diffLen, contentDownloaded, contentLenEnc)

					mode.CryptBlocks(ciphertext, p)
					cryptoRan = true
					if contentDownloaded >= contentLenEnc {
						ciphertext, _ = unpadPKCS7(ciphertext)
						w.Write(ciphertext)
					} else {
						fmt.Println("ct0 read:", ciphertext, string(ciphertext), n2, p)

						w.Write(ciphertext)
					}
				}
			}
			count++

		}
	}()
	fmt.Println("before done")
	<-done
	fmt.Println("after done")

	outD, _ := os.Create(filepathDecOutput)
	io.Copy(outD, r)

	fileText, _ := ioutil.ReadFile(filepathInput)
	fileTextOut, _ := ioutil.ReadFile(filepathDecOutput)
	fileOut, err := os.Open(filepathInput)
	fio, _ := fileOut.Stat()
	contentLenOut := fio.Size()
	fmt.Println(contentLenIn, contentLenEnc, contentLenOut, contentRead)
	assert.Equal(contentLenIn, contentLenOut)
	assert.Equal(contentLenIn, contentRead)
	// fileTextEncoded, _ := ioutil.ReadFile("file_upload_test_enc_out.txt")
	// fileTextEncoded2, _ := ioutil.ReadFile("whoami_enc.txt")
	// assert.Equal(fileTextEncoded2, fileTextEncoded)
	// fmt.Println("---FileText---")
	// fmt.Println(string(fileText))
	// fmt.Println("---FileTextOut---")
	// fmt.Println(string(fileTextOut))
	assert.Equal(fileText, fileTextOut)
}

func FileDownload3Test(t *testing.T) {
	//assert := assert.New(t)
	b := []byte{136, 230, 36, 7, 53, 165, 35, 60, 127, 128, 184, 60, 131, 24, 6, 161, 41}
	fmt.Println(string(b))

	filepathInput := "video_enc.mp4"
	filepathDecOutput := "video_enc_out.mp4"
	filepathInput = "gif_test.gif"
	filepathDecOutput = "gif_test_out.gif"
	filepathInput = "file_test.png"
	filepathDecOutput = "file_test_out.png"
	filepathInput = "whoami.txt"
	filepathDecOutput = "whoami_out.txt"
	filepathInput = "file_upload_original_encrypted.txt"
	filepathDecOutput = "file_upload_original_encrypted_out.txt"

	key := utils.EncryptCipherKey("enigma")
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	out, err := os.Open(filepathInput)

	blockSize := 16
	bufferSize := 16
	p := make([]byte, bufferSize)
	ivBuff := make([]byte, blockSize)
	emptyByteVar := make([]byte, blockSize)

	iv2 := make([]byte, blockSize)
	count := 0
	r, w := io.Pipe()

	var mode cipher.BlockMode

	done := make(chan bool)
	cryptoRan := false
	fi, _ := out.Stat()
	contentLenIn := fi.Size()
	var contentDownloaded int64

	go func() {
	ExitReadLabel:
		for {
			n2, err2 := io.ReadFull(out, p)
			fmt.Println("file contents ->", n2, len(p), p, contentDownloaded)
			if err2 != nil {
				fmt.Println("\nerr2:", err2)
				if err2 == io.EOF {
					ciphertext := make([]byte, blockSize)
					copy(ciphertext, p)
					fmt.Println("EOF EOF")
					diffLen := int(contentDownloaded - contentLenIn)
					fmt.Println("ct2 EOF ", ciphertext, string(ciphertext), p[:n2], p, diffLen)
					w.Close()
					break ExitReadLabel
				}

				if err2 == io.ErrUnexpectedEOF {
					if bytes.Equal(iv2, emptyByteVar) {
						copy(iv2, ivBuff[0:blockSize])
						fmt.Println("string IV err2:", string(iv2), n2)
						fmt.Println("Extracted IV err2:", iv2, len(iv2))
						mode = cipher.NewCBCDecrypter(block, iv2)
						done <- true
					}
					if !cryptoRan {
						text := make([]byte, blockSize)
						ciphertext := make([]byte, blockSize)
						copy(text, p[:n2])
						mode.CryptBlocks(ciphertext, text)
						extra := contentLenIn % int64(bufferSize)
						ciphertextCopy := make([]byte, extra)
						copy(ciphertextCopy[:extra], ciphertext[:extra])
						fmt.Println("ct1 ", ciphertextCopy)
						w.Write(ciphertextCopy)
						w.Close()
						break ExitReadLabel
					} else {
						text := make([]byte, blockSize)
						ciphertext := make([]byte, blockSize)
						copy(text, p[:n2])
						extra := contentLenIn % int64(bufferSize)
						ciphertextCopy := make([]byte, extra)
						copy(ciphertextCopy[:extra], ciphertext[:extra])
						fmt.Println("ct2 ", ciphertextCopy)

						fmt.Println("ct2 read", ciphertext, string(ciphertext))

						w.Write(ciphertextCopy)

						w.Close()
						break ExitReadLabel
					}

				}
				fmt.Println("Exiting For")
				break ExitReadLabel
			} else {
				if count < blockSize/bufferSize {
					fmt.Println(string(p[:n2]))

					// If error is not nil then panics
					if err != nil {
						panic(err)
					}
					copy(ivBuff[bufferSize*count:], p)
					fmt.Println(string(ivBuff), n2)
				} else {
					contentDownloaded += int64(n2)

					if bytes.Equal(iv2, emptyByteVar) {
						copy(iv2, ivBuff[0:blockSize])
						fmt.Println("string IV:", string(iv2), n2)
						fmt.Println("Extracted IV:", iv2, len(iv2))
						mode = cipher.NewCBCDecrypter(block, iv2)
						done <- true
					}

					ciphertext := make([]byte, blockSize)

					text := make([]byte, blockSize)
					copy(text, p[:n2])
					diffLen := contentDownloaded - contentLenIn
					fmt.Println("p===>", p, text, n2, diffLen, contentDownloaded, contentLenIn)

					mode.CryptBlocks(ciphertext, p)
					cryptoRan = true
					if contentDownloaded > contentLenIn {
						extra := contentLenIn % int64(bufferSize)
						ciphertextCopy := make([]byte, extra)
						copy(ciphertextCopy[:extra], ciphertext[:extra])
						fmt.Println("ciphertext ct01", ciphertextCopy, string(ciphertextCopy))
						fmt.Println("ct01 read:", ciphertext, string(ciphertext), n2, p)
						w.Write(ciphertextCopy)
					} else {
						fmt.Println("ct0 read:", ciphertext, string(ciphertext), n2, p)

						w.Write(ciphertext)
					}
				}
			}
			count++

		}
	}()
	fmt.Println("before done")
	<-done
	fmt.Println("after done")

	outD, _ := os.Create(filepathDecOutput)
	io.Copy(outD, r)

	fmt.Println("after reader")
	fmt.Println("---FileTextOut---")
	// fmt.Println(string(fileTextOut))

}
