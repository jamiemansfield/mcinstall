// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ftb

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"git.sr.ht/~jmansfield/go-modpacksch/modpacksch"
	"github.com/gammazero/workerpool"
	"github.com/google/uuid"
	"github.com/jamiemansfield/mcinstall/forge"
	"github.com/jamiemansfield/mcinstall/minecraft"
	"github.com/jamiemansfield/mcinstall/minecraft/launcher"
)

const (
	defaultDataDir = ".ftbinstall"

	settingsFile = "install.json"
)

var (
	OtherPackAlreadyInstalled = errors.New("ftb: a pack is already installed at this location")
)

type Installer struct {
	// The directory to store mcinstall-related files (settings.json, etc)
	DataDir string

	// Directories that shouldn't be touched, in any capacity, when updating
	// modpacks.
	ExcludedDirs []string

	// The Minecraft Forge installer to use, should it be needed
	ForgeInstaller *forge.Installer

	workerPool *workerpool.WorkerPool
}

func NewInstaller(maxWorkers int) *Installer {
	return &Installer{
		DataDir: defaultDataDir,
		ExcludedDirs: []string{
			"saves",
		},
		ForgeInstaller: forge.NewInstaller(),
		workerPool:     workerpool.New(maxWorkers),
	}
}

// IsExcludedDir determines whether a directory should be excluded from being
// changed as the result of any update logic.
func (i *Installer) IsExcludedDir(relPath string) bool {
	parts := strings.Split(relPath, string(filepath.Separator))

	excludedDirs := append(i.ExcludedDirs, i.DataDir)
	for _, excludedDir := range excludedDirs {
		if parts[0] == excludedDir {
			return true
		}
	}

	return false
}

// Installs the given pack version to the destination, with the
// appropriate files for that install target.
func (i *Installer) InstallPackVersion(installTarget minecraft.InstallTarget, dest string, pack *modpacksch.Pack, version *modpacksch.PackVersion) error {
	fmt.Println("Installing " + pack.Name + " v" + version.Name + "...")

	destination, err := filepath.Abs(dest)
	if err != nil {
		return err
	}

	// Find existing install (or create one)
	if err := os.MkdirAll(i.DataDir, os.ModePerm); err != nil {
		return err
	}

	var settings *InstallSettings
	if readJson(filepath.Join(destination, i.DataDir, settingsFile), &settings) != nil {
		settings = &InstallSettings{
			ID:      uuid.New().String(),
			Pack:    pack.ID,
			Version: version.ID,
			Target:  installTarget,
			Files:   map[string]string{},
		}
	} else {
		fmt.Println("Existing installation of " + strconv.Itoa(settings.Pack) + " v" + strconv.Itoa(settings.Version) + " detected")

		if pack.ID != settings.Pack {
			return OtherPackAlreadyInstalled
		}
	}
	install := &Install{
		Version:       version.ID,
		OriginalFiles: settings.Files,
		NewFiles:      map[string]string{},
	}

	if err := i.InstallTargets(installTarget, destination, version.Targets); err != nil {
		return err
	}
	if err := i.InstallFiles(install, installTarget, destination, version.Files); err != nil {
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
		if i.IsExcludedDir(relPath) {
			return nil
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
			fmt.Printf("You can remove the '%s' line from %s/%s if\n", ftbPath, i.DataDir, settingsFile)
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
	if installTarget == minecraft.Client {
		// Get the target Minecraft version for the pack
		var mcVersion *minecraft.Version
		for _, target := range version.Targets {
			if target.Type == "game" {
				ver, err := minecraft.ParseVersion(target.Version)
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

		// Create profile
		profile := &launcher.Profile{
			Name:    pack.Name + " " + version.Name,
			Type:    "custom",
			GameDir: destination,
		}

		// Add icon to pack
		icon, err := launcher.CreateIconFromURL(pack.GetIcon().URL)
		if err != nil {
			fmt.Printf("Failed to get pack icon: %e", err)
		} else {
			profile.Icon = icon
		}

		// Set profile version, based on modloader in use
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

					profile.Version = forgeVersion
					break
				}

				// todo: other modloaders
			}
		}

		// Install profile
		if err := launcher.InstallProfile(settings.ID, profile); err != nil {
			return err
		}
	}

	// Write install settings
	settings.Version = install.Version
	settings.Files = install.NewFiles
	return writeJson(filepath.Join(destination, i.DataDir, settingsFile), &settings)
}

type Install struct {
	Version       int
	OriginalFiles map[string]string
	NewFiles      map[string]string
}

// ftbinstall.json
type InstallSettings struct {
	ID      string                  `json:"id"`
	Pack    int                     `json:"pack"`
	Version int                     `json:"version"`
	Target  minecraft.InstallTarget `json:"target"`
	Files   map[string]string       `json:"files"`
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
