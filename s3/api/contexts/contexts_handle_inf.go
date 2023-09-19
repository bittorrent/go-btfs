package contexts

import (
	"net/http"
)

type handleInfo struct {
	name string
	err  error
	args interface{}
}

func SetHandleInf(r *http.Request, name string, err error, args interface{}) {
	set(r, keyOfHandleInf, handleInfo{name, err, args})
	return
}

func GetHandleInf(r *http.Request) (name string, err error, args interface{}) {
	v := get(r, keyOfHandleInf)
	inf, _ := v.(handleInfo)
	name = inf.name
	err = inf.err
	args = inf.args
	return
}
