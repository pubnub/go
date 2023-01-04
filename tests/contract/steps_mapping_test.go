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

	ctx.Step(`^I publish message with \'(.*)\' space id and \'(.*)\' message type$`, iPublishMessageWithSpaceIdAndMessageType)
	ctx.Step(`^I receive a successful response$`, iReceiveASuccessfulResponse)
	ctx.Step(`^I receive error response$`, iReceiveErrorResponse)
}
