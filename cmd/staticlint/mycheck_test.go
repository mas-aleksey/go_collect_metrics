package main

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func TestMyAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), ExitCheckAnalyzer, "./main")
}

func Example() {
	// `go build cmd/staticlint/mycheck.go`
	// `./mycheck ./...``
}
