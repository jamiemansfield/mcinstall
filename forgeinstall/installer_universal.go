// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package forgeinstall

import (
	"errors"
	"fmt"
	"github.com/jamiemansfield/ftbinstall/mcinstall"
	"github.com/jamiemansfield/ftbinstall/util"
	"os"
	"path/filepath"
)

// See InstallForge
// Installs Minecraft Forge for Minecraft 1.5 -> 1.12
func installUniversalForge(target mcinstall.InstallTarget, dest string, mcVersion *mcinstall.McVersion, forgeVersion string) error {
	fmt.Println("Using universal Forge installer...")
	version := mcVersion.String() + "-" + forgeVersion

	destination, err := filepath.Abs(dest)
	if err != nil {
		return err
	}

	// Check whether we need to install the server
	if _, err := os.Stat("forge-" + version + "-universal.jar"); err == nil && target == mcinstall.Server {
		fmt.Println("Minecraft Forge install found, skipping...")
		return nil
	}

	// Check whether we need to install the client
	// todo: implement

	// Download installer
	installerJar, err := downloadForgeInstaller(version)
	if err != nil {
		return err
	}
	defer os.Remove(installerJar.Name())

	if target == mcinstall.Client {
		// todo: implement support
		return errors.New("client support not yet implemented")
	} else {
		return util.RunCommand("java", "-jar", installerJar.Name(), "--installServer", destination)
	}
}
