// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package forgeinstall

import (
	"fmt"
	"github.com/jamiemansfield/ftbinstall/mcinstall"
	"os"
)

// See InstallForge
// Installs Minecraft Forge for Minecraft 1.5 -> 1.12
func installUniversalForge(target mcinstall.InstallTarget, dest string, mcVersion *mcinstall.McVersion, forgeVersion string) error {
	fmt.Println("Using universal Forge installer...")
	version := mcVersion.String() + "-" + forgeVersion

	// Check whether we need to install the server
	if _, err := os.Stat("forge-" + version + "-universal.jar"); err == nil {
		fmt.Println("Minecraft Forge install found, skipping...")
		return nil
	}

	// Download installer
	installFile, err := downloadForgeInstaller(version)
	if err != nil {
		return err
	}
	defer os.Remove(installFile.Name())

	// Run installer
	return runInstaller(target, dest, installFile)
}
