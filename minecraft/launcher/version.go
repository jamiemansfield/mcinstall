// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package launcher

import "time"

// Struct for creating simple version JSONs.
type Version struct {
	ID string                   `json:"id"`
	Type string                 `json:"type"`
	InheritsFrom string         `json:"inheritsFrom"`
	Time time.Time              `json:"time,omitempty"`
	ReleaseTime time.Time       `json:"releaseTime,omitempty"`
	MainClass string            `json:"mainClass,omitempty"`
	Libraries []*VersionLibrary `json:"libraries,omitempty"`
}

type VersionLibrary struct {
	Name string `json:"name"`
	URL string `json:"url,omitempty"`
}
