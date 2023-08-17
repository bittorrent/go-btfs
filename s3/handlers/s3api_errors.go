package handlers

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

// APIError structure
type APIError struct {
	Code           string
	Description    string
	HTTPStatusCode int
}

// RESTErrorResponse - error response format
type RESTErrorResponse struct {
	XMLName    xml.Name `xml:"Error" json:"-"`
	Code       string   `xml:"Code" json:"Code"`
	Message    string   `xml:"Message" json:"Message"`
	Resource   string   `xml:"Resource" json:"Resource"`
	RequestID  string   `xml:"RequestId" json:"RequestId"`
	Key        string   `xml:"Key,omitempty" json:"Key,omitempty"`
	BucketName string   `xml:"BucketName,omitempty" json:"BucketName,omitempty"`
}

// Error - Returns S3 error string.
func (e RESTErrorResponse) Error() string {
	if e.Message == "" {
		msg, ok := s3ErrorResponseMap[e.Code]
		if !ok {
			msg = fmt.Sprintf("Error response code %s.", e.Code)
		}
		return msg
	}
	return e.Message
}

// ErrorCode type of error status.
type ErrorCode int

// Error codes, non exhaustive list - http://docs.aws.amazon.com/AmazonS3/latest/API/ErrorResponses.html
const (
	ErrCodeNone ErrorCode = iota
	ErrCodeAccessDenied
	ErrCodeBadDigest
	ErrCodeEntityTooSmall
	ErrCodeEntityTooLarge
	ErrCodeIncompleteBody
	ErrCodeInternalError
	ErrCodeInvalidAccessKeyID
	ErrCodeAccessKeyDisabled
	ErrCodeInvalidBucketName
	ErrCodeInvalidDigest
	ErrCodeInvalidRange
	ErrCodeInvalidRangePartNumber
	ErrCodeInvalidCopyPartRange
	ErrCodeInvalidCopyPartRangeSource
	ErrCodeInvalidMaxKeys
	ErrCodeInvalidEncodingMethod
	ErrCodeInvalidMaxUploads
	ErrCodeInvalidMaxParts
	ErrCodeInvalidPartNumberMarker
	ErrCodeInvalidRequestBody
	ErrCodeInvalidCopySource
	ErrCodeInvalidMetadataDirective
	ErrCodeCodeInvalidCopyDest
	ErrCodeInvalidPolicyDocument
	ErrCodeInvalidObjectState
	ErrCodeMalformedXML
	ErrCodeMissingContentLength
	ErrCodeMissingContentMD5
	ErrCodeMissingRequestBodyError
	ErrCodeMissingSecurityHeader
	ErrCodeNoSuchUser
	ErrCodeUserAlreadyExists
	ErrCodeNoSuchUserPolicy
	ErrCodeUserPolicyAlreadyExists
	ErrCodeNoSuchBucket
	ErrCodeNoSuchBucketPolicy
	ErrCodeNoSuchLifecycleConfiguration
	ErrCodeNoSuchCORSConfiguration
	ErrCodeNoSuchWebsiteConfiguration
	ErrCodeReplicationConfigurationNotFoundError
	ErrCodeReplicationNeedsVersioningError
	ErrCodeReplicationBucketNeedsVersioningError
	ErrCodeObjectRestoreAlreadyInProgress
	ErrCodeNoSuchKey
	ErrCodeNoSuchUpload
	ErrCodeInvalidVersionID
	ErrCodeNoSuchVersion
	ErrCodeNotImplemented
	ErrCodePreconditionFailed
	ErrCodeRequestTimeTooSkewed
	ErrCodeSignatureDoesNotMatch
	ErrCodeMethodNotAllowed
	ErrCodeInvalidPart
	ErrCodeInvalidPartOrder
	ErrCodeAuthorizationHeaderMalformed
	ErrCodeMalformedDate
	ErrCodeMalformedPOSTRequest
	ErrCodePOSTFileRequired
	ErrCodeSignatureVersionNotSupported
	ErrCodeBucketNotEmpty
	ErrCodeAllAccessDisabled
	ErrCodeMalformedPolicy
	ErrCodeMissingFields
	ErrCodeMissingCredTag
	ErrCodeCredMalformed
	ErrCodeInvalidRegion

	ErrCodeMissingSignTag
	ErrCodeMissingSignHeadersTag

	ErrCodeAuthHeaderEmpty
	ErrCodeExpiredPresignRequest
	ErrCodeRequestNotReadyYet
	ErrCodeUnsignedHeaders
	ErrCodeMissingDateHeader

	ErrCodeBucketAlreadyOwnedByYou
	ErrCodeInvalidDuration
	ErrCodeBucketAlreadyExists
	ErrCodeMetadataTooLarge
	ErrCodeUnsupportedMetadata

	ErrCodeSlowDown
	ErrCodeBadRequest
	ErrCodeKeyTooLongError
	ErrCodeInvalidBucketObjectLockConfiguration
	ErrCodeObjectLockConfigurationNotAllowed
	ErrCodeNoSuchObjectLockConfiguration
	ErrCodeObjectLocked
	ErrCodeInvalidRetentionDate
	ErrCodePastObjectLockRetainDate
	ErrCodeUnknownWORMModeDirective
	ErrCodeBucketTaggingNotFound
	ErrCodeObjectLockInvalidHeaders
	ErrCodeInvalidTagDirective
	// Add new error codes here.

	// SSE-S3 related API errors
	ErrCodeInvalidEncryptionMethod
	ErrCodeInvalidQueryParams
	ErrCodeNoAccessKey
	ErrCodeInvalidToken

	// Bucket notification related errors.
	ErrCodeEventNotification
	ErrCodeARNNotification
	ErrCodeRegionNotification
	ErrCodeOverlappingFilterNotification
	ErrCodeFilterNameInvalid
	ErrCodeFilterNamePrefix
	ErrCodeFilterNameSuffix
	ErrCodeFilterValueInvalid
	ErrCodeOverlappingConfigs

	// S3 extended errors.
	ErrCodeContentSHA256Mismatch

	// Add new extended error codes here.
	ErrCodeInvalidObjectName
	ErrCodeInvalidObjectNamePrefixSlash
	ErrCodeClientDisconnected
	ErrCodeOperationTimedOut
	ErrCodeOperationMaxedOut
	ErrCodeInvalidRequest
	ErrCodeIncorrectContinuationToken
	ErrCodeInvalidFormatAccessKey

	// S3 Select Errors
	ErrCodeEmptyRequestBody
	ErrCodeUnsupportedFunction
	ErrCodeInvalidExpressionType
	ErrCodeBusy
	ErrCodeUnauthorizedAccess
	ErrCodeExpressionTooLong
	ErrCodeIllegalSQLFunctionArgument
	ErrCodeInvalidKeyPath
	ErrCodeInvalidCompressionFormat
	ErrCodeInvalidFileHeaderInfo
	ErrCodeInvalidJSONType
	ErrCodeInvalidQuoteFields
	ErrCodeInvalidRequestParameter
	ErrCodeInvalidDataType
	ErrCodeInvalidTextEncoding
	ErrCodeInvalidDataSource
	ErrCodeInvalidTableAlias
	ErrCodeMissingRequiredParameter
	ErrCodeObjectSerializationConflict
	ErrCodeUnsupportedSQLOperation
	ErrCodeUnsupportedSQLStructure
	ErrCodeUnsupportedSyntax
	ErrCodeUnsupportedRangeHeader
	ErrCodeLexerInvalidChar
	ErrCodeLexerInvalidOperator
	ErrCodeLexerInvalidLiteral
	ErrCodeLexerInvalidIONLiteral
	ErrCodeParseExpectedDatePart
	ErrCodeParseExpectedKeyword
	ErrCodeParseExpectedTokenType
	ErrCodeParseExpected2TokenTypes
	ErrCodeParseExpectedNumber
	ErrCodeParseExpectedRightParenBuiltinFunctionCall
	ErrCodeParseExpectedTypeName
	ErrCodeParseExpectedWhenClause
	ErrCodeParseUnsupportedToken
	ErrCodeParseUnsupportedLiteralsGroupBy
	ErrCodeParseExpectedMember
	ErrCodeParseUnsupportedSelect
	ErrCodeParseUnsupportedCase
	ErrCodeParseUnsupportedCaseClause
	ErrCodeParseUnsupportedAlias
	ErrCodeParseUnsupportedSyntax
	ErrCodeParseUnknownOperator
	ErrCodeParseMissingIdentAfterAt
	ErrCodeParseUnexpectedOperator
	ErrCodeParseUnexpectedTerm
	ErrCodeParseUnexpectedToken
	ErrCodeParseUnexpectedKeyword
	ErrCodeParseExpectedExpression
	ErrCodeParseExpectedLeftParenAfterCast
	ErrCodeParseExpectedLeftParenValueConstructor
	ErrCodeParseExpectedLeftParenBuiltinFunctionCall
	ErrCodeParseExpectedArgumentDelimiter
	ErrCodeParseCastArity
	ErrCodeParseInvalidTypeParam
	ErrCodeParseEmptySelect
	ErrCodeParseSelectMissingFrom
	ErrCodeParseExpectedIdentForGroupName
	ErrCodeParseExpectedIdentForAlias
	ErrCodeParseUnsupportedCallWithStar
	ErrCodeParseNonUnaryAgregateFunctionCall
	ErrCodeParseMalformedJoin
	ErrCodeParseExpectedIdentForAt
	ErrCodeParseAsteriskIsNotAloneInSelectList
	ErrCodeParseCannotMixSqbAndWildcardInSelectList
	ErrCodeParseInvalidContextForWildcardInSelectList
	ErrCodeIncorrectSQLFunctionArgumentType
	ErrCodeValueParseFailure
	ErrCodeEvaluatorInvalidArguments
	ErrCodeIntegerOverflow
	ErrCodeLikeInvalidInputs
	ErrCodeCastFailed
	ErrCodeInvalidCast
	ErrCodeEvaluatorInvalidTimestampFormatPattern
	ErrCodeEvaluatorInvalidTimestampFormatPatternSymbolForParsing
	ErrCodeEvaluatorTimestampFormatPatternDuplicateFields
	ErrCodeEvaluatorTimestampFormatPatternHourClockAmPmMismatch
	ErrCodeEvaluatorUnterminatedTimestampFormatPatternToken
	ErrCodeEvaluatorInvalidTimestampFormatPatternToken
	ErrCodeEvaluatorInvalidTimestampFormatPatternSymbol
	ErrCodeEvaluatorBindingDoesNotExist
	ErrCodeMissingHeaders
	ErrCodeInvalidColumnIndex

	ErrCodePostPolicyConditionInvalidFormat

	ErrCodeMalformedJSON
)

