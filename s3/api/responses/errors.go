package responses

import (
	"fmt"
	"net/http"
)

type Error struct {
	code           string
	description    string
	httpStatusCode int
}

func (err *Error) Code() string {
	return err.code
}

func (err *Error) Description() string {
	return err.description
}

func (err *Error) HTTPStatusCode() int {
	return err.httpStatusCode
}

func (err *Error) Error() string {
	return fmt.Sprintf("<%s> %s", err.code, err.description)
}

// Errors http://docs.aws.amazon.com/AmazonS3/latest/API/ErrorResponses.html
var (
	ErrInvalidCopyDest = &Error{
		code:           "InvalidRequest",
		description:    "This copy request is illegal because it is trying to copy an object to itself without changing the object's metadata, storage class, website redirect location or encryption attributes.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidCopySource = &Error{
		code:           "InvalidArgument",
		description:    "Copy Source must mention the source bucket and key: sourcebucket/sourcekey.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidMetadataDirective = &Error{
		code:           "InvalidArgument",
		description:    "Unknown metadata directive.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidRequestBody = &Error{
		code:           "InvalidArgument",
		description:    "Body shouldn't be set for this request.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidUploads = &Error{
		code:           "InvalidArgument",
		description:    "Argument max-uploads must be an integer between 0 and 2147483647",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidMaxKeys = &Error{
		code:           "InvalidArgument",
		description:    "Argument maxKeys must be an integer between 0 and 2147483647",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidEncodingMethod = &Error{
		code:           "InvalidArgument",
		description:    "Invalid Encoding Method specified in Request",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidPartNumber = &Error{
		code:           "InvalidArgument",
		description:    "Part number must be an integer between 1 and 10000, inclusive",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidPartNumberMarker = &Error{
		code:           "InvalidArgument",
		description:    "Argument partNumberMarker must be an integer.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidPolicyDocument = &Error{
		code:           "InvalidPolicyDocument",
		description:    "The content of the form does not meet the conditions specified in the policy document.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrAccessDenied = &Error{
		code:           "AccessDenied",
		description:    "Access Denied.",
		httpStatusCode: http.StatusForbidden,
	}
	ErrBadDigest = &Error{
		code:           "BadDigest",
		description:    "The Content-Md5 you specified did not match what we received.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrEntityTooSmall = &Error{
		code:           "EntityTooSmall",
		description:    "Your proposed upload is smaller than the minimum allowed object size.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrEntityTooLarge = &Error{
		code:           "EntityTooLarge",
		description:    "Your proposed upload exceeds the maximum allowed object size.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrIncompleteBody = &Error{
		code:           "IncompleteBody",
		description:    "You did not provide the number of bytes specified by the Content-Length HTTP header.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInternalError = &Error{
		code:           "InternalError",
		description:    "We encountered an internal error, please try again.",
		httpStatusCode: http.StatusInternalServerError,
	}
	ErrInvalidAccessKeyID = &Error{
		code:           "InvalidAccessKeyId",
		description:    "The Access Key Id you provided does not exist in our records.",
		httpStatusCode: http.StatusForbidden,
	}
	ErrAccessKeyDisabled = &Error{
		code:           "InvalidAccessKeyId",
		description:    "Your account is disabled; please contact your administrator.",
		httpStatusCode: http.StatusForbidden,
	}
	ErrInvalidBucketName = &Error{
		code:           "InvalidBucketName",
		description:    "The specified bucket is not valid.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidDigest = &Error{
		code:           "InvalidDigest",
		description:    "The Content-Md5 you specified is not valid.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidRange = &Error{
		code:           "InvalidRange",
		description:    "The requested range is not satisfiable",
		httpStatusCode: http.StatusRequestedRangeNotSatisfiable,
	}
	ErrInvalidRangePartNumber = &Error{
		code:           "InvalidRequest",
		description:    "Cannot specify both Range header and partNumber query parameter",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrMalformedXML = &Error{
		code:           "MalformedXML",
		description:    "The XML you provided was not well-formed or did not validate against our published schema.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrMissingContentLength = &Error{
		code:           "MissingContentLength",
		description:    "You must provide the Content-Length HTTP header.",
		httpStatusCode: http.StatusLengthRequired,
	}
	ErrMissingContentMD5 = &Error{
		code:           "MissingContentMD5",
		description:    "Missing required header for this request: Content-Md5.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrMissingSecurityHeader = &Error{
		code:           "MissingSecurityHeader",
		description:    "Your request was missing a required header",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrMissingRequestBodyError = &Error{
		code:           "MissingRequestBodyError",
		description:    "Request body is empty.",
		httpStatusCode: http.StatusLengthRequired,
	}
	ErrNoSuchBucket = &Error{
		code:           "NoSuchBucket",
		description:    "The specified bucket does not exist",
		httpStatusCode: http.StatusNotFound,
	}
	ErrNoSuchBucketPolicy = &Error{
		code:           "NoSuchBucketPolicy",
		description:    "The bucket policy does not exist",
		httpStatusCode: http.StatusNotFound,
	}
	ErrNoSuchLifecycleConfiguration = &Error{
		code:           "NoSuchLifecycleConfiguration",
		description:    "The lifecycle configuration does not exist",
		httpStatusCode: http.StatusNotFound,
	}
	ErrNoSuchUser = &Error{
		code:           "NoSuchUser",
		description:    "The specified user does not exist",
		httpStatusCode: http.StatusConflict,
	}
	ErrUserAlreadyExists = &Error{
		code:           "UserAlreadyExists",
		description:    "The request was rejected because it attempted to create a resource that already exists .",
		httpStatusCode: http.StatusConflict,
	}
	ErrNoSuchUserPolicy = &Error{
		code:           "NoSuchUserPolicy",
		description:    "The specified user policy does not exist",
		httpStatusCode: http.StatusConflict,
	}
	ErrUserPolicyAlreadyExists = &Error{
		code:           "UserPolicyAlreadyExists",
		description:    "The same user policy already exists .",
		httpStatusCode: http.StatusConflict,
	}
	ErrNoSuchKey = &Error{
		code:           "NoSuchKey",
		description:    "The specified key does not exist.",
		httpStatusCode: http.StatusNotFound,
	}
	ErrNoSuchUpload = &Error{
		code:           "NoSuchUpload",
		description:    "The specified multipart upload does not exist. The upload ID may be invalid, or the upload may have been aborted or completed.",
		httpStatusCode: http.StatusNotFound,
	}
	ErrInvalidVersionID = &Error{
		code:           "InvalidArgument",
		description:    "Invalid version id specified",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrNoSuchVersion = &Error{
		code:           "NoSuchVersion",
		description:    "The specified version does not exist.",
		httpStatusCode: http.StatusNotFound,
	}
	ErrNotImplemented = &Error{
		code:           "NotImplemented",
		description:    "A header you provided implies functionality that is not implemented",
		httpStatusCode: http.StatusNotImplemented,
	}
	ErrPreconditionFailed = &Error{
		code:           "PreconditionFailed",
		description:    "At least one of the pre-conditions you specified did not hold",
		httpStatusCode: http.StatusPreconditionFailed,
	}
	ErrRequestTimeTooSkewed = &Error{
		code:           "RequestTimeTooSkewed",
		description:    "The difference between the request time and the server's time is too large.",
		httpStatusCode: http.StatusForbidden,
	}
	ErrSignatureDoesNotMatch = &Error{
		code:           "SignatureDoesNotMatch",
		description:    "The request signature we calculated does not match the signature you provided. Check your key and signing method.",
		httpStatusCode: http.StatusForbidden,
	}
	ErrMethodNotAllowed = &Error{
		code:           "MethodNotAllowed",
		description:    "The specified method is not allowed against this resource.",
		httpStatusCode: http.StatusMethodNotAllowed,
	}
	ErrInvalidPart = &Error{
		code:           "InvalidPart",
		description:    "One or more of the specified parts could not be found.  The part may not have been uploaded, or the specified entity tag may not match the part's entity tag.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidPartOrder = &Error{
		code:           "InvalidPartOrder",
		description:    "The list of parts was not in ascending order. The parts list must be specified in order by part number.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidObjectState = &Error{
		code:           "InvalidObjectState",
		description:    "The operation is not valid for the current state of the object.",
		httpStatusCode: http.StatusForbidden,
	}
	ErrAuthorizationHeaderMalformed = &Error{
		code:           "AuthorizationHeaderMalformed",
		description:    "The authorization header is malformed; the region is wrong; expecting 'us-east-1'.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrMalformedDate = &Error{ // todo
		code:           "ErrMalformedDate",
		description:    "ErrMalformedDate",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrMalformedPOSTRequest = &Error{
		code:           "MalformedPOSTRequest",
		description:    "The body of your POST request is not well-formed multipart/form-data.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrPOSTFileRequired = &Error{
		code:           "InvalidArgument",
		description:    "POST requires exactly one file upload per request.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrSignatureVersionNotSupported = &Error{
		code:           "InvalidRequest",
		description:    "The authorization mechanism you have provided is not supported. Please use AWS4-HMAC-SHA256.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrBucketNotEmpty = &Error{
		code:           "BucketNotEmpty",
		description:    "The bucket you tried to delete is not empty",
		httpStatusCode: http.StatusConflict,
	}
	ErrBucketAlreadyExists = &Error{
		code:           "BucketAlreadyExists",
		description:    "The requested bucket name is not available. The bucket namespace is shared by all users of the system. Please select a different name and try again.",
		httpStatusCode: http.StatusConflict,
	}
	ErrAllAccessDisabled = &Error{
		code:           "AllAccessDisabled",
		description:    "All access to this resource has been disabled.",
		httpStatusCode: http.StatusForbidden,
	}
	ErrMalformedPolicy = &Error{
		code:           "MalformedPolicy",
		description:    "Policy has invalid resource.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrMissingFields = &Error{ // todo
		code:           "InvalidRequest",
		description:    "ErrMissingFields",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrMissingCredTag = &Error{
		code:           "InvalidRequest",
		description:    "Missing Credential field for this request.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrCredMalformed = &Error{ // todo
		code:           "InvalidRequest",
		description:    "ErrCredMalformed",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidRegion = &Error{
		code:           "InvalidRegion",
		description:    "Region does not match.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrMissingSignTag = &Error{
		code:           "AccessDenied",
		description:    "Signature header missing Signature field.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrMissingSignHeadersTag = &Error{
		code:           "InvalidArgument",
		description:    "Signature header missing SignedHeaders field.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrAuthHeaderEmpty = &Error{
		code:           "InvalidArgument",
		description:    "Authorization header is invalid -- one and only one ' ' (space) required.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrMissingDateHeader = &Error{
		code:           "AccessDenied",
		description:    "AWS authentication requires a valid Date or x-amz-date header",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrExpiredPresignRequest = &Error{
		code:           "AccessDenied",
		description:    "Request has expired",
		httpStatusCode: http.StatusForbidden,
	}
	ErrRequestNotReadyYet = &Error{
		code:           "AccessDenied",
		description:    "Request is not valid yet",
		httpStatusCode: http.StatusForbidden,
	}
	ErrSlowDown = &Error{
		code:           "SlowDown",
		description:    "Resource requested is unreadable, please reduce your request rate",
		httpStatusCode: http.StatusServiceUnavailable,
	}
	ErrBadRequest = &Error{
		code:           "BadRequest",
		description:    "400 BadRequest",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrKeyTooLongError = &Error{
		code:           "KeyTooLongError",
		description:    "Your key is too long",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrUnsignedHeaders = &Error{
		code:           "AccessDenied",
		description:    "There were headers present in the request which were not signed",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrBucketAlreadyOwnedByYou = &Error{
		code:           "BucketAlreadyOwnedByYou",
		description:    "Your previous request to create the named bucket succeeded and you already own it.",
		httpStatusCode: http.StatusConflict,
	}
	ErrInvalidDuration = &Error{
		code:           "InvalidDuration",
		description:    "Duration provided in the request is invalid.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidBucketObjectLockConfiguration = &Error{
		code:           "InvalidRequest",
		description:    "Bucket is missing ObjectLockConfiguration",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrBucketTaggingNotFound = &Error{
		code:           "NoSuchTagSet",
		description:    "The TagSet does not exist",
		httpStatusCode: http.StatusNotFound,
	}
	ErrObjectLockConfigurationNotAllowed = &Error{
		code:           "InvalidBucketState",
		description:    "Object Lock configuration cannot be enabled on existing buckets",
		httpStatusCode: http.StatusConflict,
	}
	ErrNoSuchCORSConfiguration = &Error{
		code:           "NoSuchCORSConfiguration",
		description:    "The CORS configuration does not exist",
		httpStatusCode: http.StatusNotFound,
	}
	ErrNoSuchWebsiteConfiguration = &Error{
		code:           "NoSuchWebsiteConfiguration",
		description:    "The specified bucket does not have a website configuration",
		httpStatusCode: http.StatusNotFound,
	}
	ErrReplicationConfigurationNotFoundError = &Error{
		code:           "ReplicationConfigurationNotFoundError",
		description:    "The replication configuration was not found",
		httpStatusCode: http.StatusNotFound,
	}
	ErrReplicationNeedsVersioningError = &Error{
		code:           "InvalidRequest",
		description:    "Versioning must be 'Enabled' on the bucket to apply a replication configuration",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrReplicationBucketNeedsVersioningError = &Error{
		code:           "InvalidRequest",
		description:    "Versioning must be 'Enabled' on the bucket to add a replication target",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrNoSuchObjectLockConfiguration = &Error{
		code:           "NoSuchObjectLockConfiguration",
		description:    "The specified object does not have a ObjectLock configuration",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrObjectLocked = &Error{
		code:           "InvalidRequest",
		description:    "Object is WORM protected and cannot be overwritten",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidRetentionDate = &Error{
		code:           "InvalidRequest",
		description:    "Date must be provided in ISO 8601 format",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrPastObjectLockRetainDate = &Error{
		code:           "InvalidRequest",
		description:    "the retain until date must be in the future",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrUnknownWORMModeDirective = &Error{
		code:           "InvalidRequest",
		description:    "unknown wormMode directive",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrObjectLockInvalidHeaders = &Error{
		code:           "InvalidRequest",
		description:    "x-amz-object-lock-retain-until-date and x-amz-object-lock-mode must both be supplied",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrObjectRestoreAlreadyInProgress = &Error{
		code:           "RestoreAlreadyInProgress",
		description:    "Object restore is already in progress",
		httpStatusCode: http.StatusConflict,
	}
	// Bucket notification related errors.
	ErrEventNotification = &Error{
		code:           "InvalidArgument",
		description:    "A specified event is not supported for notifications.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrARNNotification = &Error{
		code:           "InvalidArgument",
		description:    "A specified destination ARN does not exist or is not well-formed. Verify the destination ARN.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrRegionNotification = &Error{
		code:           "InvalidArgument",
		description:    "A specified destination is in a different region than the bucket. You must use a destination that resides in the same region as the bucket.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrOverlappingFilterNotification = &Error{
		code:           "InvalidArgument",
		description:    "An object key name filtering rule defined with overlapping prefixes, overlapping suffixes, or overlapping combinations of prefixes and suffixes for the same event types.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrFilterNameInvalid = &Error{
		code:           "InvalidArgument",
		description:    "filter rule name must be either prefix or suffix",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrFilterNamePrefix = &Error{
		code:           "InvalidArgument",
		description:    "Cannot specify more than one prefix rule in a filter.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrFilterNameSuffix = &Error{
		code:           "InvalidArgument",
		description:    "Cannot specify more than one suffix rule in a filter.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrFilterValueInvalid = &Error{
		code:           "InvalidArgument",
		description:    "Size of filter rule value cannot exceed 1024 bytes in UTF-8 representation",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrOverlappingConfigs = &Error{
		code:           "InvalidArgument",
		description:    "Configurations overlap. Configurations on the same bucket cannot share a common event type.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrContentSHA256Mismatch = &Error{ //todo
		code:           "InvalidArgument",
		description:    "ErrContentSHA256Mismatch",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidCopyPartRange = &Error{
		code:           "InvalidArgument",
		description:    "The x-amz-copy-source-range value must be of the form bytes=first-last where first and last are the zero-based offsets of the first and last bytes to copy",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidCopyPartRangeSource = &Error{
		code:           "InvalidArgument",
		description:    "Range specified is not valid for source object",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrMetadataTooLarge = &Error{
		code:           "MetadataTooLarge",
		description:    "Your metadata headers exceed the maximum allowed metadata size.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidTagDirective = &Error{
		code:           "InvalidArgument",
		description:    "Unknown tag directive.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidEncryptionMethod = &Error{
		code:           "InvalidRequest",
		description:    "The encryption method specified is not supported",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidQueryParams = &Error{
		code:           "AuthorizationQueryParametersError",
		description:    "Query-string authentication version 4 requires the X-Amz-Algorithm, X-Amz-Credential, X-Amz-Signature, X-Amz-Date, X-Amz-SignedHeaders, and X-Amz-Expires parameters.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrNoAccessKey = &Error{
		code:           "AccessDenied",
		description:    "No AWSAccessKey was presented",
		httpStatusCode: http.StatusForbidden,
	}
	ErrInvalidToken = &Error{
		code:           "InvalidTokenId",
		description:    "The security token included in the request is invalid",
		httpStatusCode: http.StatusForbidden,
	}

	// S3 extensions.
	ErrInvalidObjectName = &Error{
		code:           "InvalidObjectName",
		description:    "Object name contains unsupported characters.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidObjectNamePrefixSlash = &Error{
		code:           "InvalidObjectName",
		description:    "Object name contains a leading slash.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrClientDisconnected = &Error{
		code:           "ClientDisconnected",
		description:    "Client disconnected before response was ready",
		httpStatusCode: 499, // No official code, use nginx value.
	}
	ErrOperationTimedOut = &Error{
		code:           "RequestTimeout",
		description:    "A timeout occurred while trying to lock a resource, please reduce your request rate",
		httpStatusCode: http.StatusServiceUnavailable,
	}
	ErrOperationMaxedOut = &Error{
		code:           "SlowDown",
		description:    "A timeout exceeded while waiting to proceed with the request, please reduce your request rate",
		httpStatusCode: http.StatusServiceUnavailable,
	}
	ErrUnsupportedMetadata = &Error{
		code:           "InvalidArgument",
		description:    "Your metadata headers are not supported.",
		httpStatusCode: http.StatusBadRequest,
	}
	// Generic Invalid-Request error. Should be used for response errors only for unlikely
	// corner case errors for which introducing new APIorcode is not worth it. LogIf()
	// should be used to log the error at the source of the error for debugging purposes.
	ErrInvalidRequest = &Error{
		code:           "InvalidRequest",
		description:    "Invalid Request",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrIncorrectContinuationToken = &Error{
		code:           "InvalidArgument",
		description:    "The continuation token provided is incorrect",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidFormatAccessKey = &Error{
		code:           "InvalidAccessKeyId",
		description:    "The Access Key Id you provided contains invalid characters.",
		httpStatusCode: http.StatusBadRequest,
	}
	// S3 Select API ors
	ErrErrEmptyRequestBody = &Error{
		code:           "EmptyRequestBody",
		description:    "Request body cannot be empty.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrUnsupportedFunction = &Error{
		code:           "UnsupportedFunction",
		description:    "Encountered an unsupported SQL function.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidDataSource = &Error{
		code:           "InvalidDataSource",
		description:    "Invalid data source type. Only CSV and JSON are supported at this time.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidExpressionType = &Error{
		code:           "InvalidExpressionType",
		description:    "The ExpressionType is invalid. Only SQL expressions are supported at this time.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrBusy = &Error{
		code:           "Busy",
		description:    "The service is unavailable. Please retry.",
		httpStatusCode: http.StatusServiceUnavailable,
	}
	ErrUnauthorizedAccess = &Error{
		code:           "UnauthorizedAccess",
		description:    "You are not authorized to perform this operation",
		httpStatusCode: http.StatusUnauthorized,
	}
	ErrExpressionTooLong = &Error{
		code:           "ExpressionTooLong",
		description:    "The SQL expression is too long: The maximum byte-length for the SQL expression is 256 KB.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrIllegalSQLFunctionArgument = &Error{
		code:           "IllegalSqlFunctionArgument",
		description:    "Illegal argument was used in the SQL function.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidKeyPath = &Error{
		code:           "InvalidKeyPath",
		description:    "Key path in the SQL expression is invalid.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidCompressionFormat = &Error{
		code:           "InvalidCompressionFormat",
		description:    "The file is not in a supported compression format. Only GZIP is supported at this time.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidFileHeaderInfo = &Error{
		code:           "InvalidFileHeaderInfo",
		description:    "The FileHeaderInfo is invalid. Only NONE, USE, and IGNORE are supported.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidJSONType = &Error{
		code:           "InvalidJsonType",
		description:    "The JsonType is invalid. Only DOCUMENT and LINES are supported at this time.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidQuoteFields = &Error{
		code:           "InvalidQuoteFields",
		description:    "The QuoteFields is invalid. Only ALWAYS and ASNEEDED are supported.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidRequestParameter = &Error{
		code:           "InvalidRequestParameter",
		description:    "The value of a parameter in SelectRequest element is invalid. Check the service API documentation and try again.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidDataType = &Error{
		code:           "InvalidDataType",
		description:    "The SQL expression contains an invalid data type.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidTextEncoding = &Error{
		code:           "InvalidTextEncoding",
		description:    "Invalid encoding type. Only UTF-8 encoding is supported at this time.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidTableAlias = &Error{
		code:           "InvalidTableAlias",
		description:    "The SQL expression contains an invalid table alias.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrMissingRequiredParameter = &Error{
		code:           "MissingRequiredParameter",
		description:    "The SelectRequest entity is missing a required parameter. Check the service documentation and try again.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrObjectSerializationConflict = &Error{
		code:           "ObjectSerializationConflict",
		description:    "The SelectRequest entity can only contain one of CSV or JSON. Check the service documentation and try again.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrUnsupportedSQLOperation = &Error{
		code:           "UnsupportedSqlOperation",
		description:    "Encountered an unsupported SQL operation.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrUnsupportedSQLStructure = &Error{
		code:           "UnsupportedSqlStructure",
		description:    "Encountered an unsupported SQL structure. Check the SQL Reference.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrUnsupportedSyntax = &Error{
		code:           "UnsupportedSyntax",
		description:    "Encountered invalid syntax.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrUnsupportedRangeHeader = &Error{
		code:           "UnsupportedRangeHeader",
		description:    "Range header is not supported for this operation.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrLexerInvalidChar = &Error{
		code:           "LexerInvalidChar",
		description:    "The SQL expression contains an invalid character.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrLexerInvalidOperator = &Error{
		code:           "LexerInvalidOperator",
		description:    "The SQL expression contains an invalid literal.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrLexerInvalidLiteral = &Error{
		code:           "LexerInvalidLiteral",
		description:    "The SQL expression contains an invalid operator.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrLexerInvalidIONLiteral = &Error{
		code:           "LexerInvalidIONLiteral",
		description:    "The SQL expression contains an invalid operator.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseExpectedDatePart = &Error{
		code:           "ParseExpectedDatePart",
		description:    "Did not find the expected date part in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseExpectedKeyword = &Error{
		code:           "ParseExpectedKeyword",
		description:    "Did not find the expected keyword in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseExpectedTokenType = &Error{
		code:           "ParseExpectedTokenType",
		description:    "Did not find the expected token in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseExpected2TokenTypes = &Error{
		code:           "ParseExpected2TokenTypes",
		description:    "Did not find the expected token in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseExpectedNumber = &Error{
		code:           "ParseExpectedNumber",
		description:    "Did not find the expected number in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseExpectedRightParenBuiltinFunctionCall = &Error{
		code:           "ParseExpectedRightParenBuiltinFunctionCall",
		description:    "Did not find the expected right parenthesis character in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseExpectedTypeName = &Error{
		code:           "ParseExpectedTypeName",
		description:    "Did not find the expected type name in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseExpectedWhenClause = &Error{
		code:           "ParseExpectedWhenClause",
		description:    "Did not find the expected WHEN clause in the SQL expression. CASE is not supported.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseUnsupportedToken = &Error{
		code:           "ParseUnsupportedToken",
		description:    "The SQL expression contains an unsupported token.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseUnsupportedLiteralsGroupBy = &Error{
		code:           "ParseUnsupportedLiteralsGroupBy",
		description:    "The SQL expression contains an unsupported use of GROUP BY.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseExpectedMember = &Error{
		code:           "ParseExpectedMember",
		description:    "The SQL expression contains an unsupported use of MEMBER.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseUnsupportedSelect = &Error{
		code:           "ParseUnsupportedSelect",
		description:    "The SQL expression contains an unsupported use of SELECT.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseUnsupportedCase = &Error{
		code:           "ParseUnsupportedCase",
		description:    "The SQL expression contains an unsupported use of CASE.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseUnsupportedCaseClause = &Error{
		code:           "ParseUnsupportedCaseClause",
		description:    "The SQL expression contains an unsupported use of CASE.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseUnsupportedAlias = &Error{
		code:           "ParseUnsupportedAlias",
		description:    "The SQL expression contains an unsupported use of ALIAS.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseUnsupportedSyntax = &Error{
		code:           "ParseUnsupportedSyntax",
		description:    "The SQL expression contains unsupported syntax.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseUnknownOperator = &Error{
		code:           "ParseUnknownOperator",
		description:    "The SQL expression contains an invalid operator.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseMissingIdentAfterAt = &Error{
		code:           "ParseMissingIdentAfterAt",
		description:    "Did not find the expected identifier after the @ symbol in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseUnexpectedOperator = &Error{
		code:           "ParseUnexpectedOperator",
		description:    "The SQL expression contains an unexpected operator.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseUnexpectedTerm = &Error{
		code:           "ParseUnexpectedTerm",
		description:    "The SQL expression contains an unexpected term.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseUnexpectedToken = &Error{
		code:           "ParseUnexpectedToken",
		description:    "The SQL expression contains an unexpected token.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseUnexpectedKeyword = &Error{
		code:           "ParseUnexpectedKeyword",
		description:    "The SQL expression contains an unexpected keyword.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseExpectedExpression = &Error{
		code:           "ParseExpectedExpression",
		description:    "Did not find the expected SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseExpectedLeftParenAfterCast = &Error{
		code:           "ParseExpectedLeftParenAfterCast",
		description:    "Did not find expected the left parenthesis in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseExpectedLeftParenValueConstructor = &Error{
		code:           "ParseExpectedLeftParenValueConstructor",
		description:    "Did not find expected the left parenthesis in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseExpectedLeftParenBuiltinFunctionCall = &Error{
		code:           "ParseExpectedLeftParenBuiltinFunctionCall",
		description:    "Did not find the expected left parenthesis in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseExpectedArgumentDelimiter = &Error{
		code:           "ParseExpectedArgumentDelimiter",
		description:    "Did not find the expected argument delimiter in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseCastArity = &Error{
		code:           "ParseCastArity",
		description:    "The SQL expression CAST has incorrect arity.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseInvalidTypeParam = &Error{
		code:           "ParseInvalidTypeParam",
		description:    "The SQL expression contains an invalid parameter value.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseEmptySelect = &Error{
		code:           "ParseEmptySelect",
		description:    "The SQL expression contains an empty SELECT.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseSelectMissingFrom = &Error{
		code:           "ParseSelectMissingFrom",
		description:    "GROUP is not supported in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseExpectedIdentForGroupName = &Error{
		code:           "ParseExpectedIdentForGroupName",
		description:    "GROUP is not supported in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseExpectedIdentForAlias = &Error{
		code:           "ParseExpectedIdentForAlias",
		description:    "Did not find the expected identifier for the alias in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseUnsupportedCallWithStar = &Error{
		code:           "ParseUnsupportedCallWithStar",
		description:    "Only COUNT with (*) as a parameter is supported in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseNonUnaryAgregateFunctionCall = &Error{
		code:           "ParseNonUnaryAgregateFunctionCall",
		description:    "Only one argument is supported for aggregate functions in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseMalformedJoin = &Error{
		code:           "ParseMalformedJoin",
		description:    "JOIN is not supported in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseExpectedIdentForAt = &Error{
		code:           "ParseExpectedIdentForAt",
		description:    "Did not find the expected identifier for AT name in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseAsteriskIsNotAloneInSelectList = &Error{
		code:           "ParseAsteriskIsNotAloneInSelectList",
		description:    "Other expressions are not allowed in the SELECT list when '*' is used without dot notation in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseCannotMixSqbAndWildcardInSelectList = &Error{
		code:           "ParseCannotMixSqbAndWildcardInSelectList",
		description:    "Cannot mix [] and * in the same expression in a SELECT list in SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrParseInvalidContextForWildcardInSelectList = &Error{
		code:           "ParseInvalidContextForWildcardInSelectList",
		description:    "Invalid use of * in SELECT list in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrIncorrectSQLFunctionArgumentType = &Error{
		code:           "IncorrectSqlFunctionArgumentType",
		description:    "Incorrect type of arguments in function call in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrValueParseFailure = &Error{
		code:           "ValueParseFailure",
		description:    "Time stamp parse failure in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrEvaluatorInvalidArguments = &Error{
		code:           "EvaluatorInvalidArguments",
		description:    "Incorrect number of arguments in the function call in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrIntegerOverflow = &Error{
		code:           "IntegerOverflow",
		description:    "Int overflow or underflow in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrLikeInvalidInputs = &Error{
		code:           "LikeInvalidInputs",
		description:    "Invalid argument given to the LIKE clause in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrCastFailed = &Error{
		code:           "CastFailed",
		description:    "Attempt to convert from one data type to another using CAST failed in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidCast = &Error{
		code:           "InvalidCast",
		description:    "Attempt to convert from one data type to another using CAST failed in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrEvaluatorInvalidTimestampFormatPattern = &Error{
		code:           "EvaluatorInvalidTimestampFormatPattern",
		description:    "Time stamp format pattern requires additional fields in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrEvaluatorInvalidTimestampFormatPatternSymbolForParsing = &Error{
		code:           "EvaluatorInvalidTimestampFormatPatternSymbolForParsing",
		description:    "Time stamp format pattern contains a valid format symbol that cannot be applied to time stamp parsing in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrEvaluatorTimestampFormatPatternDuplicateFields = &Error{
		code:           "EvaluatorTimestampFormatPatternDuplicateFields",
		description:    "Time stamp format pattern contains multiple format specifiers representing the time stamp field in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrEvaluatorTimestampFormatPatternHourClockAmPmMismatch = &Error{
		code:           "EvaluatorUnterminatedTimestampFormatPatternToken",
		description:    "Time stamp format pattern contains unterminated token in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrEvaluatorUnterminatedTimestampFormatPatternToken = &Error{
		code:           "EvaluatorInvalidTimestampFormatPatternToken",
		description:    "Time stamp format pattern contains an invalid token in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrEvaluatorInvalidTimestampFormatPatternToken = &Error{
		code:           "EvaluatorInvalidTimestampFormatPatternToken",
		description:    "Time stamp format pattern contains an invalid token in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrEvaluatorInvalidTimestampFormatPatternSymbol = &Error{
		code:           "EvaluatorInvalidTimestampFormatPatternSymbol",
		description:    "Time stamp format pattern contains an invalid symbol in the SQL expression.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrEvaluatorBindingDoesNotExist = &Error{
		code:           "ErrEvaluatorBindingDoesNotExist",
		description:    "A column name or a path provided does not exist in the SQL expression",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrMissingHeaders = &Error{
		code:           "MissingHeaders",
		description:    "Some headers in the query are missing from the file. Check the file and try again.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrInvalidColumnIndex = &Error{
		code:           "InvalidColumnIndex",
		description:    "The column index is invalid. Please check the service documentation and try again.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrPostPolicyConditionInvalidFormat = &Error{
		code:           "PostPolicyInvalidKeyName",
		description:    "Invalid according to Policy: Policy Conditions failed",
		httpStatusCode: http.StatusForbidden,
	}
	ErrMalformedJSON = &Error{
		code:           "MalformedJSON",
		description:    "The JSON was not well-formed or did not validate against our published format.",
		httpStatusCode: http.StatusBadRequest,
	}
	ErrMalformedACLError = &Error{
		code:           "MalformedACLError",
		description:    "The ACL that you provided was not well formed or did not validate against our published schema.",
		httpStatusCode: http.StatusBadRequest,
	}
)
