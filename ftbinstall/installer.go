// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ftbinstall

import (
	"github.com/jamiemansfield/ftbinstall/mcinstall"
	"github.com/jamiemansfield/go-ftbmeta/ftbmeta"
)

// Installs the given pack version to the destination, with the
// appropriate files for that install target.
func InstallPackVersion(target mcinstall.InstallTarget, dest string, version *ftbmeta.Version) error {
	err := InstallTargets(target, dest, version.Targets)
	if err != nil {
		return err
	}
	err = InstallFiles(target, dest, version.Files)
	if err != nil {
		return err
	}

	return nil
}
