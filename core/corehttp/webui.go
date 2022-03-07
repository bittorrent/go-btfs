package corehttp

const WebUIPath = "/btfs/QmWDZ94ZMAjts3WSPbFdLUbfLMYbygJR7BNEygVJqxuqfw" // v2.1.1

// this is a list of all past webUI paths.
var WebUIPaths = []string{
	WebUIPath,
	"/btfs/QmVKKZuhYriR26jAkZGwR7jEPLGARtWZZV3oS2abykwT2U", // v2.1.0
	"/btfs/QmPSaMbVPTrcPg8CxW8GRrKYY9YZyxEQek9eKg5is13S9H", // v2.0.2
	"/btfs/QmRZYABH4LisEQPQmpzHAAxacDZbsKtnysjzEuZBaug6bG", // v2.0.1
	"/btfs/QmRwxtQfzpfaLKYfg3qhxsFfpRm3J3qLXJxDofsLV8ydXq", // v2.0.0
}

var HostUIOption = RedirectOption("hostui", WebUIPath)
var WebUIOption = RedirectOption("webui", WebUIPath)
var DashboardOption = RedirectOption("dashboard", WebUIPath)
