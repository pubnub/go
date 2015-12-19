package messaging

func StringResponseType(responseType ResponseType) string {
	switch responseType {
	case ChannelResponse:
		return "channel"
	case ChannelGroupResponse:
		return "channel group"
	case WildcardResponse:
		return "wildcard channel"
	default:
		return ""
	}
}

func stringResponseReason(status ResponseStatus) string {
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

func StringConnectionAction(status ConnectionAction) string {
	switch status {
	case ConnectionConnected:
		return "connect"
	case ConnectionUnsubscribed:
		return "unsubscrib"
	case ConnectionReconnected:
		return "reconnect"
	default:
		return ""
	}
}
