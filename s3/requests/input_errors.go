package requests

import (
	"fmt"
	"reflect"
)

// ErrInvalidInputValue .
type ErrInvalidInputValue struct {
	er error
}

func (err ErrInvalidInputValue) Error() string {
	return fmt.Sprintf("invalid input value: %v", err.er)
}

// ErrTypeNotSet .
type ErrTypeNotSet struct {
	typ reflect.Type
}

func (err ErrTypeNotSet) Error() string {
	return fmt.Sprintf("type <%s> not set", err.typ.String())
}

// ErrFailedDecodeXML .
type ErrFailedDecodeXML struct {
	err error
}

func (err ErrFailedDecodeXML) Error() string {
	return fmt.Sprintf("decode xml: %v", err.err)
}

// ErrWithUnsupportedParam .
type ErrWithUnsupportedParam struct {
	param string
}

func (err ErrWithUnsupportedParam) Error() string {
	return fmt.Sprintf("param %s is unsported", err.param)
}

// ErrFailedParseValue .
type ErrFailedParseValue struct {
	name string
	err  error
}

func (err ErrFailedParseValue) Name() string {
	return err.name
}

func (err ErrFailedParseValue) Error() string {
	return fmt.Sprintf("parse <%s> value: %v", err.name, err.err)
}

// ErrMissingRequiredParam .
type ErrMissingRequiredParam struct {
	param string
}

func (err ErrMissingRequiredParam) Error() string {
	return fmt.Sprintf("missing required param <%s>", err.param)
}
