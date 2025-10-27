// snippet.includes
// Replace with your package name (usually "main")
package pubnub_samples_test

import (
	"fmt"
	"os"

	pubnub "github.com/pubnub/go/v7"
)

// snippet.end

/*
IMPORTANT NOTE FOR COPYING EXAMPLES:

Throughout this file, you'll see code between "snippet.hide" and "snippet.show" comments.
These sections are used for CI/CD testing and should be SKIPPED if you're copying examples.

Example of what to skip:
	// snippet.hide
	config = setPubnubExampleConfigData(config)  // <- Skip this line (for testing only)
	defer pn.DeleteFile().Execute()              // <- Skip this line (cleanup for tests)
	// snippet.show

When copying examples to your own code:
- Use your own publish/subscribe keys instead of the "demo" keys
- Remove any statements that are between snippet.hide and snippet.show (they're only for testing purposes)
*/

// snippet.send_file
// Example_sendFile demonstrates uploading a file to a channel
func Example_sendFile() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo" // Replace with your subscribe key
	config.PublishKey = "demo"   // Replace with your publish key

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	// Create a test file for testing
	testFile, _ := os.CreateTemp("", "test_*.txt")
	testFile.WriteString("Hello, this is a test file!")
	testFile.Close()
	defer os.Remove(testFile.Name())
	// snippet.show

	// Open the file you want to upload
	file, err := os.Open(
		testFile.Name(), // Replace with testFile.Name() with your file path
	)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Upload file to channel
	response, status, err := pn.SendFile().
		Channel("my-channel").
		Name("my_text_file.txt").        // Name of the file
		File(file).                      // File to upload
		Message("Check out this file!"). // Optional message
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// snippet.hide
	// Cleanup uploaded file
	if response != nil && response.Data.ID != "" {
		pn.DeleteFile().
			Channel("my-channel").
			ID(response.Data.ID).
			Name("my_text_file.txt").
			Execute()
	}
	// snippet.show

	if status.StatusCode == 200 && response.Data.ID != "" {
		fmt.Println("File uploaded successfully")
	}

	// Output:
	// File uploaded successfully
}

// snippet.send_file_with_meta
// Example_sendFileWithMeta demonstrates uploading a file with metadata
func Example_sendFileWithMeta() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	// Create a test file for testing
	testFile, _ := os.CreateTemp("", "test_*.pdf")
	testFile.WriteString("File with metadata")
	testFile.Close()
	defer os.Remove(testFile.Name())
	// snippet.show

	// Open the file you want to upload
	file, err := os.Open(
		testFile.Name(), // Replace testFile.Name() with your file path
	)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Create metadata for the file
	metadata := map[string]interface{}{
		"fileType": "document",
		"category": "reports",
	}

	// Upload file with metadata
	response, status, err := pn.SendFile().
		Channel("my-channel").
		Name("report.pdf").
		File(file).
		Message("Monthly report").
		Meta(metadata). // Add metadata
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// snippet.hide
	if response != nil && response.Data.ID != "" {
		pn.DeleteFile().
			Channel("my-channel").
			ID(response.Data.ID).
			Name("report.pdf").
			Execute()
	}
	// snippet.show

	if status.StatusCode == 200 && response.Data.ID != "" {
		fmt.Println("File with metadata uploaded successfully")
	}

	// Output:
	// File with metadata uploaded successfully
}

// snippet.list_files
// Example_listFiles demonstrates listing all files in a channel
func Example_listFiles() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// List all files in the channel
	response, status, err := pn.ListFiles().
		Channel("files-channel").
		Limit(25). // Limit number of results
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Printf("Found %d file(s) in channel\n", len(response.Data))
	}
}

