// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mcinstall

import (
	"os"
	"path/filepath"
	"runtime"
)

// Gets the root directory of the Minecraft Launcher for the system
// mcinstall is running on.
func GetLauncherDir() string {
	if runtime.GOOS == "windows" {
		if appdata, present := os.LookupEnv("APPDATA"); present {
			return filepath.Join(appdata, ".minecraft")
		}
	} else
	if runtime.GOOS == "darwin" {
		userHome, _ := os.UserHomeDir()
		return filepath.Join(userHome, "Library", "Application Support", "minecraft")
	}

	userHome, _ := os.UserHomeDir()
	return filepath.Join(userHome, ".minecraft")
}
