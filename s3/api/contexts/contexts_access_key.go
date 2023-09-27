package contexts

import (
	"net/http"
)

func SetAccessKey(r *http.Request, ack string) {
	set(r, keyOfAccessKey, ack)
	return
}

func GetAccessKey(r *http.Request) (ack string) {
	v := get(r, keyOfAccessKey)
	ack, _ = v.(string)
	return
}
