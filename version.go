package btfs

// CurrentCommit is the current git commit, this is set as a ldflag in the Makefile
var CurrentCommit string

// CurrentVersionNumber is the current application's version literal
const CurrentVersionNumber = "2.3.2"

const ApiVersion = "/go-btfs/" + CurrentVersionNumber + "/"

// UserAgent is the libp2p user agent used by go-btfs.
//
// Note: This will end in `/` when no commit is available. This is expected.
var UserAgent = "go-btfs/" + CurrentVersionNumber + "/" + CurrentCommit
