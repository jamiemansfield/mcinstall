// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ftb

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jamiemansfield/go-ftbmeta/ftbmeta"
	"github.com/jamiemansfield/mcinstall/minecraft"
	"github.com/jamiemansfield/mcinstall/util"
)

// Installs the given files, for the target environment, to the given
// destination.
func InstallFiles(install *Install, target minecraft.InstallTarget, dest string, files []*ftbmeta.File) error {
	// Collect the target-specific files, so we can keep an accurate count
	// of how many files we've installed.
	var targetFiles []*ftbmeta.File
	for _, file := range files {
		// Ignore files for another target
		if (target == minecraft.Client && file.ServerOnly) || (target == minecraft.Server && file.ClientOnly) {
			continue
		}
		targetFiles = append(targetFiles, file)
	}

	// Install files for the target
	for i, file := range targetFiles {
		msg, err := installFile(install, dest, file)
		if err != nil {
			fmt.Printf("[%d / %d] Failed to install '%s%s', ignoring file...\n", i+1, len(targetFiles), file.Path, file.Name)
			fmt.Println(err)
			continue
		}
		fmt.Printf("[%d / %d] %s\n", i+1, len(targetFiles), msg)

		// Log the files information in the install settings
		install.NewFiles[file.Path+file.Name] = file.Sha1
	}

	return nil
}

// Installs the given file, to the destination
func installFile(install *Install, dest string, file *ftbmeta.File) (string, error) {
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
		hash := hex.EncodeToString(hasher.Sum(nil))

		// If already exists, continue to next file
		if hash == file.Sha1 {
			return fmt.Sprintf("%s%s found, skipping...", file.Path, file.Name), nil
		}

		// If the file previously existed, don't override if the player made changes
		originalHash := install.OriginalFiles[file.Path+file.Name]
		if originalHash != "" && hash != originalHash {
			fmt.Println("************************************************************************************************")
			fmt.Printf("%s%s has a sha1 has of '%s', when\n", file.Path, file.Name, hash)
			fmt.Printf("'%s' was expected.\n", originalHash)
			fmt.Printf("To prevent overriding configurations, it will be installed under %s/%s.\n", DataDir, install.Version)
			fmt.Println("Please investigate any collisions before playing!")
			fmt.Println("************************************************************************************************")

			return installFile(install, filepath.Join(dest, DataDir, install.Version), file)
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
