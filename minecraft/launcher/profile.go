// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package launcher

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jamiemansfield/mcinstall/util"
)

// Struct for quickly creating new profiles.
type Profile struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	GameDir string `json:"gameDir"`
	Icon    string `json:"icon"`
	Version string `json:"lastVersionId"`
}

// Installs the given profile to the Minecraft launcher.
func InstallProfile(id string, profile *Profile) error {
	path := filepath.Join(GetLauncherDir(), "launcher_profiles.json")

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	raw["profiles"].(map[string]interface{})[id] = profile

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	return encoder.Encode(&raw)
}

// CreateIconFromURL creates a string that can be used within a Profile
// as a profile icon, from a remote resource.
func CreateIconFromURL(url string) (string, error) {
	req, err := util.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	writer := new(bytes.Buffer)
	if err := util.Download(writer, req); err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(writer.Bytes()), nil
}