// error code to APIError structure, these fields carry respective
// descriptions for all the error responses.
var errorCodeResponse = map[ErrorCode]APIError{
	ErrCodeCodeInvalidCopyDest: {
		Code:           "InvalidRequest",
		Description:    "This copy request is illegal because it is trying to copy an object to itself without changing the object's metadata, storage class, website redirect location or encryption attributes.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidCopySource: {
		Code:           "InvalidArgument",
		Description:    "Copy Source must mention the source bucket and key: sourcebucket/sourcekey.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidMetadataDirective: {
		Code:           "InvalidArgument",
		Description:    "Unknown metadata directive.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidRequestBody: {
		Code:           "InvalidArgument",
		Description:    "Body shouldn't be set for this request.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidMaxUploads: {
		Code:           "InvalidArgument",
		Description:    "Argument max-uploads must be an integer between 0 and 2147483647",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidMaxKeys: {
		Code:           "InvalidArgument",
		Description:    "Argument maxKeys must be an integer between 0 and 2147483647",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidEncodingMethod: {
		Code:           "InvalidArgument",
		Description:    "Invalid Encoding Method specified in Request",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidMaxParts: {
		Code:           "InvalidArgument",
		Description:    "Part number must be an integer between 1 and 10000, inclusive",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidPartNumberMarker: {
		Code:           "InvalidArgument",
		Description:    "Argument partNumberMarker must be an integer.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidPolicyDocument: {
		Code:           "InvalidPolicyDocument",
		Description:    "The content of the form does not meet the conditions specified in the policy document.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeAccessDenied: {
		Code:           "AccessDenied",
		Description:    "Access Denied.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrCodeBadDigest: {
		Code:           "BadDigest",
		Description:    "The Content-Md5 you specified did not match what we received.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeEntityTooSmall: {
		Code:           "EntityTooSmall",
		Description:    "Your proposed upload is smaller than the minimum allowed object size.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeEntityTooLarge: {
		Code:           "EntityTooLarge",
		Description:    "Your proposed upload exceeds the maximum allowed object size.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeIncompleteBody: {
		Code:           "IncompleteBody",
		Description:    "You did not provide the number of bytes specified by the Content-Length HTTP header.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInternalError: {
		Code:           "InternalError",
		Description:    "We encountered an internal error, please try again.",
		HTTPStatusCode: http.StatusInternalServerError,
	},
	ErrCodeInvalidAccessKeyID: {
		Code:           "InvalidAccessKeyId",
		Description:    "The Access Key Id you provided does not exist in our records.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrCodeAccessKeyDisabled: {
		Code:           "InvalidAccessKeyId",
		Description:    "Your account is disabled; please contact your administrator.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrCodeInvalidBucketName: {
		Code:           "InvalidBucketName",
		Description:    "The specified bucket is not valid.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidDigest: {
		Code:           "InvalidDigest",
		Description:    "The Content-Md5 you specified is not valid.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidRange: {
		Code:           "InvalidRange",
		Description:    "The requested range is not satisfiable",
		HTTPStatusCode: http.StatusRequestedRangeNotSatisfiable,
	},
	ErrCodeInvalidRangePartNumber: {
		Code:           "InvalidRequest",
		Description:    "Cannot specify both Range header and partNumber query parameter",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeMalformedXML: {
		Code:           "MalformedXML",
		Description:    "The XML you provided was not well-formed or did not validate against our published schema.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeMissingContentLength: {
		Code:           "MissingContentLength",
		Description:    "You must provide the Content-Length HTTP header.",
		HTTPStatusCode: http.StatusLengthRequired,
	},
	ErrCodeMissingContentMD5: {
		Code:           "MissingContentMD5",
		Description:    "Missing required header for this request: Content-Md5.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeMissingSecurityHeader: {
		Code:           "MissingSecurityHeader",
		Description:    "Your request was missing a required header",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeMissingRequestBodyError: {
		Code:           "MissingRequestBodyError",
		Description:    "Request body is empty.",
		HTTPStatusCode: http.StatusLengthRequired,
	},
	ErrCodeNoSuchBucket: {
		Code:           "NoSuchBucket",
		Description:    "The specified bucket does not exist",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrCodeNoSuchBucketPolicy: {
		Code:           "NoSuchBucketPolicy",
		Description:    "The bucket policy does not exist",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrCodeNoSuchLifecycleConfiguration: {
		Code:           "NoSuchLifecycleConfiguration",
		Description:    "The lifecycle configuration does not exist",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrCodeNoSuchUser: {
		Code:           "NoSuchUser",
		Description:    "The specified user does not exist",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrCodeUserAlreadyExists: {
		Code:           "UserAlreadyExists",
		Description:    "The request was rejected because it attempted to create a resource that already exists .",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrCodeNoSuchUserPolicy: {
		Code:           "NoSuchUserPolicy",
		Description:    "The specified user policy does not exist",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrCodeUserPolicyAlreadyExists: {
		Code:           "UserPolicyAlreadyExists",
		Description:    "The same user policy already exists .",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrCodeNoSuchKey: {
		Code:           "NoSuchKey",
		Description:    "The specified key does not exist.",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrCodeNoSuchUpload: {
		Code:           "NoSuchUpload",
		Description:    "The specified multipart upload does not exist. The upload ID may be invalid, or the upload may have been aborted or completed.",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrCodeInvalidVersionID: {
		Code:           "InvalidArgument",
		Description:    "Invalid version id specified",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeNoSuchVersion: {
		Code:           "NoSuchVersion",
		Description:    "The specified version does not exist.",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrCodeNotImplemented: {
		Code:           "NotImplemented",
		Description:    "A header you provided implies functionality that is not implemented",
		HTTPStatusCode: http.StatusNotImplemented,
	},
	ErrCodePreconditionFailed: {
		Code:           "PreconditionFailed",
		Description:    "At least one of the pre-conditions you specified did not hold",
		HTTPStatusCode: http.StatusPreconditionFailed,
	},
	ErrCodeRequestTimeTooSkewed: {
		Code:           "RequestTimeTooSkewed",
		Description:    "The difference between the request time and the server's time is too large.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrCodeSignatureDoesNotMatch: {
		Code:           "SignatureDoesNotMatch",
		Description:    "The request signature we calculated does not match the signature you provided. Check your key and signing method.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrCodeMethodNotAllowed: {
		Code:           "MethodNotAllowed",
		Description:    "The specified method is not allowed against this resource.",
		HTTPStatusCode: http.StatusMethodNotAllowed,
	},
	ErrCodeInvalidPart: {
		Code:           "InvalidPart",
		Description:    "One or more of the specified parts could not be found.  The part may not have been uploaded, or the specified entity tag may not match the part's entity tag.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidPartOrder: {
		Code:           "InvalidPartOrder",
		Description:    "The list of parts was not in ascending order. The parts list must be specified in order by part number.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidObjectState: {
		Code:           "InvalidObjectState",
		Description:    "The operation is not valid for the current state of the object.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrCodeAuthorizationHeaderMalformed: {
		Code:           "AuthorizationHeaderMalformed",
		Description:    "The authorization header is malformed; the region is wrong; expecting 'us-east-1'.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeMalformedPOSTRequest: {
		Code:           "MalformedPOSTRequest",
		Description:    "The body of your POST request is not well-formed multipart/form-data.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodePOSTFileRequired: {
		Code:           "InvalidArgument",
		Description:    "POST requires exactly one file upload per request.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeSignatureVersionNotSupported: {
		Code:           "InvalidRequest",
		Description:    "The authorization mechanism you have provided is not supported. Please use AWS4-HMAC-SHA256.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeBucketNotEmpty: {
		Code:           "BucketNotEmpty",
		Description:    "The bucket you tried to delete is not empty",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrCodeBucketAlreadyExists: {
		Code:           "BucketAlreadyExists",
		Description:    "The requested bucket name is not available. The bucket namespace is shared by all users of the system. Please select a different name and try again.",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrCodeAllAccessDisabled: {
		Code:           "AllAccessDisabled",
		Description:    "All access to this resource has been disabled.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrCodeMalformedPolicy: {
		Code:           "MalformedPolicy",
		Description:    "Policy has invalid resource.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeMissingCredTag: {
		Code:           "InvalidRequest",
		Description:    "Missing Credential field for this request.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidRegion: {
		Code:           "InvalidRegion",
		Description:    "Region does not match.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeMissingSignTag: {
		Code:           "AccessDenied",
		Description:    "Signature header missing Signature field.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeMissingSignHeadersTag: {
		Code:           "InvalidArgument",
		Description:    "Signature header missing SignedHeaders field.",
		HTTPStatusCode: http.StatusBadRequest,
	},

	ErrCodeAuthHeaderEmpty: {
		Code:           "InvalidArgument",
		Description:    "Authorization header is invalid -- one and only one ' ' (space) required.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeMissingDateHeader: {
		Code:           "AccessDenied",
		Description:    "AWS authentication requires a valid Date or x-amz-date header",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeExpiredPresignRequest: {
		Code:           "AccessDenied",
		Description:    "Request has expired",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrCodeRequestNotReadyYet: {
		Code:           "AccessDenied",
		Description:    "Request is not valid yet",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrCodeSlowDown: {
		Code:           "SlowDown",
		Description:    "Resource requested is unreadable, please reduce your request rate",
		HTTPStatusCode: http.StatusServiceUnavailable,
	},
	ErrCodeBadRequest: {
		Code:           "BadRequest",
		Description:    "400 BadRequest",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeKeyTooLongError: {
		Code:           "KeyTooLongError",
		Description:    "Your key is too long",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeUnsignedHeaders: {
		Code:           "AccessDenied",
		Description:    "There were headers present in the request which were not signed",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeBucketAlreadyOwnedByYou: {
		Code:           "BucketAlreadyOwnedByYou",
		Description:    "Your previous request to create the named bucket succeeded and you already own it.",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrCodeInvalidDuration: {
		Code:           "InvalidDuration",
		Description:    "Duration provided in the request is invalid.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidBucketObjectLockConfiguration: {
		Code:           "InvalidRequest",
		Description:    "Bucket is missing ObjectLockConfiguration",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeBucketTaggingNotFound: {
		Code:           "NoSuchTagSet",
		Description:    "The TagSet does not exist",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrCodeObjectLockConfigurationNotAllowed: {
		Code:           "InvalidBucketState",
		Description:    "Object Lock configuration cannot be enabled on existing buckets",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrCodeNoSuchCORSConfiguration: {
		Code:           "NoSuchCORSConfiguration",
		Description:    "The CORS configuration does not exist",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrCodeNoSuchWebsiteConfiguration: {
		Code:           "NoSuchWebsiteConfiguration",
		Description:    "The specified bucket does not have a website configuration",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrCodeReplicationConfigurationNotFoundError: {
		Code:           "ReplicationConfigurationNotFoundError",
		Description:    "The replication configuration was not found",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrCodeReplicationNeedsVersioningError: {
		Code:           "InvalidRequest",
		Description:    "Versioning must be 'Enabled' on the bucket to apply a replication configuration",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeReplicationBucketNeedsVersioningError: {
		Code:           "InvalidRequest",
		Description:    "Versioning must be 'Enabled' on the bucket to add a replication target",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeNoSuchObjectLockConfiguration: {
		Code:           "NoSuchObjectLockConfiguration",
		Description:    "The specified object does not have a ObjectLock configuration",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeObjectLocked: {
		Code:           "InvalidRequest",
		Description:    "Object is WORM protected and cannot be overwritten",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidRetentionDate: {
		Code:           "InvalidRequest",
		Description:    "Date must be provided in ISO 8601 format",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodePastObjectLockRetainDate: {
		Code:           "InvalidRequest",
		Description:    "the retain until date must be in the future",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeUnknownWORMModeDirective: {
		Code:           "InvalidRequest",
		Description:    "unknown wormMode directive",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeObjectLockInvalidHeaders: {
		Code:           "InvalidRequest",
		Description:    "x-amz-object-lock-retain-until-date and x-amz-object-lock-mode must both be supplied",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeObjectRestoreAlreadyInProgress: {
		Code:           "RestoreAlreadyInProgress",
		Description:    "Object restore is already in progress",
		HTTPStatusCode: http.StatusConflict,
	},
	// Bucket notification related errors.
	ErrCodeEventNotification: {
		Code:           "InvalidArgument",
		Description:    "A specified event is not supported for notifications.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeARNNotification: {
		Code:           "InvalidArgument",
		Description:    "A specified destination ARN does not exist or is not well-formed. Verify the destination ARN.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeRegionNotification: {
		Code:           "InvalidArgument",
		Description:    "A specified destination is in a different region than the bucket. You must use a destination that resides in the same region as the bucket.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeOverlappingFilterNotification: {
		Code:           "InvalidArgument",
		Description:    "An object key name filtering rule defined with overlapping prefixes, overlapping suffixes, or overlapping combinations of prefixes and suffixes for the same event types.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeFilterNameInvalid: {
		Code:           "InvalidArgument",
		Description:    "filter rule name must be either prefix or suffix",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeFilterNamePrefix: {
		Code:           "InvalidArgument",
		Description:    "Cannot specify more than one prefix rule in a filter.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeFilterNameSuffix: {
		Code:           "InvalidArgument",
		Description:    "Cannot specify more than one suffix rule in a filter.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeFilterValueInvalid: {
		Code:           "InvalidArgument",
		Description:    "Size of filter rule value cannot exceed 1024 bytes in UTF-8 representation",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeOverlappingConfigs: {
		Code:           "InvalidArgument",
		Description:    "Configurations overlap. Configurations on the same bucket cannot share a common event type.",
		HTTPStatusCode: http.StatusBadRequest,
	},

	ErrCodeInvalidCopyPartRange: {
		Code:           "InvalidArgument",
		Description:    "The x-amz-copy-source-range value must be of the form bytes=first-last where first and last are the zero-based offsets of the first and last bytes to copy",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidCopyPartRangeSource: {
		Code:           "InvalidArgument",
		Description:    "Range specified is not valid for source object",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeMetadataTooLarge: {
		Code:           "MetadataTooLarge",
		Description:    "Your metadata headers exceed the maximum allowed metadata size.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidTagDirective: {
		Code:           "InvalidArgument",
		Description:    "Unknown tag directive.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidEncryptionMethod: {
		Code:           "InvalidRequest",
		Description:    "The encryption method specified is not supported",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidQueryParams: {
		Code:           "AuthorizationQueryParametersError",
		Description:    "Query-string authentication version 4 requires the X-Amz-Algorithm, X-Amz-Credential, X-Amz-Signature, X-Amz-Date, X-Amz-SignedHeaders, and X-Amz-Expires parameters.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeNoAccessKey: {
		Code:           "AccessDenied",
		Description:    "No AWSAccessKey was presented",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrCodeInvalidToken: {
		Code:           "InvalidTokenId",
		Description:    "The security token included in the request is invalid",
		HTTPStatusCode: http.StatusForbidden,
	},

	// S3 extensions.
	ErrCodeInvalidObjectName: {
		Code:           "InvalidObjectName",
		Description:    "Object name contains unsupported characters.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidObjectNamePrefixSlash: {
		Code:           "InvalidObjectName",
		Description:    "Object name contains a leading slash.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeClientDisconnected: {
		Code:           "ClientDisconnected",
		Description:    "Client disconnected before response was ready",
		HTTPStatusCode: 499, // No official code, use nginx value.
	},
	ErrCodeOperationTimedOut: {
		Code:           "RequestTimeout",
		Description:    "A timeout occurred while trying to lock a resource, please reduce your request rate",
		HTTPStatusCode: http.StatusServiceUnavailable,
	},
	ErrCodeOperationMaxedOut: {
		Code:           "SlowDown",
		Description:    "A timeout exceeded while waiting to proceed with the request, please reduce your request rate",
		HTTPStatusCode: http.StatusServiceUnavailable,
	},
	ErrCodeUnsupportedMetadata: {
		Code:           "InvalidArgument",
		Description:    "Your metadata headers are not supported.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	// Generic Invalid-Request error. Should be used for response errors only for unlikely
	// corner case errors for which introducing new APIErrCodeorCode is not worth it. LogIf()
	// should be used to log the error at the source of the error for debugging purposes.
	ErrCodeInvalidRequest: {
		Code:           "InvalidRequest",
		Description:    "Invalid Request",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeIncorrectContinuationToken: {
		Code:           "InvalidArgument",
		Description:    "The continuation token provided is incorrect",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidFormatAccessKey: {
		Code:           "InvalidAccessKeyId",
		Description:    "The Access Key Id you provided contains invalid characters.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	// S3 Select API Errors
	ErrCodeEmptyRequestBody: {
		Code:           "EmptyRequestBody",
		Description:    "Request body cannot be empty.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeUnsupportedFunction: {
		Code:           "UnsupportedFunction",
		Description:    "Encountered an unsupported SQL function.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidDataSource: {
		Code:           "InvalidDataSource",
		Description:    "Invalid data source type. Only CSV and JSON are supported at this time.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidExpressionType: {
		Code:           "InvalidExpressionType",
		Description:    "The ExpressionType is invalid. Only SQL expressions are supported at this time.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeBusy: {
		Code:           "Busy",
		Description:    "The service is unavailable. Please retry.",
		HTTPStatusCode: http.StatusServiceUnavailable,
	},
	ErrCodeUnauthorizedAccess: {
		Code:           "UnauthorizedAccess",
		Description:    "You are not authorized to perform this operation",
		HTTPStatusCode: http.StatusUnauthorized,
	},
	ErrCodeExpressionTooLong: {
		Code:           "ExpressionTooLong",
		Description:    "The SQL expression is too long: The maximum byte-length for the SQL expression is 256 KB.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeIllegalSQLFunctionArgument: {
		Code:           "IllegalSqlFunctionArgument",
		Description:    "Illegal argument was used in the SQL function.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidKeyPath: {
		Code:           "InvalidKeyPath",
		Description:    "Key path in the SQL expression is invalid.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidCompressionFormat: {
		Code:           "InvalidCompressionFormat",
		Description:    "The file is not in a supported compression format. Only GZIP is supported at this time.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidFileHeaderInfo: {
		Code:           "InvalidFileHeaderInfo",
		Description:    "The FileHeaderInfo is invalid. Only NONE, USE, and IGNORE are supported.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidJSONType: {
		Code:           "InvalidJsonType",
		Description:    "The JsonType is invalid. Only DOCUMENT and LINES are supported at this time.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidQuoteFields: {
		Code:           "InvalidQuoteFields",
		Description:    "The QuoteFields is invalid. Only ALWAYS and ASNEEDED are supported.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidRequestParameter: {
		Code:           "InvalidRequestParameter",
		Description:    "The value of a parameter in SelectRequest element is invalid. Check the service API documentation and try again.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidDataType: {
		Code:           "InvalidDataType",
		Description:    "The SQL expression contains an invalid data type.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidTextEncoding: {
		Code:           "InvalidTextEncoding",
		Description:    "Invalid encoding type. Only UTF-8 encoding is supported at this time.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidTableAlias: {
		Code:           "InvalidTableAlias",
		Description:    "The SQL expression contains an invalid table alias.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeMissingRequiredParameter: {
		Code:           "MissingRequiredParameter",
		Description:    "The SelectRequest entity is missing a required parameter. Check the service documentation and try again.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeObjectSerializationConflict: {
		Code:           "ObjectSerializationConflict",
		Description:    "The SelectRequest entity can only contain one of CSV or JSON. Check the service documentation and try again.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeUnsupportedSQLOperation: {
		Code:           "UnsupportedSqlOperation",
		Description:    "Encountered an unsupported SQL operation.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeUnsupportedSQLStructure: {
		Code:           "UnsupportedSqlStructure",
		Description:    "Encountered an unsupported SQL structure. Check the SQL Reference.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeUnsupportedSyntax: {
		Code:           "UnsupportedSyntax",
		Description:    "Encountered invalid syntax.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeUnsupportedRangeHeader: {
		Code:           "UnsupportedRangeHeader",
		Description:    "Range header is not supported for this operation.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeLexerInvalidChar: {
		Code:           "LexerInvalidChar",
		Description:    "The SQL expression contains an invalid character.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeLexerInvalidOperator: {
		Code:           "LexerInvalidOperator",
		Description:    "The SQL expression contains an invalid literal.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeLexerInvalidLiteral: {
		Code:           "LexerInvalidLiteral",
		Description:    "The SQL expression contains an invalid operator.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeLexerInvalidIONLiteral: {
		Code:           "LexerInvalidIONLiteral",
		Description:    "The SQL expression contains an invalid operator.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseExpectedDatePart: {
		Code:           "ParseExpectedDatePart",
		Description:    "Did not find the expected date part in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseExpectedKeyword: {
		Code:           "ParseExpectedKeyword",
		Description:    "Did not find the expected keyword in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseExpectedTokenType: {
		Code:           "ParseExpectedTokenType",
		Description:    "Did not find the expected token in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseExpected2TokenTypes: {
		Code:           "ParseExpected2TokenTypes",
		Description:    "Did not find the expected token in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseExpectedNumber: {
		Code:           "ParseExpectedNumber",
		Description:    "Did not find the expected number in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseExpectedRightParenBuiltinFunctionCall: {
		Code:           "ParseExpectedRightParenBuiltinFunctionCall",
		Description:    "Did not find the expected right parenthesis character in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseExpectedTypeName: {
		Code:           "ParseExpectedTypeName",
		Description:    "Did not find the expected type name in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseExpectedWhenClause: {
		Code:           "ParseExpectedWhenClause",
		Description:    "Did not find the expected WHEN clause in the SQL expression. CASE is not supported.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseUnsupportedToken: {
		Code:           "ParseUnsupportedToken",
		Description:    "The SQL expression contains an unsupported token.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseUnsupportedLiteralsGroupBy: {
		Code:           "ParseUnsupportedLiteralsGroupBy",
		Description:    "The SQL expression contains an unsupported use of GROUP BY.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseExpectedMember: {
		Code:           "ParseExpectedMember",
		Description:    "The SQL expression contains an unsupported use of MEMBER.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseUnsupportedSelect: {
		Code:           "ParseUnsupportedSelect",
		Description:    "The SQL expression contains an unsupported use of SELECT.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseUnsupportedCase: {
		Code:           "ParseUnsupportedCase",
		Description:    "The SQL expression contains an unsupported use of CASE.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseUnsupportedCaseClause: {
		Code:           "ParseUnsupportedCaseClause",
		Description:    "The SQL expression contains an unsupported use of CASE.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseUnsupportedAlias: {
		Code:           "ParseUnsupportedAlias",
		Description:    "The SQL expression contains an unsupported use of ALIAS.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseUnsupportedSyntax: {
		Code:           "ParseUnsupportedSyntax",
		Description:    "The SQL expression contains unsupported syntax.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseUnknownOperator: {
		Code:           "ParseUnknownOperator",
		Description:    "The SQL expression contains an invalid operator.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseMissingIdentAfterAt: {
		Code:           "ParseMissingIdentAfterAt",
		Description:    "Did not find the expected identifier after the @ symbol in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseUnexpectedOperator: {
		Code:           "ParseUnexpectedOperator",
		Description:    "The SQL expression contains an unexpected operator.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseUnexpectedTerm: {
		Code:           "ParseUnexpectedTerm",
		Description:    "The SQL expression contains an unexpected term.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseUnexpectedToken: {
		Code:           "ParseUnexpectedToken",
		Description:    "The SQL expression contains an unexpected token.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseUnexpectedKeyword: {
		Code:           "ParseUnexpectedKeyword",
		Description:    "The SQL expression contains an unexpected keyword.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseExpectedExpression: {
		Code:           "ParseExpectedExpression",
		Description:    "Did not find the expected SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseExpectedLeftParenAfterCast: {
		Code:           "ParseExpectedLeftParenAfterCast",
		Description:    "Did not find expected the left parenthesis in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseExpectedLeftParenValueConstructor: {
		Code:           "ParseExpectedLeftParenValueConstructor",
		Description:    "Did not find expected the left parenthesis in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseExpectedLeftParenBuiltinFunctionCall: {
		Code:           "ParseExpectedLeftParenBuiltinFunctionCall",
		Description:    "Did not find the expected left parenthesis in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseExpectedArgumentDelimiter: {
		Code:           "ParseExpectedArgumentDelimiter",
		Description:    "Did not find the expected argument delimiter in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseCastArity: {
		Code:           "ParseCastArity",
		Description:    "The SQL expression CAST has incorrect arity.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseInvalidTypeParam: {
		Code:           "ParseInvalidTypeParam",
		Description:    "The SQL expression contains an invalid parameter value.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseEmptySelect: {
		Code:           "ParseEmptySelect",
		Description:    "The SQL expression contains an empty SELECT.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseSelectMissingFrom: {
		Code:           "ParseSelectMissingFrom",
		Description:    "GROUP is not supported in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseExpectedIdentForGroupName: {
		Code:           "ParseExpectedIdentForGroupName",
		Description:    "GROUP is not supported in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseExpectedIdentForAlias: {
		Code:           "ParseExpectedIdentForAlias",
		Description:    "Did not find the expected identifier for the alias in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseUnsupportedCallWithStar: {
		Code:           "ParseUnsupportedCallWithStar",
		Description:    "Only COUNT with (*) as a parameter is supported in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseNonUnaryAgregateFunctionCall: {
		Code:           "ParseNonUnaryAgregateFunctionCall",
		Description:    "Only one argument is supported for aggregate functions in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseMalformedJoin: {
		Code:           "ParseMalformedJoin",
		Description:    "JOIN is not supported in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseExpectedIdentForAt: {
		Code:           "ParseExpectedIdentForAt",
		Description:    "Did not find the expected identifier for AT name in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseAsteriskIsNotAloneInSelectList: {
		Code:           "ParseAsteriskIsNotAloneInSelectList",
		Description:    "Other expressions are not allowed in the SELECT list when '*' is used without dot notation in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseCannotMixSqbAndWildcardInSelectList: {
		Code:           "ParseCannotMixSqbAndWildcardInSelectList",
		Description:    "Cannot mix [] and * in the same expression in a SELECT list in SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeParseInvalidContextForWildcardInSelectList: {
		Code:           "ParseInvalidContextForWildcardInSelectList",
		Description:    "Invalid use of * in SELECT list in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeIncorrectSQLFunctionArgumentType: {
		Code:           "IncorrectSqlFunctionArgumentType",
		Description:    "Incorrect type of arguments in function call in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeValueParseFailure: {
		Code:           "ValueParseFailure",
		Description:    "Time stamp parse failure in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeEvaluatorInvalidArguments: {
		Code:           "EvaluatorInvalidArguments",
		Description:    "Incorrect number of arguments in the function call in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeIntegerOverflow: {
		Code:           "IntegerOverflow",
		Description:    "Int overflow or underflow in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeLikeInvalidInputs: {
		Code:           "LikeInvalidInputs",
		Description:    "Invalid argument given to the LIKE clause in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeCastFailed: {
		Code:           "CastFailed",
		Description:    "Attempt to convert from one data type to another using CAST failed in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidCast: {
		Code:           "InvalidCast",
		Description:    "Attempt to convert from one data type to another using CAST failed in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeEvaluatorInvalidTimestampFormatPattern: {
		Code:           "EvaluatorInvalidTimestampFormatPattern",
		Description:    "Time stamp format pattern requires additional fields in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeEvaluatorInvalidTimestampFormatPatternSymbolForParsing: {
		Code:           "EvaluatorInvalidTimestampFormatPatternSymbolForParsing",
		Description:    "Time stamp format pattern contains a valid format symbol that cannot be applied to time stamp parsing in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeEvaluatorTimestampFormatPatternDuplicateFields: {
		Code:           "EvaluatorTimestampFormatPatternDuplicateFields",
		Description:    "Time stamp format pattern contains multiple format specifiers representing the time stamp field in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeEvaluatorTimestampFormatPatternHourClockAmPmMismatch: {
		Code:           "EvaluatorUnterminatedTimestampFormatPatternToken",
		Description:    "Time stamp format pattern contains unterminated token in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeEvaluatorUnterminatedTimestampFormatPatternToken: {
		Code:           "EvaluatorInvalidTimestampFormatPatternToken",
		Description:    "Time stamp format pattern contains an invalid token in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeEvaluatorInvalidTimestampFormatPatternToken: {
		Code:           "EvaluatorInvalidTimestampFormatPatternToken",
		Description:    "Time stamp format pattern contains an invalid token in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeEvaluatorInvalidTimestampFormatPatternSymbol: {
		Code:           "EvaluatorInvalidTimestampFormatPatternSymbol",
		Description:    "Time stamp format pattern contains an invalid symbol in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeEvaluatorBindingDoesNotExist: {
		Code:           "ErrCodeEvaluatorBindingDoesNotExist",
		Description:    "A column name or a path provided does not exist in the SQL expression",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeMissingHeaders: {
		Code:           "MissingHeaders",
		Description:    "Some headers in the query are missing from the file. Check the file and try again.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodeInvalidColumnIndex: {
		Code:           "InvalidColumnIndex",
		Description:    "The column index is invalid. Please check the service documentation and try again.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCodePostPolicyConditionInvalidFormat: {
		Code:           "PostPolicyInvalidKeyName",
		Description:    "Invalid according to Policy: Policy Conditions failed",
		HTTPStatusCode: http.StatusForbidden,
	},
	// Add your error structure here.
	ErrCodeMalformedJSON: {
		Code:           "MalformedJSON",
		Description:    "The JSON was not well-formed or did not validate against our published format.",
		HTTPStatusCode: http.StatusBadRequest,
	},
}

// GetAPIError provides API Error for input API error code.
func GetAPIError(code ErrorCode) APIError {
	return errorCodeResponse[code]
}

//
//// STSErrorCode type of error status.
//type STSErrorCode int
//
//// STSError structure
//type STSError struct {
//	Code           string
//	Description    string
//	HTTPStatusCode int
//}
//
//// Error codes,list - http://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRoleWithSAML.html
//const (
//	ErrSTSNone STSErrorCode = iota
//	ErrSTSAccessDenied
//	ErrSTSMissingParameter
//	ErrSTSInvalidParameterValue
//	ErrSTSInternalError
//)
//
//type stsErrorCodeMap map[STSErrorCode]STSError
//
////ToSTSErr code to err
//func (e stsErrorCodeMap) ToSTSErr(errCode STSErrorCode) STSError {
//	apiErr, ok := e[errCode]
//	if !ok {
//		return e[ErrSTSInternalError]
//	}
//	return apiErr
//}
//
//// StsErrCodes error code to STSError structure, these fields carry respective
//// descriptions for all the error responses.
//var StsErrCodes = stsErrorCodeMap{
//	ErrSTSAccessDenied: {
//		Code:           "AccessDenied",
//		Description:    "Generating temporary credentials not allowed for this request.",
//		HTTPStatusCode: http.StatusForbidden,
//	},
//	ErrSTSMissingParameter: {
//		Code:           "MissingParameter",
//		Description:    "A required parameter for the specified action is not supplied.",
//		HTTPStatusCode: http.StatusBadRequest,
//	},
//	ErrSTSInvalidParameterValue: {
//		Code:           "InvalidParameterValue",
//		Description:    "An invalid or out-of-range value was supplied for the input parameter.",
//		HTTPStatusCode: http.StatusBadRequest,
//	},
//	ErrSTSInternalError: {
//		Code:           "InternalError",
//		Description:    "We encountered an internal error generating credentials, please try again.",
//		HTTPStatusCode: http.StatusInternalServerError,
//	},
//}
