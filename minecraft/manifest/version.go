// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package manifest

import "time"

type Version struct {
	ID string `json:"id"`
	Type string `json:"type"`
	Time time.Time `json:"time"`
	ReleaseTime time.Time `json:"releaseTime"`
	MinimumLauncherVersion int `json:"minimumLauncherVersion"`
	MainClass string `json:"mainClass"`
	Downloads struct {
		Client *VersionDownload `json:"client"`
		ClientMappings *VersionDownload `json:"client_mappings"`
		Server *VersionDownload `json:"server"`
		ServerMappings *VersionDownload `json:"server_mappings"`
	} `json:"downloads"`
}

type VersionDownload struct {
	Sha1 string `json:"sha1"`
	Size int `json:"size"`
	URL string `json:"url"`
}
