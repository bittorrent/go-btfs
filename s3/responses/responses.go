package responses

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/private/protocol/xml/xmlutil"
	"github.com/bittorrent/go-btfs/s3/consts"
	"io"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/private/protocol"
)

const (
	mimeTypeXml = "application/xml"
	noPayload   = "nopayload"
)

const (
	floatNaN    = "NaN"
	floatInf    = "Infinity"
	floatNegInf = "-Infinity"
)

var errValueNotSet = fmt.Errorf("value not set")

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

func wrapOutput(v interface{}, locationName string) (wrapper interface{}) {
	if locationName == "" {
		wrapper = v
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
			Type: reflect.TypeOf(v),
			Tag:  reflect.StructTag(outputTag),
		},
	}
	wrapperTyp := reflect.StructOf(fields)
	wrapperVal := reflect.New(wrapperTyp)
	wrapperVal.Elem().Field(1).Set(reflect.ValueOf(v))
	wrapper = wrapperVal.Interface()
	return
}

func extractBody(v reflect.Value) (body io.ReadCloser, clen int, ctyp string, err error) {
	ptyp, plod := getPayload(v)
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

	if plod.Interface() == nil {
		return
	}

	switch pifc := plod.Interface().(type) {
	case io.ReadCloser:
		body = pifc
		clen = -1
	case io.ReadSeeker:
		var bs []byte
		bs, err = io.ReadAll(pifc)
		if err != nil {
			return
		}
		body = io.NopCloser(bytes.NewBuffer(bs))
		clen = len(bs)
		ctyp = http.DetectContentType(bs)
	case []byte:
		body = io.NopCloser(bytes.NewBuffer(pifc))
		clen = len(pifc)
	case string:
		body = io.NopCloser(bytes.NewBufferString(pifc))
		clen = len(pifc)
	default:
		err = fmt.Errorf(
			"unknown payload type %s",
			plod.Type(),
		)
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
	reqId := getRequestID()
	headers.Set(consts.ServerInfo, consts.DefaultServerInfo)
	headers.Set(consts.AcceptRanges, "bytes")
	headers.Set(consts.AmzRequestID, reqId)
}

func setLocationHeaders(header http.Header, v reflect.Value) (err error) {
	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		ft := v.Type().Field(i)
		fk := fv.Kind()

		if n := ft.Name; n[0:1] == strings.ToLower(n[0:1]) {
			continue
		}

		if !fv.IsValid() {
			continue
		}

		if fk == reflect.Ptr {
			fv = fv.Elem()
			fk = fv.Kind()
			if !fv.IsValid() {
				continue
			}
		} else if fk == reflect.Interface {
			if !fv.Elem().IsValid() {
				continue
			}
		}

		switch ft.Tag.Get("location") {
		case "header":
			name := ifemp(ft.Tag.Get("locationName"), ft.Name)
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
	if errors.Is(err, errValueNotSet) {
		err = nil
		return
	}
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

func getPayload(v reflect.Value) (ptyp string, plod reflect.Value) {
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

	member, ok := v.Type().FieldByName(payloadName)
	if !ok {
		return
	}

	ptyp = member.Tag.Get("type")
	plod = reflect.Indirect(v.FieldByName(payloadName))

	return
}

func convertType(v reflect.Value, tag reflect.StructTag) (str string, err error) {
	v = reflect.Indirect(v)
	if !v.IsValid() {
		err = errValueNotSet
		return
	}

	switch value := v.Interface().(type) {
	case string:
		if tag.Get("suppressedJSONValue") == "true" && tag.Get("location") == "header" {
			value = base64.StdEncoding.EncodeToString([]byte(value))
		}
		str = value
	case []byte:
		str = base64.StdEncoding.EncodeToString(value)
	case bool:
		str = strconv.FormatBool(value)
	case int64:
		str = strconv.FormatInt(value, 10)
	case float64:
		switch {
		case math.IsNaN(value):
			str = floatNaN
		case math.IsInf(value, 1):
			str = floatInf
		case math.IsInf(value, -1):
			str = floatNegInf
		default:
			str = strconv.FormatFloat(value, 'f', -1, 64)
		}
	case time.Time:
		format := tag.Get("timestampFormat")
		if len(format) == 0 {
			format = protocol.RFC822TimeFormatName
			if tag.Get("location") == "querystring" {
				format = protocol.ISO8601TimeFormatName
			}
		}
		str = protocol.FormatTime(format, value)
	case aws.JSONValue:
		if len(value) == 0 {
			return "", errValueNotSet
		}
		escaping := protocol.NoEscape
		if tag.Get("location") == "header" {
			escaping = protocol.Base64Escape
		}
		str, err = protocol.EncodeJSONValue(value, escaping)
	default:
		err = fmt.Errorf("unsupported value for param %v (%s)", v.Interface(), v.Type())
	}

	return
}

func ifemp(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
