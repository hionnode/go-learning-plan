package main

import (
	"fmt"
	"runtime"
)

// These are overwritten at build time via ldflags:
//   -ldflags "-X main.Version=0.1.0 -X main.Commit=abc1234 -X main.BuildDate=2026-04-18T..."
// setup.sh does this; a bare `go build` leaves the dev defaults.
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

func runVersion(ctx *appContext, args []string) error {
	_ = ctx
	_ = args
	fmt.Print(banner)
	fmt.Printf("  go-dojo     %s\n", Version)
	fmt.Printf("  commit      %s\n", Commit)
	fmt.Printf("  built       %s\n", BuildDate)
	fmt.Printf("  go runtime  %s\n", runtime.Version())
	fmt.Printf("  platform    %s/%s\n", runtime.GOOS, runtime.GOARCH)
	return nil
}
