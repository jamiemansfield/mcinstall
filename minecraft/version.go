// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package minecraft

import (
	"errors"
	"strconv"
	"strings"
)

type McVersion struct {
	Major int
	Minor int
	Revision int
}

func ParseMcVersion(version string) (*McVersion, error) {
	parts := strings.Split(version, ".")

	if len(parts) < 2 {
		return nil, errors.New("invalid Minecraft version: '" + version + "'")
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}
	var revision int
	if len(parts) > 2 {
		revision, err = strconv.Atoi(parts[2])
		if err != nil {
			return nil, err
		}
	} else {
		revision = 0
	}

	return &McVersion{
		Major:    major,
		Minor:    minor,
		Revision: revision,
	}, nil
}

func (v *McVersion) String() string {
	if v.Revision == 0 {
		return strconv.Itoa(v.Major) + "." + strconv.Itoa(v.Minor)
	}

	return strconv.Itoa(v.Major) + "." + strconv.Itoa(v.Minor) + "." + strconv.Itoa(v.Revision)
}
