// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ftbinstall

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jamiemansfield/ftbinstall/mcinstall"
	"github.com/jamiemansfield/ftbinstall/util"
	"github.com/jamiemansfield/go-ftbmeta/ftbmeta"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	DataDir = ".ftbinstall"
	SettingsFile = "install.json"
)

var (
	ExcludedDirs = []string{
		DataDir,
		"saves",
	}
)

var (
	OtherPackAlreadyInstalled = errors.New("ftbinstall: a pack is already installed at this location")
)

// Installs the given pack version to the destination, with the
// appropriate files for that install target.
func InstallPackVersion(installTarget mcinstall.InstallTarget, dest string, pack *ftbmeta.Pack, version *ftbmeta.Version) error {
	destination, err := filepath.Abs(dest)
	if err != nil {
		return err
	}

	// Find existing install (or create one)
	if err := os.MkdirAll(DataDir, os.ModePerm); err != nil {
		return err
	}

	var settings *InstallSettings
	if readJson(filepath.Join(destination, DataDir, SettingsFile), &settings) != nil {
		settings = &InstallSettings{
			ID:      uuid.New().String(),
			Pack:    pack.Slug,
			Version: version.Slug,
			Target:  installTarget,
			Files:   map[string]string{},
		}
	} else {
		fmt.Println("Existing installation of " + settings.Pack + " v" + settings.Version + " detected")

		if pack.Slug != settings.Pack {
			return OtherPackAlreadyInstalled
		}
	}
	install := &Install{
		Version:       version.Slug,
		OriginalFiles: settings.Files,
		NewFiles:      map[string]string{},
	}

	err = InstallTargets(installTarget, destination, version.Targets)
	if err != nil {
		return err
	}
	err = InstallFiles(install, installTarget, destination, version.Files)
	if err != nil {
		return err
	}

	// Remove any unmodified files that are no longer apart the pack
	err = filepath.Walk(destination, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(destination, path)
		if err != nil {
			return err
		}
		ftbPath := "./" + filepath.ToSlash(relPath)

		// While this should never be an issue anyway - for peace of mind,
		// ftbinstall will not delete ANYTHING from a protected directory.
		parts := strings.Split(relPath, string(filepath.Separator))
		for _, excludedDir := range ExcludedDirs {
			if parts[0] == excludedDir {
				return nil
			}
		}

		// Ignore if its a current file
		if install.NewFiles[ftbPath] != "" {
			return nil
		}

		// Check if its been modified
		if install.OriginalFiles[ftbPath] != "" {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			hasher := sha1.New()
			if _, err := io.Copy(hasher, f); err != nil {
				return err
			}
			hash := hex.EncodeToString(hasher.Sum(nil))

			// Delete file
			if hash == install.OriginalFiles[ftbPath] {
				fmt.Printf("%s has been removed from the modpack, as its\n", ftbPath)
				fmt.Println("sha1 hash matches the original, it has been removed.")
				return os.Remove(path)
			}

			// The file has been removed from the pack, but the player has modified it
			fmt.Printf("%s has been removed from the modpack, as its\n", ftbPath)
			fmt.Println("sha1 hash doesn't match the original - we have left it in place.")
			fmt.Println("Please investigate whether you still need the file before playing!")
			fmt.Printf("You can remove the '%s' line from %s/%s if\n", ftbPath, DataDir, SettingsFile)
			fmt.Println("still required")

			// So that this message continues on, store the original hash in the new file list
			install.NewFiles[ftbPath] = install.OriginalFiles[ftbPath]
		}

		return nil
	})
	if err != nil {
		return err
	}

	// Install profile for the Minecraft launcher
	if installTarget == mcinstall.Client {
		// Get the target Minecraft version for the pack
		var mcVersion *mcinstall.McVersion
		for _, target := range version.Targets {
			if target.Type == "game" {
				ver, err := mcinstall.ParseMcVersion(target.Version)
				if err != nil {
					return err
				}
				mcVersion = ver
				break
			}
		}

		// If we can't determine the game version, we can't really proceed
		if mcVersion == nil {
			return FailedToDetermineGameVersion
		}

		// Get icon for pack profile
		req, err := util.NewRequest(http.MethodGet, pack.Art["square"].URL, nil)
		if err != nil {
			return err
		}
		writer := new(bytes.Buffer)
		if err := util.Download(writer, req); err != nil {
			return err
		}
		icon := "data:image/png;base64," + base64.StdEncoding.EncodeToString(writer.Bytes())

		// Create profile
		for _, target := range version.Targets {
			if target.Type == "modloader" {
				// Minecraft Forge
				if target.Name == "forge" {
					var forgeVersion string
					// Minecraft 1.13 and above
					if mcVersion.Major >= 1 && mcVersion.Minor >= 13 {
						forgeVersion = mcVersion.String() + "-forge-" + target.Version
					} else {
						forgeVersion = mcVersion.String() + "-forge" + mcVersion.String() + "-" + target.Version
					}

					if err := mcinstall.InstallProfile(settings.ID, &mcinstall.Profile{
						Name:    pack.Name + " " + version.Name,
						Type:    "custom",
						GameDir: destination,
						Icon:    icon,
						Version: forgeVersion,
					}); err != nil {
						return err
					}
					break
				}

				// todo: other modloaders
			}
		}
	}

	// Write install settings
	settings.Version = install.Version
	settings.Files = install.NewFiles
	return writeJson(filepath.Join(destination, DataDir, SettingsFile), &settings)
}

type Install struct {
	Version string
	OriginalFiles map[string]string
	NewFiles map[string]string
}

// ftbinstall.json
type InstallSettings struct {
	ID string `json:"id"`
	Pack string `json:"pack"`
	Version string `json:"version"`
	Target mcinstall.InstallTarget `json:"target"`
	Files map[string]string `json:"files"`
}

func readJson(destination string, v interface{}) error {
	contents, err := ioutil.ReadFile(destination)
	if err != nil {
		return err
	}
	return json.Unmarshal(contents, v)
}

func writeJson(destination string, v interface{}) error {
	out, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(destination, out, 0644)
}
