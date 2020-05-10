// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

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
