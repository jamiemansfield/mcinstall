// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package manifest

import (
	"encoding/json"
	"net/http"

	"github.com/jamiemansfield/ftbinstall/util"
)

func GetVersionManifest(httpClient *http.Client) (*VersionManifest, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	// Create the request
	req, err := util.NewRequest(http.MethodGet, "https://launchermeta.mojang.com/mc/game/version_manifest.json", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	// Make the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Get the version manifest
	var response VersionManifest
	err = json.NewDecoder(resp.Body).Decode(&response)
	return &response, err
}

func (v *VersionManifestVersion) GetFull(httpClient *http.Client) (*Version, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	// Create the request
	req, err := util.NewRequest(http.MethodGet, v.URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	// Make the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Get the version manifest
	var response Version
	err = json.NewDecoder(resp.Body).Decode(&response)
	return &response, err
}
