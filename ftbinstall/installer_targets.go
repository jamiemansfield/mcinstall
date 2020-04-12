// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ftbinstall

import (
	"errors"
	"github.com/jamiemansfield/ftbinstall/forgeinstall"
	"github.com/jamiemansfield/ftbinstall/mcinstall"
	"github.com/jamiemansfield/go-ftbmeta/ftbmeta"
)

var (
	FailedToDetermineGameVersion = errors.New("failed to determine game version")
)

// Installs the given targets, for the target environment, to the given
// destination.
func InstallTargets(installTarget mcinstall.InstallTarget, dest string, targets []*ftbmeta.Target) error {
	// Get the target Minecraft version for the pack
	var mcVersion *mcinstall.McVersion
	for _, target := range targets {
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

	// Install mod loaders, etc
	for _, target := range targets {
		if target.Type == "game" {
			continue
		}

		if target.Type == "modloader" {
			// todo: use launcher root as dest for Client

			// Minecraft Forge
			if target.Name == "forge" {
				if err := forgeinstall.InstallForge(installTarget, dest, mcVersion, target.Version); err != nil {
					return err
				}
			}

			// todo: liteloader, fabric, etc support?
		}
	}

	return nil
}
