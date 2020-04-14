// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ftbinstall

import (
	"bytes"
	"encoding/base64"
	"github.com/jamiemansfield/ftbinstall/mcinstall"
	"github.com/jamiemansfield/ftbinstall/util"
	"github.com/jamiemansfield/go-ftbmeta/ftbmeta"
	"net/http"
	"path/filepath"
)

// Installs the given pack version to the destination, with the
// appropriate files for that install target.
func InstallPackVersion(installTarget mcinstall.InstallTarget, dest string, pack *ftbmeta.Pack, version *ftbmeta.Version) error {
	destination, err := filepath.Abs(dest)
	if err != nil {
		return err
	}

	err = InstallTargets(installTarget, destination, version.Targets)
	if err != nil {
		return err
	}
	err = InstallFiles(installTarget, destination, version.Files)
	if err != nil {
		return err
	}

	// Install profile for the Minecraft launcher
	if installTarget == mcinstall.Client {
		// Get the target Minecraft version for the pack
		var mcVersion *mcinstall.McVersion
		for _, target := range version.Targets {
			if target.Type == "game" {
				ver, err := mcinstall.ParseMcVersion(target.Version)
				if err != nil {
					return err
				}
				mcVersion = ver
				break
			}
		}

		// If we can't determine the game version, we can't really proceed
		if mcVersion == nil {
			return FailedToDetermineGameVersion
		}

		// Get icon for pack profile
		req, err := util.NewRequest(http.MethodGet, pack.Art["square"].URL, nil)
		if err != nil {
			return err
		}
		writer := new(bytes.Buffer)
		if err := util.Download(writer, req); err != nil {
			return err
		}
		icon := "data:image/png;base64," + base64.StdEncoding.EncodeToString(writer.Bytes())

		// Create profile
		for _, target := range version.Targets {
			if target.Type == "modloader" {
				// Minecraft Forge
				if target.Name == "forge" {
					var version string
					// Minecraft 1.13 and above
					if mcVersion.Major >= 1 && mcVersion.Minor >= 13 {
						version = mcVersion.String() + "-forge-" + target.Version
					} else {
						version = mcVersion.String() + "-forge" + mcVersion.String() + "-" + target.Version
					}

					return mcinstall.InstallProfile(pack.Slug, &mcinstall.Profile{
						Name:    pack.Name,
						Type:    "custom",
						GameDir: destination,
						Icon:    icon,
						Version: version,
					})
				}

				// todo: other modloaders
			}
		}
	}

	return nil
}
