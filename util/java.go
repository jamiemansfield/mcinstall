package util

import (
	"runtime"
	"strings"
)

// Gets the classpath separator for the operating system currently
// running.
func GetClasspathSeparator() string {
	if runtime.GOOS == "windows" {
		return ";"
	} else {
		return ":"
	}
}

// Creates a classpath string for the given files appropriate to the
// operating system currently running.
// For example, on macOS using BuildClasspath("x.jar", "y.jar") the
// output would be "x.jar:y.jar". However on Windows, the output
// would be "x.jar;y.jar".
func BuildClasspath(files ...string) string {
	return strings.Join(files, GetClasspathSeparator())
}
