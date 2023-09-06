package responses

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/private/protocol"
	"github.com/aws/aws-sdk-go/private/protocol/xml/xmlutil"
	"github.com/gorilla/mux"
	"io"
	"math"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var byteSliceType = reflect.TypeOf([]byte{})

func ParseRequest(r *http.Request, input interface{}) (err error) {
	inv, err := getInputValue(input)
	if err != nil {
		return
	}

	err = parseLocation(r, inv)
	if err != nil {
		return
	}

	ptyp, pftp, pfvl := getPayload(inv)
	if ptyp == noPayload {
		return
	}

	if ptyp == "structure" || ptyp == "" {
		err = parseXMLBody(r, inv)
	} else {
		err = parseBody(r, pftp, pfvl)
	}

	return

}

func parseXMLBody(r *http.Request, inv reflect.Value) (err error) {
	defer r.Body.Close()
	decoder := xml.NewDecoder(r.Body)
	err = xmlutil.UnmarshalXML(inv.Addr().Interface(), decoder, "")
	return
}

func parseBody(r *http.Request, pftp reflect.Type, pfvl reflect.Value) (err error) {
	var b []byte
	switch pfvl.Interface().(type) {
	case []byte:
		defer r.Body.Close()
		b, err = io.ReadAll(r.Body)
		if err != nil {
			return
		}
		pfvl.Set(reflect.ValueOf(b))
	case *string:
		defer r.Body.Close()
		b, err = io.ReadAll(r.Body)
		if err != nil {
			return
		}
		val := string(b)
		pfvl.Set(reflect.ValueOf(&val))
	default:
		switch pftp.String() {
		case "io.ReadSeeker":
			// keep the request body
		default:
			err = errValueNotSet
		}
	}
	return
}

func getInputValue(input interface{}) (inv reflect.Value, err error) {
	typErr := fmt.Errorf("input <%T> must be non nil <T *struct> or <T **struct>", input)

	if input == nil {
		err = typErr
		return
	}

	t := reflect.TypeOf(input)
	k := t.Kind()

	if k != reflect.Pointer {
		err = typErr
		return
	}

	inv = reflect.ValueOf(input).Elem()
	if !inv.IsValid() {
		err = typErr
		return
	}

	t = t.Elem()
	k = t.Kind()

	if k == reflect.Struct {
		return
	}

	if k != reflect.Pointer {
		err = typErr
		return
	}

	t = t.Elem()
	k = t.Kind()
	if k != reflect.Struct {
		err = typErr
		return
	}

	if inv.Elem().IsValid() {
		inv = inv.Elem()
		return
	}

	inv.Set(reflect.New(inv.Type().Elem()))
	inv = inv.Elem()

	return
}

func parseLocation(r *http.Request, inv reflect.Value) (err error) {
	query := r.URL.Query()

	for i := 0; i < inv.NumField(); i++ {
		fv := inv.Field(i)
		ft := inv.Type().Field(i)
		if ft.Name[0:1] == strings.ToLower(ft.Name[0:1]) {
			continue
		}

		if ft.Tag.Get("ignore") != "" {
			continue
		}

		name := ifemp(ft.Tag.Get("locationName"), ft.Name)

		if ft.Tag.Get("marshal-as") == "blob" {
			if fv.Kind() == reflect.Pointer {
				fv.Set(reflect.New(fv.Type().Elem()))
				fv = fv.Elem()
			}
			fv = fv.Convert(byteSliceType)
		}

		switch ft.Tag.Get("location") {
		case "headers":
			prefix := ft.Tag.Get("locationName")
			err = parseHeaderMap(r.Header, fv, prefix)
		case "header":
			locVal := r.Header.Get(name)
			err = parseLocationValue(locVal, fv, ft.Tag)
		case "uri":
			locVal := mux.Vars(r)[name]
			err = parseLocationValue(locVal, fv, ft.Tag)
		case "querystring":
			err = parseQueryString(query, fv, name, ft.Tag)
		}
	}

	return
}

func parseQueryString(query url.Values, fv reflect.Value, name string, tag reflect.StructTag) (err error) {
	switch value := fv.Interface().(type) {
	case []*string:
		vals := make([]*string, len(query[name]))
		for i, oval := range query[name] {
			val := oval
			vals[i] = &val
		}
		if len(vals) > 0 {
			fv.Set(reflect.ValueOf(vals))
		}
	case map[string]*string:
		vals := make(map[string]*string, len(query))
		for key := range query {
			val := query.Get(key)
			vals[key] = &val
		}
		if len(vals) > 0 {
			fv.Set(reflect.ValueOf(vals))
		}
	case map[string][]*string:
		for key, items := range value {
			for _, item := range items {
				query.Add(key, *item)
			}
		}
		vals := make(map[string][]*string, len(query))
		for key := range query {
			vals[key] = make([]*string, len(query[key]))
			for i := range query[key] {
				vals[key][i] = &(query[key][i])
			}
		}
		if len(vals) > 0 {
			fv.Set(reflect.ValueOf(vals))
		}
	default:
		locVal := query.Get(name)
		err = parseLocationValue(locVal, fv, tag)
		if err != nil {
			return
		}
	}

	return
}

