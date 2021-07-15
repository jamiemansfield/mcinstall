// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package forge

import (
	"strconv"
	"strings"
)

// TODO: understand Forge versions better
type Version struct {
	Build int
}

func ParseVersion(version string) (*Version, error) {
	parts := strings.Split(version, ".")

	build, err := strconv.Atoi(parts[len(parts) - 1])
	if err != nil {
		return nil, err
	}

	return &Version{
		Build: build,
	}, nil
}
