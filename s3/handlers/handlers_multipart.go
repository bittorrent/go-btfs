package handlers

//func (h *Handlers) CreateMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
//	ctx := r.Context()
//	ack := cctx.GetAccessKey(r)
//	var err error
//	defer func() {
//		cctx.SetHandleInf(r, h.name(), err)
//	}()
//
//	bucname, objname, err := requests.ParseBucketAndObject(r)
//	if err != nil {
//		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestParameter)
//		return
//	}
//
//	err = s3utils.CheckNewMultipartArgs(ctx, bucname, objname)
//	if err != nil {
//		responses.WriteErrorResponse(w, r, err)
//		return
//	}
//
//	meta, err := extractMetadata(ctx, r)
//	if err != nil {
//		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequest)
//		return
//	}
//
//	// rlock bucket
//	runlock, err := h.rlock(ctx, bucname, w, r)
//	if err != nil {
//		return
//	}
//	defer runlock()
//
//	// lock object
//	unlock, err := h.lock(ctx, bucname+"/"+objname, w, r)
//	if err != nil {
//		return
//	}
//	defer unlock()
//
//	err = h.bucsvc.CheckACL(ack, bucname, action.CreateMultipartUploadAction)
//	if errors.Is(err, object.ErrBucketNotFound) {
//		responses.WriteErrorResponse(w, r, responses.ErrNoSuchBucket)
//		return
//	}
//	if err != nil {
//		responses.WriteErrorResponse(w, r, err)
//		return
//	}
//
//	mtp, err := h.objsvc.CreateMultipartUpload(ctx, bucname, objname, meta)
//	if err != nil {
//		responses.WriteErrorResponse(w, r, err)
//		return
//	}
//
//	responses.WriteCreateMultipartUploadResponse(w, r, bucname, objname, mtp.UploadID)
//
//	return
//}
//
//func (h *Handlers) UploadPartHandler(w http.ResponseWriter, r *http.Request) {
//	ctx := r.Context()
//	ack := cctx.GetAccessKey(r)
//	var err error
//	defer func() {
//		cctx.SetHandleInf(r, h.name(), err)
//	}()
//
//	// X-Amz-Copy-Source shouldn't be set for this call.
//	if _, ok := r.Header[consts.AmzCopySource]; ok {
//		err = errors.New("shouldn't be copy")
//		responses.WriteErrorResponse(w, r, responses.ErrInvalidCopySource)
//		return
//	}
//
//	bucname, objname, err := requests.ParseBucketAndObject(r)
//	if err != nil {
//		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestParameter)
//		return
//	}
//
//	err = s3utils.CheckPutObjectPartArgs(ctx, bucname, objname)
//	if err != nil { // todo: convert error
//		responses.WriteErrorResponse(w, r, err)
//		return
//	}
//
//	uploadID := r.Form.Get(consts.UploadID)
//	partIDString := r.Form.Get(consts.PartNumber)
//	partID, err := strconv.Atoi(partIDString)
//	if err != nil {
//		responses.WriteErrorResponse(w, r, responses.ErrInvalidPart)
//		return
//	}
//	if partID > consts.MaxPartID {
//		responses.WriteErrorResponse(w, r, responses.ErrInvalidMaxParts)
//		return
//	}
//
//	if r.ContentLength == 0 {
//		responses.WriteErrorResponse(w, r, responses.ErrEntityTooSmall)
//		return
//	}
//
//	if r.ContentLength > consts.MaxPartSize {
//		responses.WriteErrorResponse(w, r, responses.ErrEntityTooLarge)
//		return
//	}
//
//	hrdr, ok := r.Body.(*hash.Reader)
//	if !ok {
//		responses.WriteErrorResponse(w, r, responses.ErrInternalError)
//		return
//	}
//
//	mtp, err := h.objsvc.GetMultipart(ctx, bucname, objname, uploadID)
//	if errors.Is(err, object.ErrUploadNotFound) {
//		responses.WriteErrorResponse(w, r, responses.ErrNoSuchUpload)
//		return
//	}
//	if err != nil {
//		responses.WriteErrorResponse(w, r, err)
//		return
//	}
//
//	// rlock bucket
//	runlock, err := h.rlock(ctx, bucname, w, r)
//	if err != nil {
//		return
//	}
//	defer runlock()
//
//	// lock object
//	unlock, err := h.lock(ctx, bucname+"/"+objname, w, r)
//	if err != nil {
//		return
//	}
//	defer unlock()
//
//	err = h.bucsvc.CheckACL(ack, bucname, action.PutObjectAction)
//	if errors.Is(err, object.ErrBucketNotFound) {
//		responses.WriteErrorResponse(w, r, responses.ErrNoSuchBucket)
//		return
//	}
//	if err != nil {
//		responses.WriteErrorResponse(w, r, err)
//		return
//	}
//
//	part, err := h.objsvc.UploadPart(ctx, bucname, objname, uploadID, partID, hrdr, r.ContentLength, mtp.MetaData)
//	if err != nil {
//		responses.WriteErrorResponse(w, r, err)
//		return
//	}
//
//	responses.WriteUploadPartResponse(w, r, part)
//
//	return
//}
//
//func (h *Handlers) AbortMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
//	ctx := r.Context()
//	ack := cctx.GetAccessKey(r)
//	var err error
//	defer func() {
//		cctx.SetHandleInf(r, h.name(), err)
//	}()
//
//	bucname, objname, err := requests.ParseBucketAndObject(r)
//	if err != nil {
//		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestParameter)
//		return
//	}
//
//	err = s3utils.CheckAbortMultipartArgs(ctx, bucname, objname)
//	if err != nil {
//		responses.WriteErrorResponse(w, r, err)
//		return
//	}
//
//	uploadID, _, _, _, rerr := h.getObjectResources(r.Form)
//	if rerr != nil {
//		err = rerr
//		responses.WriteErrorResponse(w, r, err)
//		return
//	}
//
//	// rlock bucket
//	runlock, err := h.rlock(ctx, bucname, w, r)
//	if err != nil {
//		return
//	}
//	defer runlock()
//
//	// rlock object
//	unlock, err := h.lock(ctx, bucname+"/"+objname, w, r)
//	if err != nil {
//		return
//	}
//	defer unlock()
//
//	err = h.bucsvc.CheckACL(ack, bucname, action.AbortMultipartUploadAction)
//	if errors.Is(err, object.ErrBucketNotFound) {
//		responses.WriteErrorResponse(w, r, responses.ErrNoSuchBucket)
//		return
//	}
//	if err != nil {
//		responses.WriteErrorResponse(w, r, err)
//		return
//	}
//
//	err = h.objsvc.AbortMultipartUpload(ctx, bucname, objname, uploadID)
//	if err != nil {
//		responses.WriteErrorResponse(w, r, err)
//		return
//	}
//
//	responses.WriteAbortMultipartUploadResponse(w, r)
//
//	return
//}
//
//func (h *Handlers) CompleteMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
//	ctx := r.Context()
//	ack := cctx.GetAccessKey(r)
//	var err error
//	defer func() {
//		cctx.SetHandleInf(r, h.name(), err)
//	}()
//
//	bucname, objname, err := requests.ParseBucketAndObject(r)
//	if err != nil {
//		responses.WriteErrorResponse(w, r, responses.ErrInvalidRequestParameter)
//		return
//	}
//
//	err = s3utils.CheckCompleteMultipartArgs(ctx, bucname, objname)
//	if err != nil { // todo: convert error
//		responses.WriteErrorResponse(w, r, err)
//		return
//	}
//
//	// Content-Length is required and should be non-zero
//	if r.ContentLength <= 0 {
//		responses.WriteErrorResponse(w, r, responses.ErrMissingContentLength)
//		return
//	}
//
//	// Get upload id.
//	uploadID, _, _, _, rerr := h.getObjectResources(r.Form)
//	if rerr != nil {
//		err = rerr
//		responses.WriteErrorResponse(w, r, rerr)
//		return
//	}
//
//	complMultipartUpload := &object.CompleteMultipartUpload{}
//	if err = utils.XmlDecoder(r.Body, complMultipartUpload, r.ContentLength); err != nil {
//		responses.WriteErrorResponse(w, r, responses.ErrMalformedXML)
//		return
//	}
//	if len(complMultipartUpload.Parts) == 0 {
//		responses.WriteErrorResponse(w, r, responses.ErrMalformedXML)
//		return
//	}
//	if !sort.IsSorted(object.CompletedParts(complMultipartUpload.Parts)) {
//		responses.WriteErrorResponse(w, r, responses.ErrInvalidPartOrder)
//		return
//	}
//
//	// rlock bucket
//	runlock, err := h.rlock(ctx, bucname, w, r)
//	if err != nil {
//		return
//	}
//	defer runlock()
//
//	// rlock object
//	unlock, err := h.lock(ctx, bucname+"/"+objname, w, r)
//	if err != nil {
//		return
//	}
//	defer unlock()
//
//	err = h.bucsvc.CheckACL(ack, bucname, action.CompleteMultipartUploadAction)
//	if errors.Is(err, object.ErrBucketNotFound) {
//		responses.WriteErrorResponse(w, r, responses.ErrNoSuchBucket)
//		return
//	}
//	if err != nil {
//		responses.WriteErrorResponse(w, r, err)
//		return
//	}
//
//	obj, err := h.objsvc.CompleteMultiPartUpload(ctx, bucname, objname, uploadID, complMultipartUpload.Parts)
//	if errors.Is(err, object.ErrUploadNotFound) {
//		rerr = responses.ErrNoSuchUpload
//		return
//	}
//	if err != nil {
//		responses.WriteErrorResponse(w, r, err)
//		return
//	}
//
//	buc, err := h.bucsvc.GetBucketMeta(ctx, bucname)
//	if err != nil {
//		responses.WriteErrorResponse(w, r, err)
//		return
//	}
//
//	responses.WriteCompleteMultipartUploadResponse(w, r, bucname, objname, buc.Region, obj)
//
//	return
//}
//
