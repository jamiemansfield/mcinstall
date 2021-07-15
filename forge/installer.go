// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package forge

import (
	"net/http"
	"net/url"
	"os"

	"github.com/jamiemansfield/mcinstall/minecraft"
	"github.com/jamiemansfield/mcinstall/util"
)

const (
	defaultMavenRoot = "https://files.minecraftforge.net/maven/"
)

type Installer struct {
	// The URL to Minecraft Forge's Maven, or a mirror.
	MavenRoot *url.URL
}

// NewInstaller returns a new Installer to use for installing Minecraft
// Forge.
func NewInstaller() *Installer {
	mavenRoot, _ := url.Parse(defaultMavenRoot)

	return &Installer{
		MavenRoot: mavenRoot,
	}
}

// Installs Minecraft Forge to the given destination, for the given target.
// If the target is Server, the destination will be the root directory of
// the server; if the target is Client, the destination will be the
// launcher's root directory.
func (i *Installer) InstallForge(target minecraft.InstallTarget, dest string, mcVersion *minecraft.Version, forgeVersion string) error {
	forgeVrsn, err := ParseVersion(forgeVersion)
	if err != nil {
		return err
	}

	// Use modern installer - Minecraft 1.13 and above / newer Minecraft 1.12 builds
	if (mcVersion.Major >= 1 && mcVersion.Minor >= 13) || (mcVersion.Minor == 12 && forgeVrsn.Build >= 2851) {
		return i.installModernForge(target, dest, mcVersion, forgeVersion)
	} else
	// Use universal install method - Minecraft 1.5 -> Minecraft 1.12
	if mcVersion.Major >= 1 && mcVersion.Minor >= 5 && mcVersion.Minor <= 12 {
		return i.installUniversalForge(target, dest, mcVersion, forgeVersion)
	}
	// todo: older

	return nil
}

// Downloads the Minecraft Forge installer for the given version (MC-Forge),
// to a temporary file.
// The temporary file should be removed after usage.
func (i *Installer) downloadForgeInstaller(version string) (*os.File, error) {
	u, err := i.MavenRoot.Parse("net/minecraftforge/forge/" + version + "/forge-" + version + "-installer.jar")
	if err != nil {
		return nil, err
	}

	// Download installer
	req, err := util.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/java,application/java-archive,application/x-java-archive")

	return util.DownloadTemp(req, "forge*.jar")
}
