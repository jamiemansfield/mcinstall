// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package launcher

import (
	"errors"
	"fmt"
	"github.com/jamiemansfield/ftbinstall/minecraft/manifest"
	"github.com/jamiemansfield/ftbinstall/util"
	"net/http"
	"os"
	"path/filepath"
)

// Struct for creating simple version JSONs.
type Version struct {
	ID string                   `json:"id"`
	Type string                 `json:"type"`
	InheritsFrom string         `json:"inheritsFrom"`
	MainClass string            `json:"mainClass,omitempty"`
	Libraries []*VersionLibrary `json:"libraries,omitempty"`
	Downloads struct {
	} `json:"downloads"`
}

type VersionLibrary struct {
	Name string `json:"name"`
	URL string `json:"url,omitempty"`
}

var (
	ErrVersionDoesntExist = errors.New("launcher: given version doesn't exist")
)

// InstallClientVersion installs the given client version, to the given
// launcher directory.
func InstallClientVersion(launcherDir string, versionName string) error {
	versionDir := filepath.Join(launcherDir, "versions", versionName)
	versionJar := filepath.Join(versionDir, versionName + ".jar")
	versionJson := filepath.Join(versionDir, versionName + ".json")

	// Check whether we need to install the client version
	_, versionJarExists := os.Stat(versionJar)
	_, versionJsonExists := os.Stat(versionJson)
	if versionJarExists == nil && versionJsonExists == nil {
		return nil
	}

	// Get version information
	versions, err := manifest.GetVersionManifest(nil)
	if err != nil {
		return err
	}
	versionInfo := versions.FindVersion(versionName)
	if versionInfo == nil {
		return ErrVersionDoesntExist
	}

	// Ensure version directory exists
	if err := os.MkdirAll(versionDir, os.ModePerm); err != nil {
		return err
	}

	// Download version.jar
	if versionJarExists != nil {
		fmt.Println("Downloading " + versionName + " client jar...")

		// Get full version
		version, err := versionInfo.GetFull(nil)
		if err != nil {
			return err
		}

		// Create file
		f, err := os.Create(versionJar)
		if err != nil {
			return err
		}
		defer f.Close()

		// Create request
		req, err := util.NewRequest(http.MethodGet, version.Downloads.Client.URL, nil)
		if err != nil {
			return err
		}
		req.Header.Add("Accepts", "*")

		// Download file
		if err := util.Download(f, req); err != nil {
			return err
		}
	}

	// Download version.json
	if versionJsonExists != nil {
		fmt.Println("Downloading " + versionName + " json...")

		// Create file
		f, err := os.Create(versionJson)
		if err != nil {
			return err
		}
		defer f.Close()

		// Create request
		req, err := util.NewRequest(http.MethodGet, versionInfo.URL, nil)
		if err != nil {
			return err
		}
		req.Header.Add("Accepts", "application/json")

		// Download file
		if err := util.Download(f, req); err != nil {
			return err
		}
	}

	return nil
}
