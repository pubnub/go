package main

import (
	"fmt"
	"time"

	pubnub "github.com/pubnub/go"
)

type Lists struct {
	message_list map[string]interface{}
	deleted_ids  map[string]interface{}
	message_ids  map[string]interface{}
}

var my_channel string

func displayMessages(message interface{}, channel string, lists *Lists) {
	if _, ok := lists.message_list[channel]; !ok {
		lists.message_list[channel] = map[string]interface{}{}
	}

	if _, ok := lists.message_ids[channel]; !ok {
		lists.message_ids[channel] = map[string]interface{}{}
	}

	if _, ok := lists.deleted_ids[channel]; !ok {
		lists.deleted_ids[channel] = map[string]interface{}{}
	}

	messageObj := map[string]interface{}{}

	if m, ok := message.(map[string]interface{}); ok {
		if message_id, ok := m["message_id"].(string); ok {
			if m_ids, ok := lists.message_ids[channel].(map[string]interface{}); ok {
				if _, ok := m_ids[message_id]; !ok {
					lists.message_ids[channel] = message_id
					if m_id, ok := m["message_id"].(string); ok {
						messageObj[m_id] = message
					}
					lists.message_list[channel] = messageObj
				} else {
					if del, ok := m["deleted"].(bool); ok {
						if del {
							lists.deleted_ids[channel] = m["message_id"]
							delete(lists.message_list, channel)
						} else {
							if m_id, ok := m["message_id"].(string); ok {
								messageObj[m_id] = message
							}
							lists.message_list[channel] = messageObj
						}
					}
				}
			}
		}
	}
}

func main() {
	config := pubnub.NewConfig()
	config.PublishKey = "pub-c-1bd448ed-05ba-4dbc-81a5-7d6ff5c6e2bb"
	config.SubscribeKey = "sub-c-b9ab9508-43cf-11e8-9967-869954283fb4"

	pn := pubnub.NewPubNub(config)
	my_channel = "jasdeep-status"
	lists := Lists{
		map[string]interface{}{},
		map[string]interface{}{},
		map[string]interface{}{},
	}

	listener := pubnub.NewListener()

	go func() {
		for {
			select {
			case status := <-listener.Status:
				fmt.Println(status)
			case message := <-listener.Message:
				displayMessages(message.Message, message.Channel, &lists)
			case <-listener.Presence:
			}
		}
	}()

	pn.AddListener(listener)

	pn.Subscribe().
		Channels([]string{my_channel}).
		Execute()

	data := map[string]interface{}{}

	data["message_id"] = "10001"
	data["channel"] = "jasdeep-status"
	data["original_timetoken"] = time.Now().Unix()
	data["user"] = "jasdeep"
	data["status"] = "Writing up design patterns..."
	data["usecase"] = "update"
	data["deleted"] = false
	data["is_update"] = true

	res, status, err := pn.Publish().
		Message(data).
		Channel("jasdeep-status").
		Execute()

	fmt.Println(res, status, err)

	data["is_update"] = false
	data["deleted"] = true

	res, status, err = pn.Publish().
		Message(data).
		Channel("jasdeep-status-UPDATES").
		Execute()

	fmt.Println(res, status, err)

	data["is_update"] = true

	res, status, err = pn.Publish().
		Message(data).
		Channel("jasdeep-status").
		Execute()

	fmt.Println(res, status, err)
}
