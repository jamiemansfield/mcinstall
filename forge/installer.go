// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package forge

import (
	"net/http"
	"os"

	"github.com/jamiemansfield/mcinstall/minecraft"
	"github.com/jamiemansfield/mcinstall/util"
)

const (
	// The URL to Minecraft Forge's Maven, or a mirror.
	MavenRoot = "https://files.minecraftforge.net/maven/"
)

// Installs Minecraft Forge to the given destination, for the given target.
// If the target is Server, the destination will be the root directory of
// the server; if the target is Client, the destination will be the
// launcher's root directory.
func InstallForge(target minecraft.InstallTarget, dest string, mcVersion *minecraft.Version, forgeVersion string) error {
	// Use modern installer - Minecraft 1.13 and above
	if mcVersion.Major >= 1 && mcVersion.Minor >= 13 {
		return installModernForge(target, dest, mcVersion, forgeVersion)
	} else
	// Use universal install method - Minecraft 1.5 -> Minecraft 1.12
	if mcVersion.Major >= 1 && mcVersion.Minor >= 5 && mcVersion.Minor <= 12 {
		return installUniversalForge(target, dest, mcVersion, forgeVersion)
	}
	// todo: older

	return nil
}

// Downloads the Minecraft Forge installer for the given version (MC-Forge),
// to a temporary file.
// The temporary file should be removed after usage.
func downloadForgeInstaller(version string) (*os.File, error) {
	url := MavenRoot + "net/minecraftforge/forge/" + version + "/forge-" + version + "-installer.jar"

	// Download installer
	req, err := util.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/java,application/java-archive,application/x-java-archive")

	return util.DownloadTemp(req, "forge*.jar")
}
