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
		return "internet connection issues"
	case reponseAbortMaxRetry:
		return "max retries exceeded"
	case responseTimedOut:
		return "time out"
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
