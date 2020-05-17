// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package technic

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/jamiemansfield/ftbinstall/util"
	"github.com/jamiemansfield/go-technic/platform"
	"github.com/jamiemansfield/go-technic/solder"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// Installs the given pack version to the destination, with the
// appropriate files for that install target.
func InstallPackVersion(dest string, pack *platform.Modpack, version string) error {
	fmt.Printf("Installing %s (%s)...\n", pack.DisplayName, pack.Name)

	destination, err := filepath.Abs(dest)
	if err != nil {
		return err
	}

	if pack.Solder != "" {
		client := solder.NewClient(nil)
		solderUrl, err := url.Parse(pack.Solder)
		if err != nil {
			return err
		}
		client.BaseURL = solderUrl

		solderPack, err := client.Modpack.GetModpack(pack.Name)
		if err != nil {
			return err
		}
		solderVersion, err := client.Modpack.GetBuild(pack.Name, version)
		if err != nil {
			return err
		}

		err = InstallSolderBuild(destination, solderPack, solderVersion)
		if err != nil {
			return err
		}
	} else {
		if version != pack.Version {
			return errors.New("versions do not match, pack is at " + pack.Version)
		}
		fmt.Printf("Installing from direct download...\n")

		err = downloadAndExtractZip(pack.URL, dest)
		if err != nil {
			return err
		}
	}

	return nil
}

func InstallSolderBuild(dest string, pack *solder.Modpack, build *solder.Build) error {
	fmt.Printf("Installing from Solder...\n")

	// Download the zips to a temporary directory
	tmp, err := ioutil.TempDir(dest, "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	total := len(build.Mods)
	for i, mod := range build.Mods {
		fmt.Printf("[%d / %d] Installing %s...\n", i + 1, total, mod.Name)

		err = downloadAndExtractZip(mod.URL, dest)
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadAndExtractZip(url string, dest string) error {
	// GET the file
	req, err := util.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "*")

	// Download to temporary file
	tmp, err := util.DownloadTemp(req, "*.zip")
	if err != nil {
		return err
	}
	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()
	tmpInfo, err := tmp.Stat()
	if err != nil {
		return err
	}

	// Extract zip file
	zipFile, err := zip.NewReader(tmp, tmpInfo.Size())
	if err != nil {
		return err
	}
	return util.ExtractZipFileToDisk(zipFile, dest)
}
