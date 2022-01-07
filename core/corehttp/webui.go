package corehttp

const WebUIPath = "/btfs/QmRZYABH4LisEQPQmpzHAAxacDZbsKtnysjzEuZBaug6bG" // v2.0.1

// this is a list of all past webUI paths.
var WebUIPaths = []string{
	WebUIPath,
	"/btfs/QmRwxtQfzpfaLKYfg3qhxsFfpRm3J3qLXJxDofsLV8ydXq", // v2.0.0
}

var HostUIOption = RedirectOption("hostui", WebUIPath)
var WebUIOption = RedirectOption("webui", WebUIPath)
var DashboardOption = RedirectOption("dashboard", WebUIPath)
