// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package forgeinstall

import (
	"fmt"
	"github.com/jamiemansfield/ftbinstall/mcinstall"
	"github.com/jamiemansfield/ftbinstall/util"
	"net/http"
	"os"
	"path/filepath"
)

const (
	// The URL to my client install tool
	ClientInstallTool = "https://repo.neptunepowered.org/tools/forge-client-installer-0.0.2.jar"
)

// See InstallForge
// Installs Minecraft Forge for Minecraft >= 1.13
func installModernForge(target mcinstall.InstallTarget, dest string, mcVersion *mcinstall.McVersion, forgeVersion string) error {
	fmt.Println("Using modern Forge installer...")
	version := mcVersion.String() + "-" + forgeVersion

	// Check whether we need to install Minecraft Forge
	_, serverCheck := os.Stat(filepath.Join(dest,
		"forge-" + version + ".jar",
	))
	_, clientCheck := os.Stat(filepath.Join(dest,
		"libraries", "net", "minecraftforge", "forge", version, "forge-" + version + ".jar",
	))
	if (serverCheck == nil && target == mcinstall.Server) ||
		(clientCheck == nil && target == mcinstall.Client) {
		fmt.Println("Minecraft Forge install found, skipping...")
		return nil
	}

	// Download installer
	installerJar, err := downloadForgeInstaller(version)
	if err != nil {
		return err
	}
	defer func() {
		installerJar.Close()
		os.Remove(installerJar.Name())
	}()

	// Create the appropriate arguments for the install target
	var args []string
	if target == mcinstall.Client {
		toolFile, err := downloadClientInstallTool()
		if err != nil {
			return err
		}
		defer func() {
			toolFile.Close()
			os.Remove(toolFile.Name())
		}()

		args = append(args, "-cp", toolFile.Name() + ";" + installerJar.Name(), "me.jamiemansfield.forgeclientinstaller.ClientInstaller", dest)
	} else {
		args = append(args, "-jar", installerJar.Name(), "--installServer", dest)
	}

	// Run installer
	return util.RunCommand("java", args...)
}

// Downloads the Forge Client Installer tool.
// The temporary file should be removed after usage.
func downloadClientInstallTool() (*os.File, error) {
	req, err := util.NewRequest(http.MethodGet, ClientInstallTool, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/java,application/java-archive,application/x-java-archive")
	return util.DownloadTemp(req, "installtool*.jar")
}
