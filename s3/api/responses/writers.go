package responses

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/private/protocol/xml/xmlutil"
	"github.com/bittorrent/go-btfs/s3/consts"
	"github.com/bittorrent/go-btfs/s3/utils"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/private/protocol"
)

const (
	mimeTypeXml = "application/xml"
	noPayload   = "nopayload"
)

var errValueNotSet = fmt.Errorf("value not set")

func WriteSuccessResponse(w http.ResponseWriter, output interface{}, locationName string) {
	_ = WriteResponse(w, http.StatusOK, output, locationName)
}

type ErrorOutput struct {
	_         struct{} `type:"structure"`
	Code      string   `locationName:"Code"`
	Message   string   `locationName:"Message"`
	Resource  string   `locationName:"Resource"`
	RequestID string   `locationName:"RequestID"`
}

func WriteErrorResponse(w http.ResponseWriter, r *http.Request, rerr *Error) {
	_ = WriteResponse(w, rerr.HTTPStatusCode(), &ErrorOutput{
		Code:      rerr.Code(),
		Message:   rerr.Description(),
		Resource:  r.URL.Path,
		RequestID: "",
	}, "Error")
}

func WriteResponse(w http.ResponseWriter, statusCode int, output interface{}, locationName string) (err error) {
	setCommonHeaders(w.Header())

	outv := reflect.Indirect(reflect.ValueOf(wrapOutput(output, locationName)))
	if !outv.IsValid() {
		w.WriteHeader(statusCode)
		return
	}

	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	err = setFieldRequestID(w.Header(), outv)
	if err != nil {
		return
	}

	body, clen, ctyp, err := extractBody(outv)
	if err != nil {
		return
	}

	if body == nil {
		err = setLocationHeaders(w.Header(), outv)
		if err != nil {
			return
		}
		w.WriteHeader(statusCode)
		return
	}

	defer body.Close()

	w.Header().Set(consts.ContentLength, fmt.Sprintf("%d", clen))
	w.Header().Set(consts.ContentType, ctyp)

	err = setLocationHeaders(w.Header(), outv)
	if err != nil {
		return
	}

	w.WriteHeader(statusCode)

	_, err = io.Copy(w, body)

	return
}

func wrapOutput(output interface{}, locationName string) (wrapper interface{}) {
	if locationName == "" {
		wrapper = output
		return
	}

	outputTag := fmt.Sprintf(`locationName:"%s" type:"structure"`, locationName)
	fields := []reflect.StructField{
		{
			Name:    "_",
			Type:    reflect.TypeOf(struct{}{}),
			Tag:     `payload:"Output" type:"structure"`,
			PkgPath: "responses",
		},
		{
			Name: "Output",
			Type: reflect.TypeOf(output),
			Tag:  reflect.StructTag(outputTag),
		},
	}
	wrtyp := reflect.StructOf(fields)
	wrval := reflect.New(wrtyp)
	wrval.Elem().FieldByName("Output").Set(reflect.ValueOf(output))
	wrapper = wrval.Interface()
	return
}

func extractBody(v reflect.Value) (body io.ReadCloser, clen int, ctyp string, err error) {
	ptyp, pfvl := getPayload(v)
	if ptyp == noPayload {
		return
	}

	if ptyp == "structure" || ptyp == "" {
		var buf bytes.Buffer
		buf.WriteString(xml.Header)
		err = xmlutil.BuildXML(v.Interface(), xml.NewEncoder(&buf))
		if err != nil {
			return
		}
		body = io.NopCloser(&buf)
		clen = buf.Len()
		ctyp = mimeTypeXml
		return
	}

	if pfvl.Interface() == nil {
		return
	}

	clen = -1

	body, ok := pfvl.Interface().(io.ReadCloser)
	if !ok {
		err = fmt.Errorf("unsupported payload type <%s>", pfvl.Type())
	}

	return
}

