package apierrors

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
	ErrNone ErrorCode = iota
	ErrAccessDenied
	ErrBadDigest
	ErrEntityTooSmall
	ErrEntityTooLarge
	ErrIncompleteBody
	ErrInternalError
	ErrInvalidAccessKeyID
	ErrAccessKeyDisabled
	ErrInvalidBucketName
	ErrInvalidDigest
	ErrInvalidRange
	ErrInvalidRangePartNumber
	ErrInvalidCopyPartRange
	ErrInvalidCopyPartRangeSource
	ErrInvalidMaxKeys
	ErrInvalidEncodingMethod
	ErrInvalidMaxUploads
	ErrInvalidMaxParts
	ErrInvalidPartNumberMarker
	ErrInvalidRequestBody
	ErrInvalidCopySource
	ErrInvalidMetadataDirective
	ErrInvalidCopyDest
	ErrInvalidPolicyDocument
	ErrInvalidObjectState
	ErrMalformedXML
	ErrMissingContentLength
	ErrMissingContentMD5
	ErrMissingRequestBodyError
	ErrMissingSecurityHeader
	ErrNoSuchUser
	ErrUserAlreadyExists
	ErrNoSuchUserPolicy
	ErrUserPolicyAlreadyExists
	ErrNoSuchBucket
	ErrNoSuchBucketPolicy
	ErrNoSuchLifecycleConfiguration
	ErrNoSuchCORSConfiguration
	ErrNoSuchWebsiteConfiguration
	ErrReplicationConfigurationNotFoundError
	ErrReplicationNeedsVersioningError
	ErrReplicationBucketNeedsVersioningError
	ErrObjectRestoreAlreadyInProgress
	ErrNoSuchKey
	ErrNoSuchUpload
	ErrInvalidVersionID
	ErrNoSuchVersion
	ErrNotImplemented
	ErrPreconditionFailed
	ErrRequestTimeTooSkewed
	ErrSignatureDoesNotMatch
	ErrMethodNotAllowed
	ErrInvalidPart
	ErrInvalidPartOrder
	ErrAuthorizationHeaderMalformed
	ErrMalformedDate
	ErrMalformedPOSTRequest
	ErrPOSTFileRequired
	ErrSignatureVersionNotSupported
	ErrBucketNotEmpty
	ErrAllAccessDisabled
	ErrMalformedPolicy
	ErrMissingFields
	ErrMissingCredTag
	ErrCredMalformed
	ErrInvalidRegion

	ErrMissingSignTag
	ErrMissingSignHeadersTag

	ErrAuthHeaderEmpty
	ErrExpiredPresignRequest
	ErrRequestNotReadyYet
	ErrUnsignedHeaders
	ErrMissingDateHeader

	ErrBucketAlreadyOwnedByYou
	ErrInvalidDuration
	ErrBucketAlreadyExists
	ErrMetadataTooLarge
	ErrUnsupportedMetadata

	ErrSlowDown
	ErrBadRequest
	ErrKeyTooLongError
	ErrInvalidBucketObjectLockConfiguration
	ErrObjectLockConfigurationNotAllowed
	ErrNoSuchObjectLockConfiguration
	ErrObjectLocked
	ErrInvalidRetentionDate
	ErrPastObjectLockRetainDate
	ErrUnknownWORMModeDirective
	ErrBucketTaggingNotFound
	ErrObjectLockInvalidHeaders
	ErrInvalidTagDirective
	// Add new error codes here.

	// SSE-S3 related API errors
	ErrInvalidEncryptionMethod
	ErrInvalidQueryParams
	ErrNoAccessKey
	ErrInvalidToken

	// Bucket notification related errors.
	ErrEventNotification
	ErrARNNotification
	ErrRegionNotification
	ErrOverlappingFilterNotification
	ErrFilterNameInvalid
	ErrFilterNamePrefix
	ErrFilterNameSuffix
	ErrFilterValueInvalid
	ErrOverlappingConfigs

	// S3 extended errors.
	ErrContentSHA256Mismatch

	// Add new extended error codes here.
	ErrInvalidObjectName
	ErrInvalidObjectNamePrefixSlash
	ErrClientDisconnected
	ErrOperationTimedOut
	ErrOperationMaxedOut
	ErrInvalidRequest
	ErrIncorrectContinuationToken
	ErrInvalidFormatAccessKey

	// S3 Select Errors
	ErrEmptyRequestBody
	ErrUnsupportedFunction
	ErrInvalidExpressionType
	ErrBusy
	ErrUnauthorizedAccess
	ErrExpressionTooLong
	ErrIllegalSQLFunctionArgument
	ErrInvalidKeyPath
	ErrInvalidCompressionFormat
	ErrInvalidFileHeaderInfo
	ErrInvalidJSONType
	ErrInvalidQuoteFields
	ErrInvalidRequestParameter
	ErrInvalidDataType
	ErrInvalidTextEncoding
	ErrInvalidDataSource
	ErrInvalidTableAlias
	ErrMissingRequiredParameter
	ErrObjectSerializationConflict
	ErrUnsupportedSQLOperation
	ErrUnsupportedSQLStructure
	ErrUnsupportedSyntax
	ErrUnsupportedRangeHeader
	ErrLexerInvalidChar
	ErrLexerInvalidOperator
	ErrLexerInvalidLiteral
	ErrLexerInvalidIONLiteral
	ErrParseExpectedDatePart
	ErrParseExpectedKeyword
	ErrParseExpectedTokenType
	ErrParseExpected2TokenTypes
	ErrParseExpectedNumber
	ErrParseExpectedRightParenBuiltinFunctionCall
	ErrParseExpectedTypeName
	ErrParseExpectedWhenClause
	ErrParseUnsupportedToken
	ErrParseUnsupportedLiteralsGroupBy
	ErrParseExpectedMember
	ErrParseUnsupportedSelect
	ErrParseUnsupportedCase
	ErrParseUnsupportedCaseClause
	ErrParseUnsupportedAlias
	ErrParseUnsupportedSyntax
	ErrParseUnknownOperator
	ErrParseMissingIdentAfterAt
	ErrParseUnexpectedOperator
	ErrParseUnexpectedTerm
	ErrParseUnexpectedToken
	ErrParseUnexpectedKeyword
	ErrParseExpectedExpression
	ErrParseExpectedLeftParenAfterCast
	ErrParseExpectedLeftParenValueConstructor
	ErrParseExpectedLeftParenBuiltinFunctionCall
	ErrParseExpectedArgumentDelimiter
	ErrParseCastArity
	ErrParseInvalidTypeParam
	ErrParseEmptySelect
	ErrParseSelectMissingFrom
	ErrParseExpectedIdentForGroupName
	ErrParseExpectedIdentForAlias
	ErrParseUnsupportedCallWithStar
	ErrParseNonUnaryAgregateFunctionCall
	ErrParseMalformedJoin
	ErrParseExpectedIdentForAt
	ErrParseAsteriskIsNotAloneInSelectList
	ErrParseCannotMixSqbAndWildcardInSelectList
	ErrParseInvalidContextForWildcardInSelectList
	ErrIncorrectSQLFunctionArgumentType
	ErrValueParseFailure
	ErrEvaluatorInvalidArguments
	ErrIntegerOverflow
	ErrLikeInvalidInputs
	ErrCastFailed
	ErrInvalidCast
	ErrEvaluatorInvalidTimestampFormatPattern
	ErrEvaluatorInvalidTimestampFormatPatternSymbolForParsing
	ErrEvaluatorTimestampFormatPatternDuplicateFields
	ErrEvaluatorTimestampFormatPatternHourClockAmPmMismatch
	ErrEvaluatorUnterminatedTimestampFormatPatternToken
	ErrEvaluatorInvalidTimestampFormatPatternToken
	ErrEvaluatorInvalidTimestampFormatPatternSymbol
	ErrEvaluatorBindingDoesNotExist
	ErrMissingHeaders
	ErrInvalidColumnIndex
	ErrPostPolicyConditionInvalidFormat

	ErrMalformedJSON
)

