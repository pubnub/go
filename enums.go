package pubnub

type PNStatusCategory int

const (
	PNUnknownCategory PNStatusCategory = 1 << iota
	PNAcknowledgmentCategory
	PNAccessDeniedCategory
	PNTimeoutCategory
	PNNetworkIssuesCategory
	PNConnectedCategory
	PNReconnectedCategory
	PNDisconnectedCategory
	PNUnexpectedDisconnectCategory
	PNCancelledCategory
	PNBadRequestCategory
	PNMalformedFilterExpressionCategory
	PNMalformedResponseCategory
	PNDecryptionErrorCategory
	PNTLSConnectionFailedCategory
	PNTLSUntrustedCertificateCategory
	PNRequestMessageCountExceededCategory
)
