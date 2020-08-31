// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ftb

import (
	"errors"

	"github.com/jamiemansfield/go-modpacksch/modpacksch"
	"github.com/jamiemansfield/mcinstall/forge"
	"github.com/jamiemansfield/mcinstall/minecraft"
	"github.com/jamiemansfield/mcinstall/minecraft/launcher"
)

var (
	FailedToDetermineGameVersion = errors.New("failed to determine game version")
)

// Installs the given targets, for the target environment, to the given
// destination.
func InstallTargets(installTarget minecraft.InstallTarget, dest string, targets []*modpacksch.Target) error {
	// Get the target Minecraft version for the pack
	var mcVersion *minecraft.Version
	for _, target := range targets {
		if target.Type == "game" {
			ver, err := minecraft.ParseVersion(target.Version)
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
		} else
		if target.Type == "modloader" {
			var loaderDest string
			if installTarget == minecraft.Client {
				loaderDest = launcher.GetLauncherDir()
			} else {
				loaderDest = dest
			}

			// Minecraft Forge
			if target.Name == "forge" {
				if err := forge.InstallForge(installTarget, loaderDest, mcVersion, target.Version); err != nil {
					return err
				}
			}

			// todo: liteloader, fabric, etc support?
		}
	}

	return nil
}
