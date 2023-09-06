package s3utils

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"reflect"
	"testing"
)

type req struct {
	_                 struct{} `embed:"PutObjectInput"`
	s3.PutObjectInput `location:"embed"`
	Body              io.ReadCloser `type:"blob"`
}

func TestParseRequest(t *testing.T) {
	var r req
	v := reflect.ValueOf(r)
	p := v.Type()
	n := v.NumField()
	for i := 0; i < n; i++ {
		ft := p.Field(i)
		fmt.Println(ft.Name)
	}

}