func parseHeaderMap(headers http.Header, fv reflect.Value, prefix string) (err error) {
	if len(headers) == 0 {
		return
	}
	switch fv.Interface().(type) {
	case map[string]*string:
		vals := map[string]*string{}
		for key := range headers {
			if !hasPrefixFold(key, prefix) {
				continue
			}
			key = strings.ToLower(key)
			val := headers.Get(key)
			vals[key[len(prefix):]] = &val
		}
		if len(vals) != 0 {
			fv.Set(reflect.ValueOf(vals))
		}
	default:
		err = errValueNotSet
	}
	return
}

func parseLocationValue(locVal string, v reflect.Value, tag reflect.StructTag) (err error) {
	switch tag.Get("type") {
	case "jsonvalue":
		if len(locVal) == 0 {
			return
		}
	case "blob":
		if len(locVal) == 0 {
			return
		}
	default:
		if !v.IsValid() || (locVal == "" && (v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.String)) {
			return
		}
	}

	switch v.Interface().(type) {
	case *string:
		if tag.Get("suppressedJSONValue") == "true" && tag.Get("location") == "header" {
			var b []byte
			b, err = base64.StdEncoding.DecodeString(locVal)
			if err != nil {
				return
			}
			locVal = string(b)
		}
		v.Set(reflect.ValueOf(&locVal))
	case []*string:
		if tag.Get("location") != "header" || tag.Get("enum") == "" {
			return fmt.Errorf("%T is only supported with location header and enum shapes", v)
		}
		var vals []*string
		vals, err = splitHeaderVal(locVal)
		if err != nil {
			return
		}
		if len(vals) > 0 {
			v.Set(reflect.ValueOf(vals))
		}
	case []byte:
		var b []byte
		b, err = base64.StdEncoding.DecodeString(locVal)
		if err != nil {
			return
		}
		v.Set(reflect.ValueOf(b))
	case *bool:
		var b bool
		b, err = strconv.ParseBool(locVal)
		if err != nil {
			return
		}
		v.Set(reflect.ValueOf(&b))
	case *int64:
		var i int64
		i, err = strconv.ParseInt(locVal, 10, 64)
		if err != nil {
			return
		}
		v.Set(reflect.ValueOf(&i))
	case *float64:
		var f float64
		switch {
		case strings.EqualFold(locVal, floatNaN):
			f = math.NaN()
		case strings.EqualFold(locVal, floatInf):
			f = math.Inf(1)
		case strings.EqualFold(locVal, floatNegInf):
			f = math.Inf(-1)
		default:
			f, err = strconv.ParseFloat(locVal, 64)
			if err != nil {
				return
			}
		}
		v.Set(reflect.ValueOf(&f))
	case *time.Time:
		format := tag.Get("timestampFormat")
		if len(format) == 0 {
			format = protocol.RFC822TimeFormatName
			if tag.Get("location") == "querystring" {
				format = protocol.ISO8601TimeFormatName
			}
		}
		var t time.Time
		t, err = protocol.ParseTime(format, locVal)
		if err != nil {
			return
		}
		v.Set(reflect.ValueOf(&t))
	case aws.JSONValue:
		escaping := protocol.NoEscape
		if tag.Get("location") == "header" {
			escaping = protocol.Base64Escape
		}
		var m aws.JSONValue
		m, err = protocol.DecodeJSONValue(locVal, escaping)
		if err != nil {
			return
		}
		v.Set(reflect.ValueOf(m))
	default:
		err = fmt.Errorf("unsupported value for input %v (%s)", v.Interface(), v.Type())
		return
	}

	return
}

func hasPrefixFold(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[0:len(prefix)], prefix)
}

func splitHeaderVal(header string) (vals []*string, err error) {
	pv := ' '
	start := 0
	quote := false
	for i, v := range header {
		opv := pv
		pv = v
		if quote {
			if v == '"' && opv != '\\' {
				quote = false
				val := header[start : i+1]
				val, err = strconv.Unquote(val)
				if err != nil {
					return
				}
				vals = append(vals, &val)
				start = i + 1
			}
			continue
		}

		if v == '"' && opv != '\\' {
			quote = true
			continue
		}

		if v == ',' && opv == '"' {
			start += 1
			continue
		}

		if v == ',' {
			val := header[start:i]
			vals = append(vals, &val)
			start = i + 1
		}

		continue
	}

	if quote {
		err = errors.New("unquote part")
		return
	}

	if start < len(header) || pv == ',' {
		val := header[start:]
		vals = append(vals, &val)
	}

	return
}
