// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ftbinstall

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/jamiemansfield/ftbinstall/mcinstall"
	"github.com/jamiemansfield/ftbinstall/util"
	"github.com/jamiemansfield/go-ftbmeta/ftbmeta"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Installs the given files, for the target environment, to the given
// destination.
func InstallFiles(install *InstallSettings, target mcinstall.InstallTarget, dest string, files []*ftbmeta.File) error {
	// Collect the target-specific files, so we can keep an accurate count
	// of how many files we've installed.
	var targetFiles []*ftbmeta.File
	for _, file := range files {
		// Ignore files for another target
		if (target == mcinstall.Client && file.ServerOnly) || (target == mcinstall.Server && file.ClientOnly) {
			continue
		}
		targetFiles = append(targetFiles, file)
	}

	// Install files for the target
	for i, file := range targetFiles {
		msg, err := installFile(dest, file)
		if err != nil {
			fmt.Printf("[%d / %d] Failed to install '%s%s', ignoring file...`n", i + 1, len(targetFiles), file.Path, file.Name)
			fmt.Println(err)
			continue
		} else if msg != "" {
			fmt.Printf("[%d / %d] %s\n", i + 1, len(targetFiles), msg)
		}

		// Log the files information in the install settings
		install.Files[file.Path + file.Name] = file.Sha1
	}

	return nil
}

// Installs the given file, to the destination
func installFile(dest string, file *ftbmeta.File) (string, error) {
	dirPath := filepath.Join(dest, filepath.FromSlash(file.Path))
	fileDest := filepath.Join(dirPath, file.Name)

	// If file already exists, check the checksum
	if _, err := os.Stat(fileDest); err == nil {
		f, err := os.Open(fileDest)
		if err != nil {
			return "", err
		}
		defer f.Close()

		hasher := sha1.New()
		if _, err := io.Copy(hasher, f); err != nil {
			return "", err
		}

		// If already exists, continue to next file
		if hex.EncodeToString(hasher.Sum(nil)) == file.Sha1 {
			return fmt.Sprintf("%s%s found, skipping...", file.Path, file.Name), nil
		}
	}

	// Ensure directory exists
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return "", err
	}

	// GET the file
	req, err := util.NewRequest(http.MethodGet, file.URL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "*")

	// Write file to disk
	f, err := os.Create(fileDest)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return fmt.Sprintf("Installed '%s' to '%s'", file.Name, file.Path), util.Download(f, req)
}
