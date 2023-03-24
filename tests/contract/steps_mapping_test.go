package contract

import "github.com/cucumber/godog"

func MapSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^grant (.*) permission READ$`, grantPermissionREAD)
	ctx.Step(`^grant (.*) permission MANAGE$`, grantPermissionMANAGE)
	ctx.Step(`^grant (.*) permission DELETE$`, grantPermissionDELETE)
	ctx.Step(`^grant (.*) permission GET$`, grantPermissionGET)
	ctx.Step(`^grant (.*) permission JOIN$`, grantPermissionJOIN)
	ctx.Step(`^grant (.*) permission WRITE$`, grantPermissionWRITE)
	ctx.Step(`^grant (.*) permission UPDATE$`, grantPermissionUPDATE)
	ctx.Step(`^I have a known token .*$`, iHaveAKnownTokenWithEverything)

	ctx.Step(`^I grant a token specifying those permissions$`, iGrantATokenSpecifyingThosePermissions)
	ctx.Step(`^I have a keyset with access manager enabled$`, iHaveAKeysetWithAccessManagerEnabled)
	ctx.Step(`^the \'(.*)\' CHANNEL resource access permissions$`, theCHANNELResourceAccessPermissions)
	ctx.Step(`^the \'(.*)\' CHANNEL pattern access permissions$`, theCHANNELPatternAccessPermissions)
	ctx.Step(`^the \'(.*)\' CHANNEL_GROUP resource access permissions$`, theCHANNEL_GROUPResourceAccessPermissions)
	ctx.Step(`^the \'(.*)\' CHANNEL_GROUP pattern access permissions$`, theCHANNEL_GROUPPatternAccessPermissions)
	ctx.Step(`^the \'(.*)\' UUID resource access permissions$`, theUUIDResourceAccessPermissions)
	ctx.Step(`^the \'(.*)\' UUID pattern access permissions$`, theUUIDPatternAccessPermissions)
	ctx.Step(`^the TTL (\d+)$`, theTTL)
	ctx.Step(`^the token contains the TTL (\d+)$`, theTokenContainsTheTTL)
	ctx.Step(`^the token does not contain an authorized uuid$`, theTokenDoesNotContainAnAuthorizedUuid)
	ctx.Step(`^the token has \'(.*)\' CHANNEL resource access permissions$`, theTokenHasCHANNELResourceAccessPermissions)
	ctx.Step(`^token .* permission (.*)$`, resourceHasPermission)
	ctx.Step(`^the authorized UUID "([^"]*)"$`, theAuthorizedUUID)
	ctx.Step(`^the token contains the authorized UUID "([^"]*)"$`, theTokenContainsTheAuthorizedUUID)
	ctx.Step(`^the parsed token output contains the authorized UUID "(.*)"$`, theTokenContainsTheAuthorizedUUID)

	ctx.Step(`^the token has \'(.*)\' CHANNEL pattern access permissions$`, theTokenHasCHANNELPatternAccessPermissions)
	ctx.Step(`^the token has \'(.*)\' CHANNEL_GROUP resource access permissions$`, theTokenHasCHANNEL_GROUPResourceAccessPermissions)
	ctx.Step(`^the token has \'(.*)\' CHANNEL_GROUP pattern access permissions$`, theTokenHasCHANNEL_GROUPPatternAccessPermissions)
	ctx.Step(`^the token has \'(.*)\' UUID pattern access permissions$`, theTokenHasUUIDPatternAccessPermissions)
	ctx.Step(`^the token has \'(.*)\' UUID resource access permissions$`, theTokenHasUUIDResourceAccessPermissions)
	ctx.Step(`^I parse the token$`, iParseTheToken)
	ctx.Step(`^deny resource permission GET$`, denyResourcePermissionGET)

	ctx.Step(`^the error .* is \'(.*)\'$`, theErrorContains)
	ctx.Step(`^the error status code is (\d+)$`, theErrorStatusCodeIs)
	ctx.Step(`^an error is returned$`, anErrorIsReturned)
	ctx.Step(`^I attempt to grant a token specifying those permissions$`, iAttemptToGrantATokenSpecifyingThosePermissions)

	ctx.Step(`^a token$`, aToken)
	ctx.Step(`^a valid token with permissions to publish with channel \'channel-(\d+)\'$`, aValidTokenWithPermissionsToPublishWithChannelChannel)
	ctx.Step(`^an auth error is returned$`, anErrorIsReturned)
	ctx.Step(`^an expired token with permissions to publish with channel \'channel-(\d+)\'$`, anExpiredTokenWithPermissionsToPublishWithChannelChannel)
	ctx.Step(`^I attempt to publish a message using that auth token with channel \'(.*)\'$`, iPublishAMessageUsingThatAuthTokenWithChannelChannel)
	ctx.Step(`^I get confirmation that token has been revoked$`, iGetConfirmationThatTokenHasBeenRevoked)
	ctx.Step(`^I have a keyset with access manager enabled - without secret key$`, iHaveAKeysetWithAccessManagerEnabledWithoutSecretKey)
	ctx.Step(`^I publish a message using that auth token with channel \'(.*)\'$`, iPublishAMessageUsingThatAuthTokenWithChannelChannel)
	ctx.Step(`^I revoke a token$`, iRevokeAToken)
	ctx.Step(`^the auth error message is \'(.*)\'$`, theErrorContains)
	ctx.Step(`^the error detail message is not empty$`, theErrorDetailMessageIsNotEmpty)
	ctx.Step(`^the result is successful$`, theResultIsSuccessful)
	ctx.Step(`^the token string \'(.*)\'$`, theTokenString)

	ctx.Step(`^the demo keyset$`, theDemoKeyset)
	ctx.Step(`^the demo keyset with enabled storage$`, theDemoKeyset)

	ctx.Step(`^I publish message with \'(.*)\' space id and \'(.*)\' type$`, iPublishMessageWithSpaceIdAndType)
	ctx.Step(`^I send a signal with \'(.*)\' space id and \'(.*)\' type$`, iSendASignalWithSpaceidSpaceIdAndType)
	ctx.Step(`^I receive a successful response$`, iReceiveASuccessfulResponse)
	ctx.Step(`^I receive error response$`, iReceiveErrorResponse)
	ctx.Step(`^I receive an error response$`, iReceiveErrorResponse)

	ctx.Step(`^I receive the message in my subscribe response$`, iReceiveTheMessageInMySubscribeResponse)
	ctx.Step(`^I subscribe to \'(.*)\' channel$`, iSubscribeToChannel)
	ctx.Step(`^subscribe response contains messages with space ids$`, subscribeResponseContainsMessagesWithSpaceIds)
	ctx.Step(`^subscribe response contains messages without space ids$`, subscribeResponseContainsMessagesWithoutSpaceIds)

	ctx.Step(`^history response contains messages with \'(.*)\' and \'(.*)\' types$`, historyResponseContainsMessagesWithProvidedTypes)
	ctx.Step(`^history response contains messages with \'(\d+)\' and \'(\d+)\' message types$`, historyResponseContainsMessagesWithProvidedMessageTypes)
	ctx.Step(`^history response contains messages with space ids$`, historyResponseContainsMessagesWithSpaceIds)
	ctx.Step(`^history response contains messages without types$`, historyResponseContainsMessagesWithoutType)
	ctx.Step(`^history response contains messages without space ids$`, historyResponseContainsMessagesWithoutSpaceIds)
	ctx.Step(`^I fetch message history for \'(.*)\' channel$`, iFetchMessageHistoryForChannel)
	ctx.Step(`^I fetch message history with \'includeType\' set to \'false\' for \'(.*)\' channel$`, iFetchMessageHistoryWithIncludeTypeSetToFalseForChannel)
	ctx.Step(`^I fetch message history with \'includeSpaceId\' set to \'true\' for \'(.*)\' channel$`, iFetchMessageHistoryWithIncludeSpaceIdSetToTrueForChannel)

	ctx.Step(`^I send a file with \'(.*)\' space id and \'(.*)\' message type$`, iSendAFileWithSpaceidAndType)

	ctx.Step(`^I receive (\d+) messages in my subscribe response$`, iReceiveMessagesInMySubscribeResponse)
	ctx.Step(`^response contains messages with \'(.*)\' and \'(.*)\' types$`, responseContainsMessagesWithTypes)
	ctx.Step(`^response contains messages with space ids$`, responseContainsMessagesWithSpaceIds)

}
