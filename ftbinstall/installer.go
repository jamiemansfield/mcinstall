// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ftbinstall

import (
	"github.com/jamiemansfield/ftbinstall/mcinstall"
	"github.com/jamiemansfield/go-ftbmeta/ftbmeta"
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
		var mcVersion string
		for _, target := range version.Targets {
			if target.Type == "game" {
				mcVersion = target.Version
				break
			}
		}

		for _, target := range version.Targets {
			if target.Type == "modloader" {
				// Minecraft Forge
				if target.Name == "forge" {
					versionId := mcVersion + "-forge" + mcVersion + "-" + target.Version

					err = mcinstall.InstallProfile(pack.Slug, &mcinstall.Profile{
						Name:    pack.Name,
						Type:    "custom",
						GameDir: destination,
						Icon:    "Grass", // todo:
						Version: versionId,
					})
					if err != nil {
						return err
					}
				}

				// todo: other modloaders
			}
		}
	}

	return nil
}
