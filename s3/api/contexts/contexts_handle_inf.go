package contexts

import (
	"net/http"
)

type handleInfo struct {
	name string
	err  error
}

func SetHandleInf(r *http.Request, name string, err error) {
	set(r, keyOfHandleInf, handleInfo{name, err})
	return
}

func GetHandleInf(r *http.Request) (name string, err error) {
	v := get(r, keyOfHandleInf)
	inf, _ := v.(handleInfo)
	name = inf.name
	err = inf.err
	return
}