// snippet.get_file_url
// Example_getFileURL demonstrates generating a download URL for a file
func Example_getFileURL() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	// Upload a test file first for testing
	testFile, _ := os.CreateTemp("", "test_*.txt")
	testFile.WriteString("Test file for URL generation")
	testFile.Close()
	defer os.Remove(testFile.Name())

	file, _ := os.Open(testFile.Name())
	uploadResp, _, _ := pn.SendFile().
		Channel("files-channel").
		Name("sample.txt").
		File(file).
		Execute()
	file.Close()

	if uploadResp == nil || uploadResp.Data.ID == "" {
		return
	}
	// snippet.show

	// Get the download URL for a specific file
	// You would get these values from your file upload response or ListFiles
	fileID := uploadResp.Data.ID // Replace uploadResp.Data.ID with your file ID
	fileName := "sample.txt"     // Replace with your file name

	// snippet.hide
	defer pn.DeleteFile().
		Channel("files-channel").
		ID(fileID).
		Name(fileName).
		Execute()
	// snippet.show

	response, status, err := pn.GetFileURL().
		Channel("files-channel").
		ID(fileID).     // File ID from upload response
		Name(fileName). // File name from upload response
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 && response.URL != "" {
		fmt.Println("File URL generated successfully")
	}

	// Output:
	// File URL generated successfully
}

// snippet.download_file
// Example_downloadFile demonstrates downloading a file from a channel
func Example_downloadFile() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	// Upload a test file first for testing
	testFile, _ := os.CreateTemp("", "test_*.txt")
	testFile.WriteString("Test file for download")
	testFile.Close()
	defer os.Remove(testFile.Name())

	file, _ := os.Open(testFile.Name())
	uploadResp, _, _ := pn.SendFile().
		Channel("files-channel").
		Name("download.txt").
		File(file).
		Execute()
	file.Close()

	if uploadResp == nil || uploadResp.Data.ID == "" {
		return
	}
	// snippet.show

	// Download a specific file
	// You would get these values from your file upload response or ListFiles
	fileID := uploadResp.Data.ID // Replace uploadResp.Data.ID with your file ID
	fileName := "download.txt"   // Replace with your file name

	// snippet.hide
	defer pn.DeleteFile().
		Channel("files-channel").
		ID(fileID).
		Name(fileName).
		Execute()
	// snippet.show

	_, status, err := pn.DownloadFile().
		Channel("files-channel").
		ID(fileID).     // File ID from upload or list response
		Name(fileName). // File name
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("File downloaded successfully")
	}

	// Output:
	// File downloaded successfully
}

// snippet.delete_file
// Example_deleteFile demonstrates deleting a file from a channel
func Example_deleteFile() {
	config := pubnub.NewConfigWithUserId(pubnub.UserId("demo-user"))
	config.SubscribeKey = "demo"
	config.PublishKey = "demo"

	// snippet.hide
	config = setPubnubExampleConfigData(config)
	// snippet.show

	pn := pubnub.NewPubNub(config)

	// snippet.hide
	// Upload a test file first for testing
	testFile, _ := os.CreateTemp("", "test_*.txt")
	testFile.WriteString("Test file for deletion")
	testFile.Close()
	defer os.Remove(testFile.Name())

	file, _ := os.Open(testFile.Name())
	uploadResp, _, _ := pn.SendFile().
		Channel("files-channel").
		Name("to-delete.txt").
		File(file).
		Execute()
	file.Close()

	if uploadResp == nil || uploadResp.Data.ID == "" {
		return
	}
	// snippet.show

	// Delete a specific file
	// You would get these values from your file upload response or ListFiles
	fileID := uploadResp.Data.ID // Replace uploadResp.Data.ID with your file ID
	fileName := "to-delete.txt"  // Replace with your file name

	_, status, err := pn.DeleteFile().
		Channel("files-channel").
		ID(fileID).     // File ID from upload or list response
		Name(fileName). // File name
		Execute()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if status.StatusCode == 200 {
		fmt.Println("File deleted successfully")
	}

	// Output:
	// File deleted successfully
}

// snippet.end
