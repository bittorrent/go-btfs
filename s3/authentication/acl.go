package authentication

type ACLKey string

const (
	ACLKeyPrivate         ACLKey = "private"
	ACLKeyPublicRead      ACLKey = "public-read"
	ACLKeyPublicReadWrite ACLKey = "public-read-write"
)
