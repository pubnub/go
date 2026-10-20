# PubNub Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/pubnub/go/v10.svg)](https://pkg.go.dev/github.com/pubnub/go/v10)

PubNub provides global infrastructure for real-time, interactive applications.

Publish and receive messages in Go. Use this SDK for Go server, CLI, desktop, and backend applications.

[Documentation](https://www.pubnub.com/docs/sdks/go) · [API reference](https://www.pubnub.com/docs/sdks/go/api-reference/publish-and-subscribe) · [Changelog](https://www.pubnub.com/docs/sdks/go/changelog) · [Support](https://support.pubnub.com/)

## Requirements

| Requirement | Supported version or setup |
| --- | --- |
| Go | [Go 1.25 or later](https://www.pubnub.com/docs/sdks/go) |
| Module | [`github.com/pubnub/go/v10`](https://github.com/pubnub/go) |
| Go modules | Initialize a Go module before adding the SDK. |

## Installation

From the root of the Go module:

```sh
go get github.com/pubnub/go/v10
```

For other installation methods, see [environment setup](https://www.pubnub.com/docs/sdks/go).

## Quickstart

This example runs in a Go 1.25+ command-line application. It subscribes to `hello_world`, publishes one message, and prints the received text.

### Get your keys

1. Open the [PubNub Admin Portal](https://admin.pubnub.com/).
2. Create an app and a keyset for development, or select an existing development keyset.
3. Copy its publish key and subscribe key.

If Access Manager is enabled, obtain a token from your trusted backend and call `pn.SetToken(token)`. Keep the secret key on the backend.

### Send and receive a message

Create `main.go` in the module root and add:

Replace the key placeholders. Use a User ID that identifies the user or device in your app.

```go
package main

import (
	"fmt"

	pubnub "github.com/pubnub/go/v10"
)

func main() {
	// Replace with your keys from the Admin Portal
	config := pubnub.NewConfigWithUserId(pubnub.UserId("hello-world-user"))
	config.PublishKey = "YOUR_PUBLISH_KEY"
	config.SubscribeKey = "YOUR_SUBSCRIBE_KEY"

	pn := pubnub.NewPubNub(config)

	listener := pubnub.NewListener()
	pn.AddListener(listener)

	connected := make(chan struct{})
	received := make(chan struct{})

	// Print the message when it arrives, then stop listening
	go func() {
		for {
			select {
			case status := <-listener.Status:
				// The connection is ready a moment after subscribe
				if status.Category == pubnub.PNConnectedCategory {
					close(connected)
				}
			case message := <-listener.Message:
				// Print the content of the received message
				fmt.Println(message.Message)

				// Unsubscribe to stop receiving messages. This sample needs only one message.
				pn.UnsubscribeAll()
				close(received)
				return
			}
		}
	}()

	// Start receiving messages on hello_world
	pn.Subscribe().
		Channels([]string{"hello_world"}).
		Execute()

	// Publish only after the subscription is connected, so this client receives the message
	<-connected

	// Publish a "Hello world" message to the hello_world channel
	_, _, err := pn.Publish().
		Channel("hello_world").
		Message("Hello world").
		Execute()
	if err != nil {
		fmt.Println(err)
		return
	}

	<-received
}
```

### Run the example

```sh
go run main.go
```

The terminal should show:

```text
Hello world
```

The example calls `pn.UnsubscribeAll()` after receiving the message.

For a complete application, see the [getting started guide](https://www.pubnub.com/docs/sdks/go).

## Next steps

| Task | Documentation |
| --- | --- |
| Configure the client | [Configuration](https://www.pubnub.com/docs/sdks/go/api-reference/configuration) |
| Work with subscriptions and messages | [Publish & Subscribe](https://www.pubnub.com/docs/sdks/go/api-reference/publish-and-subscribe) |
| Check channel occupancy | [Presence](https://www.pubnub.com/docs/sdks/go/api-reference/presence) |
| Read message history | [Message Persistence](https://www.pubnub.com/docs/sdks/go/api-reference/storage-and-playback) |

## Build with an AI coding assistant

The [PubNub MCP server](https://www.pubnub.com/docs/ai/pubnub-mcp-server) gives an AI coding assistant access to PubNub SDK documentation and PubNub APIs. Connect the assistant to the hosted server at `https://mcp.pubnub.com`, or run `npx @pubnub/mcp@latest` locally.

The [server repository](https://github.com/pubnub/pubnub-mcp-server) has setup steps for VS Code, Cursor, Claude Code, Claude Desktop, Codex, Gemini CLI, and OpenCode.

## Before production

Use [Access Manager](https://www.pubnub.com/docs/sdks/go/api-reference/access-manager) to grant each client the permissions it needs. Keep the secret key on your backend. Never include it in a distributed client.

Reuse the PubNub client for the service/session lifetime. Call `pn.UnsubscribeAll()` when the client should stop its subscriptions.

Live delivery through PubNub SDKs is at-most-once. A subscriber can miss messages while disconnected or if its buffer overflows. For longer-gap recovery, see [Message Persistence](https://www.pubnub.com/docs/sdks/go/api-reference/storage-and-playback).

## Troubleshooting

| Symptom | Check |
| --- | --- |
| Go cannot resolve the PubNub import | Use the current v10 module path `github.com/pubnub/go/v10`, then run `go mod tidy`. See the [Go repository](https://github.com/pubnub/go). |
| A publish succeeds but no message appears | Check the message handler, subscription readiness, keyset, and channel name. |
| The program exits before the subscribed message arrives | Wait for `PNConnectedCategory` before publishing and keep the process alive until the message is consumed. See [Publish & Subscribe](https://www.pubnub.com/docs/sdks/go/api-reference/publish-and-subscribe). |
| Extra clients or duplicate messages during development | Own one PubNub client at service/application scope instead of creating a new client for each server request or development reload. Unsubscribe during shutdown. |

For setup help, see the [Go SDK documentation](https://www.pubnub.com/docs/sdks/go). Check [network status](https://status.pubnub.com/) for service incidents.

## Releases

Read the [changelog](https://www.pubnub.com/docs/sdks/go/changelog) before upgrading.

For a major version change, review the migration guidance linked from the [Go SDK documentation](https://www.pubnub.com/docs/sdks/go) and the [Go SDK changelog](https://www.pubnub.com/docs/sdks/go/changelog).

## Support and contributions

For setup or account help, contact [PubNub Support](https://support.pubnub.com/).

For a reproducible SDK bug, use [GitHub Issues](https://github.com/pubnub/go/issues). Include the SDK version, runtime, and a small reproduction with credentials removed.

Run `go test ./...` and the repository's lint/build workflow for the affected packages, and include tests for behavioral changes before opening a pull request.

## License

See the [PubNub Software Development Kit License](https://github.com/pubnub/go/blob/master/LICENSE).
