package messaging

func stringResponseType(responseType responseType) string {
	switch responseType {
	case channelResponse:
		return "channel"
	case channelGroupResponse:
		return "channel group"
	case wildcardResponse:
		return "wildcard channel"
	default:
		return ""
	}
}

func stringResponseReason(status responseStatus) string {
	switch status {
	case responseAlreadySubscribed:
		return "already subscribed"
	case responseNotSubscribed:
		return "not subscribed"
	case responseInternetConnIssues:
		return "disconnected due to internet connection issues, trying to reconnect."
	case reponseAbortMaxRetry:
		return "aborted due to max retry limit"
	case responseTimedOut:
		return "timed out."
	default:
		return "unknown error"
	}
}

func stringConnectionAction(status connectionAction) string {
	switch status {
	case connectionConnected:
		return "connect"
	case connectionUnsubscribed:
		return "unsubscrib"
	case connectionReconnected:
		return "reconnect"
	default:
		return ""
	}
}

var responseTypeName = map[responseType]string{
	channelResponse:      "channel",
	channelGroupResponse: "channel group",
	wildcardResponse:     "wildcard channel",
}

func (r responseType) String() string {
	return responseTypeName[r]
}

var subscribeLoopActionName = map[subscribeLoopAction]string{
	subscribeLoopStart:     "start",
	subscribeLoopRestart:   "restart",
	subscribeLoopDoNothing: "do nothing",
}

func (s subscribeLoopAction) String() string {
	return subscribeLoopActionName[s]
}
