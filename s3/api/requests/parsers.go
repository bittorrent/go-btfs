package requests

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"github.com/aws/aws-sdk-go/private/protocol"
	"github.com/aws/aws-sdk-go/private/protocol/xml/xmlutil"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type fields map[string]bool

func ParseLocation(r *http.Request, input interface{}, supports fields) (err error) {
	inv, err := valueOf(input)
	if err != nil {
		return
	}
	err = parseLocation(r, inv, supports)
	return
}

func ParseXMLBody(r *http.Request, input interface{}) (err error) {
	inv, err := valueOf(input)
	if err != nil {
		return
	}
	pft, ok := getPayloadField(inv)
	if !ok {
		err = ErrPayloadNotSet{"field"}
		return
	}
	ptyp := pft.Tag.Get("type")
	if ptyp != "structure" {
		err = ErrPayloadNotSet{"structure"}
		return
	}
	decoder := xml.NewDecoder(r.Body)
	err = xmlutil.UnmarshalXML(inv.Addr().Interface(), decoder, "")
	if err != nil {
		err = ErrFailedDecodeXML{err}
	}
	return
}

func valueOf(input interface{}) (inv reflect.Value, err error) {
	inv = reflect.Indirect(reflect.ValueOf(input))
	if !inv.IsValid() {
		err = ErrInvalidInputValue{"input is nil"}
		return
	}
	if inv.Kind() != reflect.Struct {
		err = ErrInvalidInputValue{"input is not point to struct"}
	}
	return
}

func parseLocation(r *http.Request, inv reflect.Value, supports fields) (err error) {
	vars := mux.Vars(r)
	headers := r.Header
	query := r.URL.Query()
	for i := 0; i < inv.NumField(); i++ {
		fv := inv.Field(i)
		ft := inv.Type().Field(i)
		err = parseLocationField(vars, query, headers, fv, ft, supports)
		if err != nil {
			return
		}
	}
	return
}

func parseLocationField(vars map[string]string, query url.Values, headers http.Header, fv reflect.Value,
	ft reflect.StructField, supports fields) (err error) {
	if ft.Name[0:1] == strings.ToLower(ft.Name[0:1]) {
		return
	}
	ftag := ft.Tag
	loca := ftag.Get("location")
	name := ftag.Get("locationName")
	requ := ftag.Get("required") == "true"
	supp := supports[ft.Name]
	var (
		vals   map[string]*string
		isVals bool
		val    string
		has    bool
	)
	switch loca {
	case "querystring":
		val, has = query.Get(name), query.Has(name)
	case "uri":
		val, has = vars[name]
	case "header":
		val, has = headers.Get(name), len(headers.Values(name)) > 0
	case "headers":
		vals, has = getHeaderValues(headers, name)
		isVals = true
	default:
		return
	}
	if !supp && has {
		err = ErrWithUnsupportedParam{name}
		return
	}
	if requ && !has {
		err = ErrMissingRequiredParam{name}
		return
	}
	if !has {
		return
	}
	if isVals {
		err = parseValues(vals, fv)
	} else {
		err = parseValue(val, fv, ftag)
		if err != nil && !errors.As(err, new(ErrTypeNotSet)) {
			err = ErrFailedParseValue{name, err}
		}
	}
	return
}

func getPayloadField(inv reflect.Value) (ft reflect.StructField, ok bool) {
	mt, ok := inv.Type().FieldByName("_")
	if !ok {
		return
	}
	if mt.Tag.Get("nopayload") != "" {
		return
	}
	pname := mt.Tag.Get("payload")
	if pname == "" {
		return
	}
	ft, ok = inv.Type().FieldByName(pname)
	return
}

func getHeaderValues(header http.Header, prefix string) (vals map[string]*string, has bool) {
	defer func() {
		has = len(vals) > 0
	}()
	vals = make(map[string]*string)
	if len(header) == 0 {
		return
	}
	for key := range header {
		if len(key) >= len(prefix) && strings.EqualFold(key[:len(prefix)], prefix) {
			val := header.Get(key)
			k := strings.ToLower(key[len(prefix):])
			vals[k] = &val
		}
	}
	return
}

func parseValues(values map[string]*string, fv reflect.Value) (err error) {
	_, ok := fv.Interface().(map[string]*string)
	if !ok {
		err = ErrTypeNotSet{fv.Type()}
		return
	}
	fv.Set(reflect.ValueOf(values))
	return
}

func parseValue(value string, rv reflect.Value, tag reflect.StructTag) (err error) {
	switch rv.Interface().(type) {
	case *string:
		rv.Set(reflect.ValueOf(&value))
		return
	case []*string:
		var val []*string
		val, err = split(value)
		if err != nil {
			return
		}
		rv.Set(reflect.ValueOf(&val))
		return
	case []byte:
		var val []byte
		val, err = base64.StdEncoding.DecodeString(value)
		if err != nil {
			return
		}
		rv.Set(reflect.ValueOf(val))
		return
	case *bool:
		var val bool
		val, err = strconv.ParseBool(value)
		if err != nil {
			return
		}
		rv.Set(reflect.ValueOf(&val))
		return
	case *int64:
		var val int64
		val, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			return
		}
		rv.Set(reflect.ValueOf(&val))
		return
	case *time.Time:
		var val time.Time
		format := getTimeFormat(tag)
		val, err = protocol.ParseTime(format, value)
		if err != nil {
			return
		}
		rv.Set(reflect.ValueOf(&val))
		return
	default:
		err = ErrTypeNotSet{rv.Type()}
		return
	}
}

func getTimeFormat(tag reflect.StructTag) (format string) {
	format = tag.Get("timestampFormat")
	if format != "" {
		return
	}
	if tag.Get("location") == "querystring" {
		format = protocol.ISO8601TimeFormatName
		return
	}
	format = protocol.RFC822TimeFormatName
	return
}

func split(value string) (vals []*string, err error) {
	pv := ' '
	start := 0
	quote := false
	for i, v := range value {
		opv := pv
		pv = v
		if quote {
			if v == '"' && opv != '\\' {
				quote = false
				val := value[start : i+1]
				val, err = strconv.Unquote(val)
				if err != nil {
					return
				}
				val = strings.TrimSpace(val)
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
			val := value[start:i]
			val = strings.TrimSpace(val)
			vals = append(vals, &val)
			start = i + 1
		}
		continue
	}
	if quote {
		err = errors.New("unquote part")
		return
	}
	if start < len(value) || pv == ',' {
		val := value[start:]
		val = strings.TrimSpace(val)
		vals = append(vals, &val)
	}
	return
}
