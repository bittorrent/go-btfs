package responses

//func WriteCreateMultipartUploadResponse(w http.ResponseWriter, r *http.Request, bucname, objname, uploadID string) {
//	resp := GenerateInitiateMultipartUploadResponse(bucname, objname, uploadID)
//	WriteSuccessResponse(w, resp, "")
//}
//
//func WriteAbortMultipartUploadResponse(w http.ResponseWriter, r *http.Request) {
//	WriteSuccessResponse(w, nil, "")
//}
//
//func WriteUploadPartResponse(w http.ResponseWriter, r *http.Request, part object.Part) {
//	setPutObjHeaders(w, part.ETag, part.CID, false)
//	WriteSuccessResponse(w, nil, "")
//}
//
//func WriteCompleteMultipartUploadResponse(w http.ResponseWriter, r *http.Request, bucname, objname, region string, obj object.Object) {
//	resp := GenerateCompleteMultipartUploadResponse(bucname, objname, region, obj)
//	setPutObjHeaders(w, obj.ETag, obj.CID, false)
//	WriteSuccessResponse(w, resp, "")
//}