// error code to APIError structure, these fields carry respective
// descriptions for all the error responses.
var errorCodeResponse = map[ErrorCode]APIError{
	ErrInvalidCopyDest: {
		Code:           "InvalidRequest",
		Description:    "This copy request is illegal because it is trying to copy an object to itself without changing the object's metadata, storage class, website redirect location or encryption attributes.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidCopySource: {
		Code:           "InvalidArgument",
		Description:    "Copy Source must mention the source bucket and key: sourcebucket/sourcekey.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidMetadataDirective: {
		Code:           "InvalidArgument",
		Description:    "Unknown metadata directive.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidRequestBody: {
		Code:           "InvalidArgument",
		Description:    "Body shouldn't be set for this request.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidMaxUploads: {
		Code:           "InvalidArgument",
		Description:    "Argument max-uploads must be an integer between 0 and 2147483647",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidMaxKeys: {
		Code:           "InvalidArgument",
		Description:    "Argument maxKeys must be an integer between 0 and 2147483647",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidEncodingMethod: {
		Code:           "InvalidArgument",
		Description:    "Invalid Encoding Method specified in Request",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidMaxParts: {
		Code:           "InvalidArgument",
		Description:    "Part number must be an integer between 1 and 10000, inclusive",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidPartNumberMarker: {
		Code:           "InvalidArgument",
		Description:    "Argument partNumberMarker must be an integer.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidPolicyDocument: {
		Code:           "InvalidPolicyDocument",
		Description:    "The content of the form does not meet the conditions specified in the policy document.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrAccessDenied: {
		Code:           "AccessDenied",
		Description:    "Access Denied.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrBadDigest: {
		Code:           "BadDigest",
		Description:    "The Content-Md5 you specified did not match what we received.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrEntityTooSmall: {
		Code:           "EntityTooSmall",
		Description:    "Your proposed upload is smaller than the minimum allowed object size.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrEntityTooLarge: {
		Code:           "EntityTooLarge",
		Description:    "Your proposed upload exceeds the maximum allowed object size.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrIncompleteBody: {
		Code:           "IncompleteBody",
		Description:    "You did not provide the number of bytes specified by the Content-Length HTTP header.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInternalError: {
		Code:           "InternalError",
		Description:    "We encountered an internal error, please try again.",
		HTTPStatusCode: http.StatusInternalServerError,
	},
	ErrInvalidAccessKeyID: {
		Code:           "InvalidAccessKeyId",
		Description:    "The Access Key Id you provided does not exist in our records.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrAccessKeyDisabled: {
		Code:           "InvalidAccessKeyId",
		Description:    "Your account is disabled; please contact your administrator.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrInvalidBucketName: {
		Code:           "InvalidBucketName",
		Description:    "The specified bucket is not valid.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidDigest: {
		Code:           "InvalidDigest",
		Description:    "The Content-Md5 you specified is not valid.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidRange: {
		Code:           "InvalidRange",
		Description:    "The requested range is not satisfiable",
		HTTPStatusCode: http.StatusRequestedRangeNotSatisfiable,
	},
	ErrInvalidRangePartNumber: {
		Code:           "InvalidRequest",
		Description:    "Cannot specify both Range header and partNumber query parameter",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMalformedXML: {
		Code:           "MalformedXML",
		Description:    "The XML you provided was not well-formed or did not validate against our published schema.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingContentLength: {
		Code:           "MissingContentLength",
		Description:    "You must provide the Content-Length HTTP header.",
		HTTPStatusCode: http.StatusLengthRequired,
	},
	ErrMissingContentMD5: {
		Code:           "MissingContentMD5",
		Description:    "Missing required header for this request: Content-Md5.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingSecurityHeader: {
		Code:           "MissingSecurityHeader",
		Description:    "Your request was missing a required header",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingRequestBodyError: {
		Code:           "MissingRequestBodyError",
		Description:    "Request body is empty.",
		HTTPStatusCode: http.StatusLengthRequired,
	},
	ErrNoSuchBucket: {
		Code:           "NoSuchBucket",
		Description:    "The specified bucket does not exist",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrNoSuchBucketPolicy: {
		Code:           "NoSuchBucketPolicy",
		Description:    "The bucket policy does not exist",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrNoSuchLifecycleConfiguration: {
		Code:           "NoSuchLifecycleConfiguration",
		Description:    "The lifecycle configuration does not exist",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrNoSuchUser: {
		Code:           "NoSuchUser",
		Description:    "The specified user does not exist",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrUserAlreadyExists: {
		Code:           "UserAlreadyExists",
		Description:    "The request was rejected because it attempted to create a resource that already exists .",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrNoSuchUserPolicy: {
		Code:           "NoSuchUserPolicy",
		Description:    "The specified user policy does not exist",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrUserPolicyAlreadyExists: {
		Code:           "UserPolicyAlreadyExists",
		Description:    "The same user policy already exists .",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrNoSuchKey: {
		Code:           "NoSuchKey",
		Description:    "The specified key does not exist.",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrNoSuchUpload: {
		Code:           "NoSuchUpload",
		Description:    "The specified multipart upload does not exist. The upload ID may be invalid, or the upload may have been aborted or completed.",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrInvalidVersionID: {
		Code:           "InvalidArgument",
		Description:    "Invalid version id specified",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrNoSuchVersion: {
		Code:           "NoSuchVersion",
		Description:    "The specified version does not exist.",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrNotImplemented: {
		Code:           "NotImplemented",
		Description:    "A header you provided implies functionality that is not implemented",
		HTTPStatusCode: http.StatusNotImplemented,
	},
	ErrPreconditionFailed: {
		Code:           "PreconditionFailed",
		Description:    "At least one of the pre-conditions you specified did not hold",
		HTTPStatusCode: http.StatusPreconditionFailed,
	},
	ErrRequestTimeTooSkewed: {
		Code:           "RequestTimeTooSkewed",
		Description:    "The difference between the request time and the server's time is too large.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrSignatureDoesNotMatch: {
		Code:           "SignatureDoesNotMatch",
		Description:    "The request signature we calculated does not match the signature you provided. Check your key and signing method.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrMethodNotAllowed: {
		Code:           "MethodNotAllowed",
		Description:    "The specified method is not allowed against this resource.",
		HTTPStatusCode: http.StatusMethodNotAllowed,
	},
	ErrInvalidPart: {
		Code:           "InvalidPart",
		Description:    "One or more of the specified parts could not be found.  The part may not have been uploaded, or the specified entity tag may not match the part's entity tag.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidPartOrder: {
		Code:           "InvalidPartOrder",
		Description:    "The list of parts was not in ascending order. The parts list must be specified in order by part number.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidObjectState: {
		Code:           "InvalidObjectState",
		Description:    "The operation is not valid for the current state of the object.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrAuthorizationHeaderMalformed: {
		Code:           "AuthorizationHeaderMalformed",
		Description:    "The authorization header is malformed; the region is wrong; expecting 'us-east-1'.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMalformedPOSTRequest: {
		Code:           "MalformedPOSTRequest",
		Description:    "The body of your POST request is not well-formed multipart/form-data.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrPOSTFileRequired: {
		Code:           "InvalidArgument",
		Description:    "POST requires exactly one file upload per request.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrSignatureVersionNotSupported: {
		Code:           "InvalidRequest",
		Description:    "The authorization mechanism you have provided is not supported. Please use AWS4-HMAC-SHA256.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrBucketNotEmpty: {
		Code:           "BucketNotEmpty",
		Description:    "The bucket you tried to delete is not empty",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrBucketAlreadyExists: {
		Code:           "BucketAlreadyExists",
		Description:    "The requested bucket name is not available. The bucket namespace is shared by all users of the system. Please select a different name and try again.",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrAllAccessDisabled: {
		Code:           "AllAccessDisabled",
		Description:    "All access to this resource has been disabled.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrMalformedPolicy: {
		Code:           "MalformedPolicy",
		Description:    "Policy has invalid resource.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingCredTag: {
		Code:           "InvalidRequest",
		Description:    "Missing Credential field for this request.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidRegion: {
		Code:           "InvalidRegion",
		Description:    "Region does not match.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingSignTag: {
		Code:           "AccessDenied",
		Description:    "Signature header missing Signature field.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingSignHeadersTag: {
		Code:           "InvalidArgument",
		Description:    "Signature header missing SignedHeaders field.",
		HTTPStatusCode: http.StatusBadRequest,
	},

	ErrAuthHeaderEmpty: {
		Code:           "InvalidArgument",
		Description:    "Authorization header is invalid -- one and only one ' ' (space) required.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingDateHeader: {
		Code:           "AccessDenied",
		Description:    "AWS authentication requires a valid Date or x-amz-date header",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrExpiredPresignRequest: {
		Code:           "AccessDenied",
		Description:    "Request has expired",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrRequestNotReadyYet: {
		Code:           "AccessDenied",
		Description:    "Request is not valid yet",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrSlowDown: {
		Code:           "SlowDown",
		Description:    "Resource requested is unreadable, please reduce your request rate",
		HTTPStatusCode: http.StatusServiceUnavailable,
	},
	ErrBadRequest: {
		Code:           "BadRequest",
		Description:    "400 BadRequest",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrKeyTooLongError: {
		Code:           "KeyTooLongError",
		Description:    "Your key is too long",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrUnsignedHeaders: {
		Code:           "AccessDenied",
		Description:    "There were headers present in the request which were not signed",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrBucketAlreadyOwnedByYou: {
		Code:           "BucketAlreadyOwnedByYou",
		Description:    "Your previous request to create the named bucket succeeded and you already own it.",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrInvalidDuration: {
		Code:           "InvalidDuration",
		Description:    "Duration provided in the request is invalid.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidBucketObjectLockConfiguration: {
		Code:           "InvalidRequest",
		Description:    "Bucket is missing ObjectLockConfiguration",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrBucketTaggingNotFound: {
		Code:           "NoSuchTagSet",
		Description:    "The TagSet does not exist",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrObjectLockConfigurationNotAllowed: {
		Code:           "InvalidBucketState",
		Description:    "Object Lock configuration cannot be enabled on existing buckets",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrNoSuchCORSConfiguration: {
		Code:           "NoSuchCORSConfiguration",
		Description:    "The CORS configuration does not exist",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrNoSuchWebsiteConfiguration: {
		Code:           "NoSuchWebsiteConfiguration",
		Description:    "The specified bucket does not have a website configuration",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrReplicationConfigurationNotFoundError: {
		Code:           "ReplicationConfigurationNotFoundError",
		Description:    "The replication configuration was not found",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrReplicationNeedsVersioningError: {
		Code:           "InvalidRequest",
		Description:    "Versioning must be 'Enabled' on the bucket to apply a replication configuration",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrReplicationBucketNeedsVersioningError: {
		Code:           "InvalidRequest",
		Description:    "Versioning must be 'Enabled' on the bucket to add a replication target",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrNoSuchObjectLockConfiguration: {
		Code:           "NoSuchObjectLockConfiguration",
		Description:    "The specified object does not have a ObjectLock configuration",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrObjectLocked: {
		Code:           "InvalidRequest",
		Description:    "Object is WORM protected and cannot be overwritten",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidRetentionDate: {
		Code:           "InvalidRequest",
		Description:    "Date must be provided in ISO 8601 format",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrPastObjectLockRetainDate: {
		Code:           "InvalidRequest",
		Description:    "the retain until date must be in the future",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrUnknownWORMModeDirective: {
		Code:           "InvalidRequest",
		Description:    "unknown wormMode directive",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrObjectLockInvalidHeaders: {
		Code:           "InvalidRequest",
		Description:    "x-amz-object-lock-retain-until-date and x-amz-object-lock-mode must both be supplied",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrObjectRestoreAlreadyInProgress: {
		Code:           "RestoreAlreadyInProgress",
		Description:    "Object restore is already in progress",
		HTTPStatusCode: http.StatusConflict,
	},
	// Bucket notification related errors.
	ErrEventNotification: {
		Code:           "InvalidArgument",
		Description:    "A specified event is not supported for notifications.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrARNNotification: {
		Code:           "InvalidArgument",
		Description:    "A specified destination ARN does not exist or is not well-formed. Verify the destination ARN.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrRegionNotification: {
		Code:           "InvalidArgument",
		Description:    "A specified destination is in a different region than the bucket. You must use a destination that resides in the same region as the bucket.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrOverlappingFilterNotification: {
		Code:           "InvalidArgument",
		Description:    "An object key name filtering rule defined with overlapping prefixes, overlapping suffixes, or overlapping combinations of prefixes and suffixes for the same event types.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrFilterNameInvalid: {
		Code:           "InvalidArgument",
		Description:    "filter rule name must be either prefix or suffix",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrFilterNamePrefix: {
		Code:           "InvalidArgument",
		Description:    "Cannot specify more than one prefix rule in a filter.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrFilterNameSuffix: {
		Code:           "InvalidArgument",
		Description:    "Cannot specify more than one suffix rule in a filter.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrFilterValueInvalid: {
		Code:           "InvalidArgument",
		Description:    "Size of filter rule value cannot exceed 1024 bytes in UTF-8 representation",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrOverlappingConfigs: {
		Code:           "InvalidArgument",
		Description:    "Configurations overlap. Configurations on the same bucket cannot share a common event type.",
		HTTPStatusCode: http.StatusBadRequest,
	},

	ErrInvalidCopyPartRange: {
		Code:           "InvalidArgument",
		Description:    "The x-amz-copy-source-range value must be of the form bytes=first-last where first and last are the zero-based offsets of the first and last bytes to copy",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidCopyPartRangeSource: {
		Code:           "InvalidArgument",
		Description:    "Range specified is not valid for source object",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMetadataTooLarge: {
		Code:           "MetadataTooLarge",
		Description:    "Your metadata headers exceed the maximum allowed metadata size.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidTagDirective: {
		Code:           "InvalidArgument",
		Description:    "Unknown tag directive.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidEncryptionMethod: {
		Code:           "InvalidRequest",
		Description:    "The encryption method specified is not supported",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidQueryParams: {
		Code:           "AuthorizationQueryParametersError",
		Description:    "Query-string authentication version 4 requires the X-Amz-Algorithm, X-Amz-Credential, X-Amz-Signature, X-Amz-Date, X-Amz-SignedHeaders, and X-Amz-Expires parameters.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrNoAccessKey: {
		Code:           "AccessDenied",
		Description:    "No AWSAccessKey was presented",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrInvalidToken: {
		Code:           "InvalidTokenId",
		Description:    "The security token included in the request is invalid",
		HTTPStatusCode: http.StatusForbidden,
	},

	// S3 extensions.
	ErrInvalidObjectName: {
		Code:           "InvalidObjectName",
		Description:    "Object name contains unsupported characters.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidObjectNamePrefixSlash: {
		Code:           "InvalidObjectName",
		Description:    "Object name contains a leading slash.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrClientDisconnected: {
		Code:           "ClientDisconnected",
		Description:    "Client disconnected before response was ready",
		HTTPStatusCode: 499, // No official code, use nginx value.
	},
	ErrOperationTimedOut: {
		Code:           "RequestTimeout",
		Description:    "A timeout occurred while trying to lock a resource, please reduce your request rate",
		HTTPStatusCode: http.StatusServiceUnavailable,
	},
	ErrOperationMaxedOut: {
		Code:           "SlowDown",
		Description:    "A timeout exceeded while waiting to proceed with the request, please reduce your request rate",
		HTTPStatusCode: http.StatusServiceUnavailable,
	},
	ErrUnsupportedMetadata: {
		Code:           "InvalidArgument",
		Description:    "Your metadata headers are not supported.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	// Generic Invalid-Request error. Should be used for response errors only for unlikely
	// corner case errors for which introducing new APIErrorCode is not worth it. LogIf()
	// should be used to log the error at the source of the error for debugging purposes.
	ErrInvalidRequest: {
		Code:           "InvalidRequest",
		Description:    "Invalid Request",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrIncorrectContinuationToken: {
		Code:           "InvalidArgument",
		Description:    "The continuation token provided is incorrect",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidFormatAccessKey: {
		Code:           "InvalidAccessKeyId",
		Description:    "The Access Key Id you provided contains invalid characters.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	// S3 Select API Errors
	ErrEmptyRequestBody: {
		Code:           "EmptyRequestBody",
		Description:    "Request body cannot be empty.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrUnsupportedFunction: {
		Code:           "UnsupportedFunction",
		Description:    "Encountered an unsupported SQL function.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidDataSource: {
		Code:           "InvalidDataSource",
		Description:    "Invalid data source type. Only CSV and JSON are supported at this time.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidExpressionType: {
		Code:           "InvalidExpressionType",
		Description:    "The ExpressionType is invalid. Only SQL expressions are supported at this time.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrBusy: {
		Code:           "Busy",
		Description:    "The service is unavailable. Please retry.",
		HTTPStatusCode: http.StatusServiceUnavailable,
	},
	ErrUnauthorizedAccess: {
		Code:           "UnauthorizedAccess",
		Description:    "You are not authorized to perform this operation",
		HTTPStatusCode: http.StatusUnauthorized,
	},
	ErrExpressionTooLong: {
		Code:           "ExpressionTooLong",
		Description:    "The SQL expression is too long: The maximum byte-length for the SQL expression is 256 KB.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrIllegalSQLFunctionArgument: {
		Code:           "IllegalSqlFunctionArgument",
		Description:    "Illegal argument was used in the SQL function.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidKeyPath: {
		Code:           "InvalidKeyPath",
		Description:    "Key path in the SQL expression is invalid.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidCompressionFormat: {
		Code:           "InvalidCompressionFormat",
		Description:    "The file is not in a supported compression format. Only GZIP is supported at this time.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidFileHeaderInfo: {
		Code:           "InvalidFileHeaderInfo",
		Description:    "The FileHeaderInfo is invalid. Only NONE, USE, and IGNORE are supported.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidJSONType: {
		Code:           "InvalidJsonType",
		Description:    "The JsonType is invalid. Only DOCUMENT and LINES are supported at this time.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidQuoteFields: {
		Code:           "InvalidQuoteFields",
		Description:    "The QuoteFields is invalid. Only ALWAYS and ASNEEDED are supported.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidRequestParameter: {
		Code:           "InvalidRequestParameter",
		Description:    "The value of a parameter in SelectRequest element is invalid. Check the service API documentation and try again.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidDataType: {
		Code:           "InvalidDataType",
		Description:    "The SQL expression contains an invalid data type.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidTextEncoding: {
		Code:           "InvalidTextEncoding",
		Description:    "Invalid encoding type. Only UTF-8 encoding is supported at this time.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidTableAlias: {
		Code:           "InvalidTableAlias",
		Description:    "The SQL expression contains an invalid table alias.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingRequiredParameter: {
		Code:           "MissingRequiredParameter",
		Description:    "The SelectRequest entity is missing a required parameter. Check the service documentation and try again.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrObjectSerializationConflict: {
		Code:           "ObjectSerializationConflict",
		Description:    "The SelectRequest entity can only contain one of CSV or JSON. Check the service documentation and try again.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrUnsupportedSQLOperation: {
		Code:           "UnsupportedSqlOperation",
		Description:    "Encountered an unsupported SQL operation.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrUnsupportedSQLStructure: {
		Code:           "UnsupportedSqlStructure",
		Description:    "Encountered an unsupported SQL structure. Check the SQL Reference.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrUnsupportedSyntax: {
		Code:           "UnsupportedSyntax",
		Description:    "Encountered invalid syntax.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrUnsupportedRangeHeader: {
		Code:           "UnsupportedRangeHeader",
		Description:    "Range header is not supported for this operation.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrLexerInvalidChar: {
		Code:           "LexerInvalidChar",
		Description:    "The SQL expression contains an invalid character.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrLexerInvalidOperator: {
		Code:           "LexerInvalidOperator",
		Description:    "The SQL expression contains an invalid literal.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrLexerInvalidLiteral: {
		Code:           "LexerInvalidLiteral",
		Description:    "The SQL expression contains an invalid operator.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrLexerInvalidIONLiteral: {
		Code:           "LexerInvalidIONLiteral",
		Description:    "The SQL expression contains an invalid operator.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedDatePart: {
		Code:           "ParseExpectedDatePart",
		Description:    "Did not find the expected date part in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedKeyword: {
		Code:           "ParseExpectedKeyword",
		Description:    "Did not find the expected keyword in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedTokenType: {
		Code:           "ParseExpectedTokenType",
		Description:    "Did not find the expected token in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpected2TokenTypes: {
		Code:           "ParseExpected2TokenTypes",
		Description:    "Did not find the expected token in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedNumber: {
		Code:           "ParseExpectedNumber",
		Description:    "Did not find the expected number in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedRightParenBuiltinFunctionCall: {
		Code:           "ParseExpectedRightParenBuiltinFunctionCall",
		Description:    "Did not find the expected right parenthesis character in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedTypeName: {
		Code:           "ParseExpectedTypeName",
		Description:    "Did not find the expected type name in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedWhenClause: {
		Code:           "ParseExpectedWhenClause",
		Description:    "Did not find the expected WHEN clause in the SQL expression. CASE is not supported.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseUnsupportedToken: {
		Code:           "ParseUnsupportedToken",
		Description:    "The SQL expression contains an unsupported token.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseUnsupportedLiteralsGroupBy: {
		Code:           "ParseUnsupportedLiteralsGroupBy",
		Description:    "The SQL expression contains an unsupported use of GROUP BY.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedMember: {
		Code:           "ParseExpectedMember",
		Description:    "The SQL expression contains an unsupported use of MEMBER.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseUnsupportedSelect: {
		Code:           "ParseUnsupportedSelect",
		Description:    "The SQL expression contains an unsupported use of SELECT.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseUnsupportedCase: {
		Code:           "ParseUnsupportedCase",
		Description:    "The SQL expression contains an unsupported use of CASE.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseUnsupportedCaseClause: {
		Code:           "ParseUnsupportedCaseClause",
		Description:    "The SQL expression contains an unsupported use of CASE.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseUnsupportedAlias: {
		Code:           "ParseUnsupportedAlias",
		Description:    "The SQL expression contains an unsupported use of ALIAS.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseUnsupportedSyntax: {
		Code:           "ParseUnsupportedSyntax",
		Description:    "The SQL expression contains unsupported syntax.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseUnknownOperator: {
		Code:           "ParseUnknownOperator",
		Description:    "The SQL expression contains an invalid operator.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseMissingIdentAfterAt: {
		Code:           "ParseMissingIdentAfterAt",
		Description:    "Did not find the expected identifier after the @ symbol in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseUnexpectedOperator: {
		Code:           "ParseUnexpectedOperator",
		Description:    "The SQL expression contains an unexpected operator.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseUnexpectedTerm: {
		Code:           "ParseUnexpectedTerm",
		Description:    "The SQL expression contains an unexpected term.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseUnexpectedToken: {
		Code:           "ParseUnexpectedToken",
		Description:    "The SQL expression contains an unexpected token.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseUnexpectedKeyword: {
		Code:           "ParseUnexpectedKeyword",
		Description:    "The SQL expression contains an unexpected keyword.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedExpression: {
		Code:           "ParseExpectedExpression",
		Description:    "Did not find the expected SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedLeftParenAfterCast: {
		Code:           "ParseExpectedLeftParenAfterCast",
		Description:    "Did not find expected the left parenthesis in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedLeftParenValueConstructor: {
		Code:           "ParseExpectedLeftParenValueConstructor",
		Description:    "Did not find expected the left parenthesis in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedLeftParenBuiltinFunctionCall: {
		Code:           "ParseExpectedLeftParenBuiltinFunctionCall",
		Description:    "Did not find the expected left parenthesis in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedArgumentDelimiter: {
		Code:           "ParseExpectedArgumentDelimiter",
		Description:    "Did not find the expected argument delimiter in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseCastArity: {
		Code:           "ParseCastArity",
		Description:    "The SQL expression CAST has incorrect arity.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseInvalidTypeParam: {
		Code:           "ParseInvalidTypeParam",
		Description:    "The SQL expression contains an invalid parameter value.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseEmptySelect: {
		Code:           "ParseEmptySelect",
		Description:    "The SQL expression contains an empty SELECT.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseSelectMissingFrom: {
		Code:           "ParseSelectMissingFrom",
		Description:    "GROUP is not supported in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedIdentForGroupName: {
		Code:           "ParseExpectedIdentForGroupName",
		Description:    "GROUP is not supported in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedIdentForAlias: {
		Code:           "ParseExpectedIdentForAlias",
		Description:    "Did not find the expected identifier for the alias in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseUnsupportedCallWithStar: {
		Code:           "ParseUnsupportedCallWithStar",
		Description:    "Only COUNT with (*) as a parameter is supported in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseNonUnaryAgregateFunctionCall: {
		Code:           "ParseNonUnaryAgregateFunctionCall",
		Description:    "Only one argument is supported for aggregate functions in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseMalformedJoin: {
		Code:           "ParseMalformedJoin",
		Description:    "JOIN is not supported in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseExpectedIdentForAt: {
		Code:           "ParseExpectedIdentForAt",
		Description:    "Did not find the expected identifier for AT name in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseAsteriskIsNotAloneInSelectList: {
		Code:           "ParseAsteriskIsNotAloneInSelectList",
		Description:    "Other expressions are not allowed in the SELECT list when '*' is used without dot notation in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseCannotMixSqbAndWildcardInSelectList: {
		Code:           "ParseCannotMixSqbAndWildcardInSelectList",
		Description:    "Cannot mix [] and * in the same expression in a SELECT list in SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrParseInvalidContextForWildcardInSelectList: {
		Code:           "ParseInvalidContextForWildcardInSelectList",
		Description:    "Invalid use of * in SELECT list in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrIncorrectSQLFunctionArgumentType: {
		Code:           "IncorrectSqlFunctionArgumentType",
		Description:    "Incorrect type of arguments in function call in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrValueParseFailure: {
		Code:           "ValueParseFailure",
		Description:    "Time stamp parse failure in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrEvaluatorInvalidArguments: {
		Code:           "EvaluatorInvalidArguments",
		Description:    "Incorrect number of arguments in the function call in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrIntegerOverflow: {
		Code:           "IntegerOverflow",
		Description:    "Int overflow or underflow in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrLikeInvalidInputs: {
		Code:           "LikeInvalidInputs",
		Description:    "Invalid argument given to the LIKE clause in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCastFailed: {
		Code:           "CastFailed",
		Description:    "Attempt to convert from one data type to another using CAST failed in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidCast: {
		Code:           "InvalidCast",
		Description:    "Attempt to convert from one data type to another using CAST failed in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrEvaluatorInvalidTimestampFormatPattern: {
		Code:           "EvaluatorInvalidTimestampFormatPattern",
		Description:    "Time stamp format pattern requires additional fields in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrEvaluatorInvalidTimestampFormatPatternSymbolForParsing: {
		Code:           "EvaluatorInvalidTimestampFormatPatternSymbolForParsing",
		Description:    "Time stamp format pattern contains a valid format symbol that cannot be applied to time stamp parsing in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrEvaluatorTimestampFormatPatternDuplicateFields: {
		Code:           "EvaluatorTimestampFormatPatternDuplicateFields",
		Description:    "Time stamp format pattern contains multiple format specifiers representing the time stamp field in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrEvaluatorTimestampFormatPatternHourClockAmPmMismatch: {
		Code:           "EvaluatorUnterminatedTimestampFormatPatternToken",
		Description:    "Time stamp format pattern contains unterminated token in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrEvaluatorUnterminatedTimestampFormatPatternToken: {
		Code:           "EvaluatorInvalidTimestampFormatPatternToken",
		Description:    "Time stamp format pattern contains an invalid token in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrEvaluatorInvalidTimestampFormatPatternToken: {
		Code:           "EvaluatorInvalidTimestampFormatPatternToken",
		Description:    "Time stamp format pattern contains an invalid token in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrEvaluatorInvalidTimestampFormatPatternSymbol: {
		Code:           "EvaluatorInvalidTimestampFormatPatternSymbol",
		Description:    "Time stamp format pattern contains an invalid symbol in the SQL expression.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrEvaluatorBindingDoesNotExist: {
		Code:           "ErrEvaluatorBindingDoesNotExist",
		Description:    "A column name or a path provided does not exist in the SQL expression",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingHeaders: {
		Code:           "MissingHeaders",
		Description:    "Some headers in the query are missing from the file. Check the file and try again.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidColumnIndex: {
		Code:           "InvalidColumnIndex",
		Description:    "The column index is invalid. Please check the service documentation and try again.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrPostPolicyConditionInvalidFormat: {
		Code:           "PostPolicyInvalidKeyName",
		Description:    "Invalid according to Policy: Policy Conditions failed",
		HTTPStatusCode: http.StatusForbidden,
	},
	// Add your error structure here.
	ErrMalformedJSON: {
		Code:           "MalformedJSON",
		Description:    "The JSON was not well-formed or did not validate against our published format.",
		HTTPStatusCode: http.StatusBadRequest,
	},
}

// GetAPIError provides API Error for input API error code.
func GetAPIError(code ErrorCode) APIError {
	return errorCodeResponse[code]
}

// STSErrorCode type of error status.
type STSErrorCode int

// STSError structure
type STSError struct {
	Code           string
	Description    string
	HTTPStatusCode int
}

// Error codes,list - http://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRoleWithSAML.html
const (
	ErrSTSNone STSErrorCode = iota
	ErrSTSAccessDenied
	ErrSTSMissingParameter
	ErrSTSInvalidParameterValue
	ErrSTSInternalError
)

type stsErrorCodeMap map[STSErrorCode]STSError

//ToSTSErr code to err
func (e stsErrorCodeMap) ToSTSErr(errCode STSErrorCode) STSError {
	apiErr, ok := e[errCode]
	if !ok {
		return e[ErrSTSInternalError]
	}
	return apiErr
}

// StsErrCodes error code to STSError structure, these fields carry respective
// descriptions for all the error responses.
var StsErrCodes = stsErrorCodeMap{
	ErrSTSAccessDenied: {
		Code:           "AccessDenied",
		Description:    "Generating temporary credentials not allowed for this request.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrSTSMissingParameter: {
		Code:           "MissingParameter",
		Description:    "A required parameter for the specified action is not supplied.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrSTSInvalidParameterValue: {
		Code:           "InvalidParameterValue",
		Description:    "An invalid or out-of-range value was supplied for the input parameter.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrSTSInternalError: {
		Code:           "InternalError",
		Description:    "We encountered an internal error generating credentials, please try again.",
		HTTPStatusCode: http.StatusInternalServerError,
	},
}
