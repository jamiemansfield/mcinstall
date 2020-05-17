// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package manifest

import "time"

type VersionManifest struct {
	Latest struct {
		Release string `json:"release"`
		Snapshot string `json:"snapshot"`
	} `json:"latest"`
	Versions []*VersionManifestVersion `json:"versions"`
}

func (m *VersionManifest) FindVersion(id string) *VersionManifestVersion {
	for _, version := range m.Versions {
		if version.ID == id {
			return version
		}
	}

	return nil
}

type VersionManifestVersion struct {
	ID string `json:"id"`
	Type string `json:"type"`
	URL string `json:"url"`
	Time time.Time `json:"time"`
	ReleaseTime time.Time `json:"releaseTime"`
}
