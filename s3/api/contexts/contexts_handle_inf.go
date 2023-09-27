package contexts

import (
	"net/http"
)

type handleInfo struct {
	name string
	args interface{}
	err  error
}

func SetHandleInf(r *http.Request, name string, args interface{}, err error) {
	set(r, keyOfHandleInf, handleInfo{name, args, err})
	return
}

func GetHandleInf(r *http.Request) (name string, args interface{}, err error) {
	v := get(r, keyOfHandleInf)
	inf, _ := v.(handleInfo)
	name = inf.name
	err = inf.err
	args = inf.args
	return
}
