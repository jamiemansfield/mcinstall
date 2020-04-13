// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package forgeinstall

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/jamiemansfield/ftbinstall/mcinstall"
	"github.com/jamiemansfield/ftbinstall/util"
	"io"
	"os"
	"path/filepath"
)

// See InstallForge
// Installs Minecraft Forge for Minecraft 1.5 -> 1.12
func installUniversalForge(target mcinstall.InstallTarget, dest string, mcVersion *mcinstall.McVersion, forgeVersion string) error {
	fmt.Println("Using universal Forge installer...")
	version := mcVersion.String() + "-" + forgeVersion

	// Check whether we need to install the server
	if _, err := os.Stat("forge-" + version + "-universal.jar"); err == nil && target == mcinstall.Server {
		fmt.Println("Minecraft Forge install found, skipping...")
		return nil
	}

	// Check whether we need to install the client
	// todo: implement

	// Download installer
	installerJar, err := downloadForgeInstaller(version)
	if err != nil {
		return err
	}
	defer os.Remove(installerJar.Name())

	installerInfo, err := installerJar.Stat()
	if err != nil {
		return err
	}

	if target == mcinstall.Client {
		versionName := mcVersion.String() + "-forge" + version

		// Open installer jar, so we can pull files
		reader, err := zip.NewReader(installerJar, installerInfo.Size())
		if err != nil {
			return err
		}

		// Create directories for install
		versionDir := filepath.Join(dest, "versions", versionName)
		if err := os.MkdirAll(versionDir, os.ModePerm); err != nil {
			return err
		}
		libraryDir := filepath.Join(dest, "libraries", "net", "minecraftforge", "forge", version)
		if err := os.MkdirAll(libraryDir, os.ModePerm); err != nil {
			return err
		}

		// Save version info to disk
		installProfile, err := util.GetFileInZip(reader, "install_profile.json")
		if err != nil {
			return err
		}
		profileReader, err := installProfile.Open()
		if err != nil {
			return err
		}
		defer profileReader.Close()
		versionInfo, err := getVersionInfo(profileReader)
		if err != nil {
			return err
		}
		infoFile, err := os.Create(filepath.Join(versionDir, versionName + ".json"))
		if err != nil {
			return err
		}
		defer infoFile.Close()
		encoder := json.NewEncoder(infoFile)
		encoder.SetIndent("", "\t")
		err = encoder.Encode(versionInfo)
		if err != nil {
			return err
		}

		// Save Forge universal jar to disk
		universalJar, err := util.GetFileInZip(reader, "forge-" + version + "-universal.jar")
		if err != nil {
			return err
		}
		return util.CopyZipFileToDisk(universalJar, filepath.Join(libraryDir, "forge-" + version + ".jar"))
	} else {
		return util.RunCommand("java", "-jar", installerJar.Name(), "--installServer", dest)
	}
}

// Extracts the version information from Forge's install profile.
func getVersionInfo(r io.Reader) (map[string]interface{}, error) {
	var profile map[string]interface{}
	if err := json.NewDecoder(r).Decode(&profile); err != nil {
		return nil, err
	}
	return profile["versionInfo"].(map[string]interface{}), nil
}
