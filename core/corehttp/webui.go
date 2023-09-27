package corehttp

const WebUIPath = "/btfs/QmPFT7PscyJ1FZ4FeFPFgikszBugQSBFNycPpy5zpK2pZe" // v2.3.3

// this is a list of all past webUI paths.
var WebUIPaths = []string{
	WebUIPath,
	"/btfs/QmUKCyDc4h9KN93AdZ7ZVqgPProsKs8NAbVJkK3ux9788d", // v2.3.2
	"/btfs/QmRdt8SzRBz5px7KfU4hFveJSKzBMFqv73YE4xXJBsVdDJ", // v2.3.1
	"/btfs/QmbNHqcL9PEhFdT5mXjNnkaAE8SEFkjr2jD7we2ckTL4Li", // v2.3.0
	"/btfs/QmZvpBNMribwdjNMrA9gXz27t2gzbae3N2tbCLtjpRTqJn", // v2.2.1.1
	"/btfs/QmaK77EYUHxKweLFvRY8gbcMTx2qEb7p4S5aWPN6EHX7T1", // v2.2.1
	"/btfs/QmXsvmvTTzciHEdbDCCMo55MfrEQ6qnct8B4Wt9aJwHoMY", // v2.2.0
	"/btfs/QmZM3CcoWPHiu9E8ugS76a2csooKEAp5YQitgjQC849h4b", // v2.1.3
	"/btfs/QmW3VGCuvfhAJcJZRYQeEjJnjQG27kHNhBasrF2TwGniTT", // v2.1.2
	"/btfs/QmWDZ94ZMAjts3WSPbFdLUbfLMYbygJR7BNEygVJqxuqfw", // v2.1.1
	"/btfs/QmVKKZuhYriR26jAkZGwR7jEPLGARtWZZV3oS2abykwT2U", // v2.1.0
	"/btfs/QmPSaMbVPTrcPg8CxW8GRrKYY9YZyxEQek9eKg5is13S9H", // v2.0.2
	"/btfs/QmRZYABH4LisEQPQmpzHAAxacDZbsKtnysjzEuZBaug6bG", // v2.0.1
	"/btfs/QmRwxtQfzpfaLKYfg3qhxsFfpRm3J3qLXJxDofsLV8ydXq", // v2.0.0
}

var HostUIOption = RedirectOption("hostui", WebUIPath)
var WebUIOption = RedirectOption("webui", WebUIPath)
var DashboardOption = RedirectOption("dashboard", WebUIPath)
