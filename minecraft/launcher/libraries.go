// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package launcher

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
)

//go:generate go run github.com/wlbr/mule -o legacylaunch.mule.go -p launcher legacylaunch/build/legacylaunch-1.0.0.jar

// Installs LegacyLaunch to the given launcher directory
func InstallLegacyLaunch(launcherDir string) (*VersionLibrary, string, error) {
	libraryDir := filepath.Join(launcherDir, "libraries", "me", "jamiemansfield", "mcinstall", "legacylaunch", "1.0.0")
	libraryJar := filepath.Join(libraryDir, "legacylaunch-1.0.0.jar")

	// Install LegacyLaunch, if it doesn't exist
	if _, check := os.Stat(libraryJar); check != nil {
		// Ensure the directory exists
		if err := os.MkdirAll(libraryDir, os.ModePerm); err != nil {
			return nil, "", err
		}

		client, err := legacylaunchResource()
		if err != nil {
			return nil, "", err
		}

		f, err := os.Create(libraryJar)
		if err != nil {
			return nil, "", err
		}
		defer f.Close()

		if _, err := io.Copy(f, bytes.NewReader(client)); err != nil {
			return nil, "", err
		}
	}

	return &VersionLibrary{
		Name: "me.jamiemansfield.mcinstall:legacylaunch:1.0.0",
	}, "me.jamiemansfield.mcinstall.LegacyLauncher", nil
}
