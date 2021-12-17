package corehttp

const WebUIPath = "/btns/16Uiu2HAmU9ysnuasmdyq1rRePYTwHntmyhZdfC9wm4qCPQMAh9Qq"

// this is a list of all past webUI paths.
var WebUIPaths = []string{
	WebUIPath,
}

var HostUIOption = RedirectOption("hostui", WebUIPath)
var WebUIOption = RedirectOption("webui", WebUIPath)
var DashboardOption = RedirectOption("dashboard", WebUIPath)
