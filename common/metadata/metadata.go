package metadata

import "fmt"

// package-scoped variables

// Package version
var Version string

// package-scoped constants
const LogModule = "application"

func GetVersionInfo() string {
	if Version == "" {
		Version = "development build"
	}

	return fmt.Sprintf("Version: %s", Version)
}
