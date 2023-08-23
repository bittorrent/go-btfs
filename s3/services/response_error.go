package services

import (
	"fmt"
	"net/http"
)

type ResponseError struct {
	code           string
	description    string
	httpStatusCode int
}

func (err *ResponseError) Code() string {
	return err.code
}

func (err *ResponseError) Description() string {
	return err.description
}

func (err *ResponseError) HTTPStatusCode() int {
	return err.httpStatusCode
}

func (err *ResponseError) Error() string {
	return fmt.Sprintf("[%s]%s", err.code, err.description)
}

// Errors http://docs.aws.amazon.com/AmazonS3/latest/API/ErrorResponses.html
var (
	RespErrInvalidCopyDest = &ResponseError{
		code:           "InvalidRequest",
		description:    "This copy request is illegal because it is trying to copy an object to itself without changing the object's metadata, storage class, website redirect location or encryption attributes.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidCopySource = &ResponseError{
		code:           "InvalidArgument",
		description:    "Copy Source must mention the source bucket and key: sourcebucket/sourcekey.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidMetadataDirective = &ResponseError{
		code:           "InvalidArgument",
		description:    "Unknown metadata directive.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidRequestBody = &ResponseError{
		code:           "InvalidArgument",
		description:    "Body shouldn't be set for this request.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidMaxUploads = &ResponseError{
		code:           "InvalidArgument",
		description:    "Argument max-uploads must be an integer between 0 and 2147483647",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidMaxKeys = &ResponseError{
		code:           "InvalidArgument",
		description:    "Argument maxKeys must be an integer between 0 and 2147483647",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidEncodingMethod = &ResponseError{
		code:           "InvalidArgument",
		description:    "Invalid Encoding Method specified in Request",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidMaxParts = &ResponseError{
		code:           "InvalidArgument",
		description:    "Part number must be an integer between 1 and 10000, inclusive",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidPartNumberMarker = &ResponseError{
		code:           "InvalidArgument",
		description:    "Argument partNumberMarker must be an integer.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidPolicyDocument = &ResponseError{
		code:           "InvalidPolicyDocument",
		description:    "The content of the form does not meet the conditions specified in the policy document.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrAccessDenied = &ResponseError{
		code:           "AccessDenied",
		description:    "Access Denied.",
		httpStatusCode: http.StatusForbidden,
	}
	RespErrBadDigest = &ResponseError{
		code:           "BadDigest",
		description:    "The Content-Md5 you specified did not match what we received.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrEntityTooSmall = &ResponseError{
		code:           "EntityTooSmall",
		description:    "Your proposed upload is smaller than the minimum allowed object size.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrEntityTooLarge = &ResponseError{
		code:           "EntityTooLarge",
		description:    "Your proposed upload exceeds the maximum allowed object size.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrIncompleteBody = &ResponseError{
		code:           "IncompleteBody",
		description:    "You did not provide the number of bytes specified by the Content-Length HTTP header.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInternalError = &ResponseError{
		code:           "InternalError",
		description:    "We encountered an internal error, please try again.",
		httpStatusCode: http.StatusInternalServerError,
	}
	RespErrInvalidAccessKeyID = &ResponseError{
		code:           "InvalidAccessKeyId",
		description:    "The Access Key Id you provided does not exist in our records.",
		httpStatusCode: http.StatusForbidden,
	}
	RespErrAccessKeyDisabled = &ResponseError{
		code:           "InvalidAccessKeyId",
		description:    "Your account is disabled; please contact your administrator.",
		httpStatusCode: http.StatusForbidden,
	}
	RespErrInvalidBucketName = &ResponseError{
		code:           "InvalidBucketName",
		description:    "The specified bucket is not valid.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidDigest = &ResponseError{
		code:           "InvalidDigest",
		description:    "The Content-Md5 you specified is not valid.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidRange = &ResponseError{
		code:           "InvalidRange",
		description:    "The requested range is not satisfiable",
		httpStatusCode: http.StatusRequestedRangeNotSatisfiable,
	}
	RespErrInvalidRangePartNumber = &ResponseError{
		code:           "InvalidRequest",
		description:    "Cannot specify both Range header and partNumber query parameter",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrMalformedXML = &ResponseError{
		code:           "MalformedXML",
		description:    "The XML you provided was not well-formed or did not validate against our published schema.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrMissingContentLength = &ResponseError{
		code:           "MissingContentLength",
		description:    "You must provide the Content-Length HTTP header.",
		httpStatusCode: http.StatusLengthRequired,
	}
	RespErrMissingContentMD5 = &ResponseError{
		code:           "MissingContentMD5",
		description:    "Missing required header for this request: Content-Md5.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrMissingSecurityHeader = &ResponseError{
		code:           "MissingSecurityHeader",
		description:    "Your request was missing a required header",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrMissingRequestBodyError = &ResponseError{
		code:           "MissingRequestBodyError",
		description:    "Request body is empty.",
		httpStatusCode: http.StatusLengthRequired,
	}
	RespErrNoSuchBucket = &ResponseError{
		code:           "NoSuchBucket",
		description:    "The specified bucket does not exist",
		httpStatusCode: http.StatusNotFound,
	}
	RespErrNoSuchBucketPolicy = &ResponseError{
		code:           "NoSuchBucketPolicy",
		description:    "The bucket policy does not exist",
		httpStatusCode: http.StatusNotFound,
	}
	RespErrNoSuchLifecycleConfiguration = &ResponseError{
		code:           "NoSuchLifecycleConfiguration",
		description:    "The lifecycle configuration does not exist",
		httpStatusCode: http.StatusNotFound,
	}
	RespErrNoSuchUser = &ResponseError{
		code:           "NoSuchUser",
		description:    "The specified user does not exist",
		httpStatusCode: http.StatusConflict,
	}
	RespErrUserAlreadyExists = &ResponseError{
		code:           "UserAlreadyExists",
		description:    "The request was rejected because it attempted to create a resource that already exists .",
		httpStatusCode: http.StatusConflict,
	}
	RespErrNoSuchUserPolicy = &ResponseError{
		code:           "NoSuchUserPolicy",
		description:    "The specified user policy does not exist",
		httpStatusCode: http.StatusConflict,
	}
	RespErrUserPolicyAlreadyExists = &ResponseError{
		code:           "UserPolicyAlreadyExists",
		description:    "The same user policy already exists .",
		httpStatusCode: http.StatusConflict,
	}
	RespErrNoSuchKey = &ResponseError{
		code:           "NoSuchKey",
		description:    "The specified key does not exist.",
		httpStatusCode: http.StatusNotFound,
	}
	RespErrNoSuchUpload = &ResponseError{
		code:           "NoSuchUpload",
		description:    "The specified multipart upload does not exist. The upload ID may be invalid, or the upload may have been aborted or completed.",
		httpStatusCode: http.StatusNotFound,
	}
	RespErrInvalidVersionID = &ResponseError{
		code:           "InvalidArgument",
		description:    "Invalid version id specified",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrNoSuchVersion = &ResponseError{
		code:           "NoSuchVersion",
		description:    "The specified version does not exist.",
		httpStatusCode: http.StatusNotFound,
	}
	RespErrNotImplemented = &ResponseError{
		code:           "NotImplemented",
		description:    "A header you provided implies functionality that is not implemented",
		httpStatusCode: http.StatusNotImplemented,
	}
	RespErrPreconditionFailed = &ResponseError{
		code:           "PreconditionFailed",
		description:    "At least one of the pre-conditions you specified did not hold",
		httpStatusCode: http.StatusPreconditionFailed,
	}
	RespErrRequestTimeTooSkewed = &ResponseError{
		code:           "RequestTimeTooSkewed",
		description:    "The difference between the request time and the server's time is too large.",
		httpStatusCode: http.StatusForbidden,
	}
	RespErrSignatureDoesNotMatch = &ResponseError{
		code:           "SignatureDoesNotMatch",
		description:    "The request signature we calculated does not match the signature you provided. Check your key and signing method.",
		httpStatusCode: http.StatusForbidden,
	}
	RespErrMethodNotAllowed = &ResponseError{
		code:           "MethodNotAllowed",
		description:    "The specified method is not allowed against this resource.",
		httpStatusCode: http.StatusMethodNotAllowed,
	}
	RespErrInvalidPart = &ResponseError{
		code:           "InvalidPart",
		description:    "One or more of the specified parts could not be found.  The part may not have been uploaded, or the specified entity tag may not match the part's entity tag.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidPartOrder = &ResponseError{
		code:           "InvalidPartOrder",
		description:    "The list of parts was not in ascending order. The parts list must be specified in order by part number.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidObjectState = &ResponseError{
		code:           "InvalidObjectState",
		description:    "The operation is not valid for the current state of the object.",
		httpStatusCode: http.StatusForbidden,
	}
	RespErrAuthorizationHeaderMalformed = &ResponseError{
		code:           "AuthorizationHeaderMalformed",
		description:    "The authorization header is malformed; the region is wrong; expecting 'us-east-1'.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrMalformedPOSTRequest = &ResponseError{
		code:           "MalformedPOSTRequest",
		description:    "The body of your POST request is not well-formed multipart/form-data.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrPOSTFileRequired = &ResponseError{
		code:           "InvalidArgument",
		description:    "POST requires exactly one file upload per request.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrSignatureVersionNotSupported = &ResponseError{
		code:           "InvalidRequest",
		description:    "The authorization mechanism you have provided is not supported. Please use AWS4-HMAC-SHA256.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrBucketNotEmpty = &ResponseError{
		code:           "BucketNotEmpty",
		description:    "The bucket you tried to delete is not empty",
		httpStatusCode: http.StatusConflict,
	}
	RespErrBucketAlreadyExists = &ResponseError{
		code:           "BucketAlreadyExists",
		description:    "The requested bucket name is not available. The bucket namespace is shared by all users of the system. Please select a different name and try again.",
		httpStatusCode: http.StatusConflict,
	}
	RespErrAllAccessDisabled = &ResponseError{
		code:           "AllAccessDisabled",
		description:    "All access to this resource has been disabled.",
		httpStatusCode: http.StatusForbidden,
	}
	RespErrMalformedPolicy = &ResponseError{
		code:           "MalformedPolicy",
		description:    "Policy has invalid resource.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrMissingFields = &ResponseError{ // todo
		code:           "InvalidRequest",
		description:    "ErrMissingFields",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrMissingCredTag = &ResponseError{
		code:           "InvalidRequest",
		description:    "Missing Credential field for this request.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrCredMalformed = &ResponseError{ // todo
		code:           "InvalidRequest",
		description:    "ErrCredMalformed",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidRegion = &ResponseError{
		code:           "InvalidRegion",
		description:    "Region does not match.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrMissingSignTag = &ResponseError{
		code:           "AccessDenied",
		description:    "Signature header missing Signature field.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrMissingSignHeadersTag = &ResponseError{
		code:           "InvalidArgument",
		description:    "Signature header missing SignedHeaders field.",
		httpStatusCode: http.StatusBadRequest,
	}

	RespErrAuthHeaderEmpty = &ResponseError{
		code:           "InvalidArgument",
		description:    "Authorization header is invalid -- one and only one ' ' (space) required.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrMissingDateHeader = &ResponseError{
		code:           "AccessDenied",
		description:    "AWS authentication requires a valid Date or x-amz-date header",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrExpiredPresignRequest = &ResponseError{
		code:           "AccessDenied",
		description:    "Request has expired",
		httpStatusCode: http.StatusForbidden,
	}
	RespErrRequestNotReadyYet = &ResponseError{
		code:           "AccessDenied",
		description:    "Request is not valid yet",
		httpStatusCode: http.StatusForbidden,
	}
	RespErrSlowDown = &ResponseError{
		code:           "SlowDown",
		description:    "Resource requested is unreadable, please reduce your request rate",
		httpStatusCode: http.StatusServiceUnavailable,
	}
	RespErrBadRequest = &ResponseError{
		code:           "BadRequest",
		description:    "400 BadRequest",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrKeyTooLongError = &ResponseError{
		code:           "KeyTooLongError",
		description:    "Your key is too long",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrUnsignedHeaders = &ResponseError{
		code:           "AccessDenied",
		description:    "There were headers present in the request which were not signed",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrBucketAlreadyOwnedByYou = &ResponseError{
		code:           "BucketAlreadyOwnedByYou",
		description:    "Your previous request to create the named bucket succeeded and you already own it.",
		httpStatusCode: http.StatusConflict,
	}
	RespErrInvalidDuration = &ResponseError{
		code:           "InvalidDuration",
		description:    "Duration provided in the request is invalid.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidBucketObjectLockConfiguration = &ResponseError{
		code:           "InvalidRequest",
		description:    "Bucket is missing ObjectLockConfiguration",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrBucketTaggingNotFound = &ResponseError{
		code:           "NoSuchTagSet",
		description:    "The TagSet does not exist",
		httpStatusCode: http.StatusNotFound,
	}
	RespErrObjectLockConfigurationNotAllowed = &ResponseError{
		code:           "InvalidBucketState",
		description:    "Object Lock configuration cannot be enabled on existing buckets",
		httpStatusCode: http.StatusConflict,
	}
	RespErrNoSuchCORSConfiguration = &ResponseError{
		code:           "NoSuchCORSConfiguration",
		description:    "The CORS configuration does not exist",
		httpStatusCode: http.StatusNotFound,
	}
	RespErrNoSuchWebsiteConfiguration = &ResponseError{
		code:           "NoSuchWebsiteConfiguration",
		description:    "The specified bucket does not have a website configuration",
		httpStatusCode: http.StatusNotFound,
	}
	RespErrReplicationConfigurationNotFoundError = &ResponseError{
		code:           "ReplicationConfigurationNotFoundError",
		description:    "The replication configuration was not found",
		httpStatusCode: http.StatusNotFound,
	}
	RespErrReplicationNeedsVersioningError = &ResponseError{
		code:           "InvalidRequest",
		description:    "Versioning must be 'Enabled' on the bucket to apply a replication configuration",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrReplicationBucketNeedsVersioningError = &ResponseError{
		code:           "InvalidRequest",
		description:    "Versioning must be 'Enabled' on the bucket to add a replication target",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrNoSuchObjectLockConfiguration = &ResponseError{
		code:           "NoSuchObjectLockConfiguration",
		description:    "The specified object does not have a ObjectLock configuration",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrObjectLocked = &ResponseError{
		code:           "InvalidRequest",
		description:    "Object is WORM protected and cannot be overwritten",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidRetentionDate = &ResponseError{
		code:           "InvalidRequest",
		description:    "Date must be provided in ISO 8601 format",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrPastObjectLockRetainDate = &ResponseError{
		code:           "InvalidRequest",
		description:    "the retain until date must be in the future",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrUnknownWORMModeDirective = &ResponseError{
		code:           "InvalidRequest",
		description:    "unknown wormMode directive",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrObjectLockInvalidHeaders = &ResponseError{
		code:           "InvalidRequest",
		description:    "x-amz-object-lock-retain-until-date and x-amz-object-lock-mode must both be supplied",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrObjectRestoreAlreadyInProgress = &ResponseError{
		code:           "RestoreAlreadyInProgress",
		description:    "Object restore is already in progress",
		httpStatusCode: http.StatusConflict,
	}
	// Bucket notification related errors.
	RespErrEventNotification = &ResponseError{
		code:           "InvalidArgument",
		description:    "A specified event is not supported for notifications.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrARNNotification = &ResponseError{
		code:           "InvalidArgument",
		description:    "A specified destination ARN does not exist or is not well-formed. Verify the destination ARN.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrRegionNotification = &ResponseError{
		code:           "InvalidArgument",
		description:    "A specified destination is in a different region than the bucket. You must use a destination that resides in the same region as the bucket.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrOverlappingFilterNotification = &ResponseError{
		code:           "InvalidArgument",
		description:    "An object key name filtering rule defined with overlapping prefixes, overlapping suffixes, or overlapping combinations of prefixes and suffixes for the same event types.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrFilterNameInvalid = &ResponseError{
		code:           "InvalidArgument",
		description:    "filter rule name must be either prefix or suffix",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrFilterNamePrefix = &ResponseError{
		code:           "InvalidArgument",
		description:    "Cannot specify more than one prefix rule in a filter.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrFilterNameSuffix = &ResponseError{
		code:           "InvalidArgument",
		description:    "Cannot specify more than one suffix rule in a filter.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrFilterValueInvalid = &ResponseError{
		code:           "InvalidArgument",
		description:    "Size of filter rule value cannot exceed 1024 bytes in UTF-8 representation",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrOverlappingConfigs = &ResponseError{
		code:           "InvalidArgument",
		description:    "Configurations overlap. Configurations on the same bucket cannot share a common event type.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrContentSHA256Mismatch = &ResponseError{ //todo
		code:           "InvalidArgument",
		description:    "ErrContentSHA256Mismatch",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidCopyPartRange = &ResponseError{
		code:           "InvalidArgument",
		description:    "The x-amz-copy-source-range value must be of the form bytes=first-last where first and last are the zero-based offsets of the first and last bytes to copy",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidCopyPartRangeSource = &ResponseError{
		code:           "InvalidArgument",
		description:    "Range specified is not valid for source object",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrMetadataTooLarge = &ResponseError{
		code:           "MetadataTooLarge",
		description:    "Your metadata headers exceed the maximum allowed metadata size.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidTagDirective = &ResponseError{
		code:           "InvalidArgument",
		description:    "Unknown tag directive.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidEncryptionMethod = &ResponseError{
		code:           "InvalidRequest",
		description:    "The encryption method specified is not supported",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidQueryParams = &ResponseError{
		code:           "AuthorizationQueryParametersError",
		description:    "Query-string authentication version 4 requires the X-Amz-Algorithm, X-Amz-Credential, X-Amz-Signature, X-Amz-Date, X-Amz-SignedHeaders, and X-Amz-Expires parameters.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrNoAccessKey = &ResponseError{
		code:           "AccessDenied",
		description:    "No AWSAccessKey was presented",
		httpStatusCode: http.StatusForbidden,
	}
	RespErrInvalidToken = &ResponseError{
		code:           "InvalidTokenId",
		description:    "The security token included in the request is invalid",
		httpStatusCode: http.StatusForbidden,
	}

	// S3 extensions.
	RespErrInvalidObjectName = &ResponseError{
		code:           "InvalidObjectName",
		description:    "Object name contains unsupported characters.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidObjectNamePrefixSlash = &ResponseError{
		code:           "InvalidObjectName",
		description:    "Object name contains a leading slash.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrClientDisconnected = &ResponseError{
		code:           "ClientDisconnected",
		description:    "Client disconnected before response was ready",
		httpStatusCode: 499, // No official code, use nginx value.
	}
	RespErrOperationTimedOut = &ResponseError{
		code:           "RequestTimeout",
		description:    "A timeout occurred while trying to lock a resource, please reduce your request rate",
		httpStatusCode: http.StatusServiceUnavailable,
	}
	RespErrOperationMaxedOut = &ResponseError{
		code:           "SlowDown",
		description:    "A timeout exceeded while waiting to proceed with the request, please reduce your request rate",
		httpStatusCode: http.StatusServiceUnavailable,
	}
	RespErrUnsupportedMetadata = &ResponseError{
		code:           "InvalidArgument",
		description:    "Your metadata headers are not supported.",
		httpStatusCode: http.StatusBadRequest,
	}
	// Generic Invalid-Request error. Should be used for response errors only for unlikely
	// corner case errors for which introducing new APIRespErrorcode is not worth it. LogIf()
	// should be used to log the error at the source of the error for debugging purposes.
	ErrInvalidRequest = &ResponseError{
		code:           "InvalidRequest",
		description:    "Invalid Request",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrIncorrectContinuationToken = &ResponseError{
		code:           "InvalidArgument",
		description:    "The continuation token provided is incorrect",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidFormatAccessKey = &ResponseError{
		code:           "InvalidAccessKeyId",
		description:    "The Access Key Id you provided contains invalid characters.",
		httpStatusCode: http.StatusBadRequest,
	}
	// S3 Select API RespErrors
	ErrEmptyRequestBody = &ResponseError{
		code:           "EmptyRequestBody",
		description:    "Request body cannot be empty.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrUnsupportedFunction = &ResponseError{
		code:           "UnsupportedFunction",
		description:    "Encountered an unsupported SQL function.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidDataSource = &ResponseError{
		code:           "InvalidDataSource",
		description:    "Invalid data source type. Only CSV and JSON are supported at this time.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidExpressionType = &ResponseError{
		code:           "InvalidExpressionType",
		description:    "The ExpressionType is invalid. Only SQL expressions are supported at this time.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrBusy = &ResponseError{
		code:           "Busy",
		description:    "The service is unavailable. Please retry.",
		httpStatusCode: http.StatusServiceUnavailable,
	}
	RespErrUnauthorizedAccess = &ResponseError{
		code:           "UnauthorizedAccess",
		description:    "You are not authorized to perform this operation",
		httpStatusCode: http.StatusUnauthorized,
	}
	RespErrExpressionTooLong = &ResponseError{
		code:           "ExpressionTooLong",
		description:    "The SQL expression is too long: The maximum byte-length for the SQL expression is 256 KB.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrIllegalSQLFunctionArgument = &ResponseError{
		code:           "IllegalSqlFunctionArgument",
		description:    "Illegal argument was used in the SQL function.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidKeyPath = &ResponseError{
		code:           "InvalidKeyPath",
		description:    "Key path in the SQL expression is invalid.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidCompressionFormat = &ResponseError{
		code:           "InvalidCompressionFormat",
		description:    "The file is not in a supported compression format. Only GZIP is supported at this time.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidFileHeaderInfo = &ResponseError{
		code:           "InvalidFileHeaderInfo",
		description:    "The FileHeaderInfo is invalid. Only NONE, USE, and IGNORE are supported.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidJSONType = &ResponseError{
		code:           "InvalidJsonType",
		description:    "The JsonType is invalid. Only DOCUMENT and LINES are supported at this time.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidQuoteFields = &ResponseError{
		code:           "InvalidQuoteFields",
		description:    "The QuoteFields is invalid. Only ALWAYS and ASNEEDED are supported.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidRequestParameter = &ResponseError{
		code:           "InvalidRequestParameter",
		description:    "The value of a parameter in SelectRequest element is invalid. Check the service API documentation and try again.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidDataType = &ResponseError{
		code:           "InvalidDataType",
		description:    "The SQL expression contains an invalid data type.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidTextEncoding = &ResponseError{
		code:           "InvalidTextEncoding",
		description:    "Invalid encoding type. Only UTF-8 encoding is supported at this time.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidTableAlias = &ResponseError{
		code:           "InvalidTableAlias",
		description:    "The SQL expression contains an invalid table alias.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrMissingRequiredParameter = &ResponseError{
		code:           "MissingRequiredParameter",
		description:    "The SelectRequest entity is missing a required parameter. Check the service documentation and try again.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrObjectSerializationConflict = &ResponseError{
		code:           "ObjectSerializationConflict",
		description:    "The SelectRequest entity can only contain one of CSV or JSON. Check the service documentation and try again.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrUnsupportedSQLOperation = &ResponseError{
		code:           "UnsupportedSqlOperation",
		description:    "Encountered an unsupported SQL operation.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrUnsupportedSQLStructure = &ResponseError{
		code:           "UnsupportedSqlStructure",
		description:    "Encountered an unsupported SQL structure. Check the SQL Reference.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrUnsupportedSyntax = &ResponseError{
		code:           "UnsupportedSyntax",
		description:    "Encountered invalid syntax.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrUnsupportedRangeHeader = &ResponseError{
		code:           "UnsupportedRangeHeader",
		description:    "Range header is not supported for this operation.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrLexerInvalidChar = &ResponseError{
		code:           "LexerInvalidChar",
		description:    "The SQL expression contains an invalid character.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrLexerInvalidOperator = &ResponseError{
		code:           "LexerInvalidOperator",
		description:    "The SQL expression contains an invalid literal.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrLexerInvalidLiteral = &ResponseError{
		code:           "LexerInvalidLiteral",
		description:    "The SQL expression contains an invalid operator.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrLexerInvalidIONLiteral = &ResponseError{
		code:           "LexerInvalidIONLiteral",
		description:    "The SQL expression contains an invalid operator.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseExpectedDatePart = &ResponseError{
		code:           "ParseExpectedDatePart",
		description:    "Did not find the expected date part in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseExpectedKeyword = &ResponseError{
		code:           "ParseExpectedKeyword",
		description:    "Did not find the expected keyword in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseExpectedTokenType = &ResponseError{
		code:           "ParseExpectedTokenType",
		description:    "Did not find the expected token in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseExpected2TokenTypes = &ResponseError{
		code:           "ParseExpected2TokenTypes",
		description:    "Did not find the expected token in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseExpectedNumber = &ResponseError{
		code:           "ParseExpectedNumber",
		description:    "Did not find the expected number in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseExpectedRightParenBuiltinFunctionCall = &ResponseError{
		code:           "ParseExpectedRightParenBuiltinFunctionCall",
		description:    "Did not find the expected right parenthesis character in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseExpectedTypeName = &ResponseError{
		code:           "ParseExpectedTypeName",
		description:    "Did not find the expected type name in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseExpectedWhenClause = &ResponseError{
		code:           "ParseExpectedWhenClause",
		description:    "Did not find the expected WHEN clause in the SQL expression. CASE is not supported.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseUnsupportedToken = &ResponseError{
		code:           "ParseUnsupportedToken",
		description:    "The SQL expression contains an unsupported token.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseUnsupportedLiteralsGroupBy = &ResponseError{
		code:           "ParseUnsupportedLiteralsGroupBy",
		description:    "The SQL expression contains an unsupported use of GROUP BY.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseExpectedMember = &ResponseError{
		code:           "ParseExpectedMember",
		description:    "The SQL expression contains an unsupported use of MEMBER.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseUnsupportedSelect = &ResponseError{
		code:           "ParseUnsupportedSelect",
		description:    "The SQL expression contains an unsupported use of SELECT.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseUnsupportedCase = &ResponseError{
		code:           "ParseUnsupportedCase",
		description:    "The SQL expression contains an unsupported use of CASE.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseUnsupportedCaseClause = &ResponseError{
		code:           "ParseUnsupportedCaseClause",
		description:    "The SQL expression contains an unsupported use of CASE.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseUnsupportedAlias = &ResponseError{
		code:           "ParseUnsupportedAlias",
		description:    "The SQL expression contains an unsupported use of ALIAS.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseUnsupportedSyntax = &ResponseError{
		code:           "ParseUnsupportedSyntax",
		description:    "The SQL expression contains unsupported syntax.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseUnknownOperator = &ResponseError{
		code:           "ParseUnknownOperator",
		description:    "The SQL expression contains an invalid operator.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseMissingIdentAfterAt = &ResponseError{
		code:           "ParseMissingIdentAfterAt",
		description:    "Did not find the expected identifier after the @ symbol in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseUnexpectedOperator = &ResponseError{
		code:           "ParseUnexpectedOperator",
		description:    "The SQL expression contains an unexpected operator.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseUnexpectedTerm = &ResponseError{
		code:           "ParseUnexpectedTerm",
		description:    "The SQL expression contains an unexpected term.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseUnexpectedToken = &ResponseError{
		code:           "ParseUnexpectedToken",
		description:    "The SQL expression contains an unexpected token.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseUnexpectedKeyword = &ResponseError{
		code:           "ParseUnexpectedKeyword",
		description:    "The SQL expression contains an unexpected keyword.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseExpectedExpression = &ResponseError{
		code:           "ParseExpectedExpression",
		description:    "Did not find the expected SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseExpectedLeftParenAfterCast = &ResponseError{
		code:           "ParseExpectedLeftParenAfterCast",
		description:    "Did not find expected the left parenthesis in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseExpectedLeftParenValueConstructor = &ResponseError{
		code:           "ParseExpectedLeftParenValueConstructor",
		description:    "Did not find expected the left parenthesis in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseExpectedLeftParenBuiltinFunctionCall = &ResponseError{
		code:           "ParseExpectedLeftParenBuiltinFunctionCall",
		description:    "Did not find the expected left parenthesis in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseExpectedArgumentDelimiter = &ResponseError{
		code:           "ParseExpectedArgumentDelimiter",
		description:    "Did not find the expected argument delimiter in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseCastArity = &ResponseError{
		code:           "ParseCastArity",
		description:    "The SQL expression CAST has incorrect arity.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseInvalidTypeParam = &ResponseError{
		code:           "ParseInvalidTypeParam",
		description:    "The SQL expression contains an invalid parameter value.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseEmptySelect = &ResponseError{
		code:           "ParseEmptySelect",
		description:    "The SQL expression contains an empty SELECT.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseSelectMissingFrom = &ResponseError{
		code:           "ParseSelectMissingFrom",
		description:    "GROUP is not supported in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseExpectedIdentForGroupName = &ResponseError{
		code:           "ParseExpectedIdentForGroupName",
		description:    "GROUP is not supported in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseExpectedIdentForAlias = &ResponseError{
		code:           "ParseExpectedIdentForAlias",
		description:    "Did not find the expected identifier for the alias in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseUnsupportedCallWithStar = &ResponseError{
		code:           "ParseUnsupportedCallWithStar",
		description:    "Only COUNT with (*) as a parameter is supported in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseNonUnaryAgregateFunctionCall = &ResponseError{
		code:           "ParseNonUnaryAgregateFunctionCall",
		description:    "Only one argument is supported for aggregate functions in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseMalformedJoin = &ResponseError{
		code:           "ParseMalformedJoin",
		description:    "JOIN is not supported in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseExpectedIdentForAt = &ResponseError{
		code:           "ParseExpectedIdentForAt",
		description:    "Did not find the expected identifier for AT name in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseAsteriskIsNotAloneInSelectList = &ResponseError{
		code:           "ParseAsteriskIsNotAloneInSelectList",
		description:    "Other expressions are not allowed in the SELECT list when '*' is used without dot notation in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseCannotMixSqbAndWildcardInSelectList = &ResponseError{
		code:           "ParseCannotMixSqbAndWildcardInSelectList",
		description:    "Cannot mix [] and * in the same expression in a SELECT list in SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrParseInvalidContextForWildcardInSelectList = &ResponseError{
		code:           "ParseInvalidContextForWildcardInSelectList",
		description:    "Invalid use of * in SELECT list in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrIncorrectSQLFunctionArgumentType = &ResponseError{
		code:           "IncorrectSqlFunctionArgumentType",
		description:    "Incorrect type of arguments in function call in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrValueParseFailure = &ResponseError{
		code:           "ValueParseFailure",
		description:    "Time stamp parse failure in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrEvaluatorInvalidArguments = &ResponseError{
		code:           "EvaluatorInvalidArguments",
		description:    "Incorrect number of arguments in the function call in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrIntegerOverflow = &ResponseError{
		code:           "IntegerOverflow",
		description:    "Int overflow or underflow in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrLikeInvalidInputs = &ResponseError{
		code:           "LikeInvalidInputs",
		description:    "Invalid argument given to the LIKE clause in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrCastFailed = &ResponseError{
		code:           "CastFailed",
		description:    "Attempt to convert from one data type to another using CAST failed in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidCast = &ResponseError{
		code:           "InvalidCast",
		description:    "Attempt to convert from one data type to another using CAST failed in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrEvaluatorInvalidTimestampFormatPattern = &ResponseError{
		code:           "EvaluatorInvalidTimestampFormatPattern",
		description:    "Time stamp format pattern requires additional fields in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrEvaluatorInvalidTimestampFormatPatternSymbolForParsing = &ResponseError{
		code:           "EvaluatorInvalidTimestampFormatPatternSymbolForParsing",
		description:    "Time stamp format pattern contains a valid format symbol that cannot be applied to time stamp parsing in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrEvaluatorTimestampFormatPatternDuplicateFields = &ResponseError{
		code:           "EvaluatorTimestampFormatPatternDuplicateFields",
		description:    "Time stamp format pattern contains multiple format specifiers representing the time stamp field in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrEvaluatorTimestampFormatPatternHourClockAmPmMismatch = &ResponseError{
		code:           "EvaluatorUnterminatedTimestampFormatPatternToken",
		description:    "Time stamp format pattern contains unterminated token in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrEvaluatorUnterminatedTimestampFormatPatternToken = &ResponseError{
		code:           "EvaluatorInvalidTimestampFormatPatternToken",
		description:    "Time stamp format pattern contains an invalid token in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrEvaluatorInvalidTimestampFormatPatternToken = &ResponseError{
		code:           "EvaluatorInvalidTimestampFormatPatternToken",
		description:    "Time stamp format pattern contains an invalid token in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrEvaluatorInvalidTimestampFormatPatternSymbol = &ResponseError{
		code:           "EvaluatorInvalidTimestampFormatPatternSymbol",
		description:    "Time stamp format pattern contains an invalid symbol in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrEvaluatorBindingDoesNotExist = &ResponseError{
		code:           "ErrEvaluatorBindingDoesNotExist",
		description:    "A column name or a path provided does not exist in the SQL expression",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrMissingHeaders = &ResponseError{
		code:           "MissingHeaders",
		description:    "Some headers in the query are missing from the file. Check the file and try again.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrInvalidColumnIndex = &ResponseError{
		code:           "InvalidColumnIndex",
		description:    "The column index is invalid. Please check the service documentation and try again.",
		httpStatusCode: http.StatusBadRequest,
	}
	RespErrPostPolicyConditionInvalidFormat = &ResponseError{
		code:           "PostPolicyInvalidKeyName",
		description:    "Invalid according to Policy: Policy Conditions failed",
		httpStatusCode: http.StatusForbidden,
	}
	RespErrMalformedJSON = &ResponseError{
		code:           "MalformedJSON",
		description:    "The JSON was not well-formed or did not validate against our published format.",
		httpStatusCode: http.StatusBadRequest,
	}
)