func setFieldRequestID(headers http.Header, outv reflect.Value) (err error) {
	reqId := headers.Get(consts.AmzRequestID)

	idv := outv.FieldByName("RequestID")
	if !idv.IsValid() {
		return
	}

	switch idv.Interface().(type) {
	case *string:
		idv.Set(reflect.ValueOf(&reqId))
	case string:
		idv.Set(reflect.ValueOf(reqId))
	default:
		err = errValueNotSet
	}

	return
}

func setCommonHeaders(headers http.Header) {
	headers.Set(consts.ServerInfo, consts.DefaultServerInfo)
	headers.Set(consts.AcceptRanges, "bytes")
	headers.Set(consts.AmzRequestID, getRequestID())
}

func getRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func setLocationHeaders(header http.Header, v reflect.Value) (err error) {
	for i := 0; i < v.NumField(); i++ {
		fv := reflect.Indirect(v.Field(i))
		ft := v.Type().Field(i)

		if n := ft.Name; n[0:1] == strings.ToLower(n[0:1]) {
			continue
		}

		if !fv.IsValid() {
			continue
		}

		if fv.Kind() == reflect.Interface && !fv.Elem().IsValid() {
			continue
		}

		switch ft.Tag.Get("location") {
		case "header":
			name := utils.CoalesceStr(ft.Tag.Get("locationName"), ft.Name)
			err = setHeaders(&header, fv, name, ft.Tag)
		case "headers":
			err = setHeadersMap(&header, fv, ft.Tag)
		}

		if err != nil {
			return
		}
	}

	return
}

func setHeaders(header *http.Header, v reflect.Value, name string, tag reflect.StructTag) (err error) {
	str, err := convertType(v, tag)
	if err != nil {
		return
	}
	name = strings.TrimSpace(name)
	str = strings.TrimSpace(str)
	header.Add(name, str)
	return
}

func setHeadersMap(header *http.Header, v reflect.Value, tag reflect.StructTag) (err error) {
	prefix := tag.Get("locationName")
	for _, key := range v.MapKeys() {
		var str string
		str, err = convertType(v.MapIndex(key), tag)
		if errors.Is(err, errValueNotSet) {
			err = nil
			continue
		}
		if err != nil {
			return
		}
		keyStr := strings.TrimSpace(key.String())
		str = strings.TrimSpace(str)
		header.Add(prefix+keyStr, str)
	}
	return
}

func getPayload(v reflect.Value) (ptyp string, pfvl reflect.Value) {
	ptyp = noPayload

	field, ok := v.Type().FieldByName("_")
	if !ok {
		return
	}

	noPayloadValue := field.Tag.Get(noPayload)
	if noPayloadValue != "" {
		return
	}

	payloadName := field.Tag.Get("payload")
	if payloadName == "" {
		return
	}

	pfld, ok := v.Type().FieldByName(payloadName)
	if !ok {
		return
	}

	ptyp = pfld.Tag.Get("type")
	pfvl = reflect.Indirect(v.FieldByName(payloadName))

	return
}

func convertType(v reflect.Value, tag reflect.StructTag) (str string, err error) {
	v = reflect.Indirect(v)
	if !v.IsValid() {
		return
	}

	switch value := v.Interface().(type) {
	case string:
		str = value
	case []byte:
		str = base64.StdEncoding.EncodeToString(value)
	case bool:
		str = strconv.FormatBool(value)
	case int64:
		str = strconv.FormatInt(value, 10)
	case time.Time:
		str = protocol.FormatTime(getTimeFormat(tag), value)
	default:
		err = fmt.Errorf("unsupported value type <%s>", v.Type())
	}

	return
}

func getTimeFormat(tag reflect.StructTag) string {
	format := tag.Get("timestampFormat")
	if len(format) == 0 {
		format = protocol.RFC822TimeFormatName
		if tag.Get("location") == "querystring" {
			format = protocol.ISO8601TimeFormatName
		}
	}
	return format
}
