// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package technic

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jamiemansfield/ftbinstall/minecraft"
	"github.com/jamiemansfield/ftbinstall/minecraft/launcher"
	"github.com/jamiemansfield/ftbinstall/util"
	"github.com/jamiemansfield/go-technic/platform"
	"github.com/jamiemansfield/go-technic/solder"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Installs the given pack version to the destination, with the
// appropriate files for that install target.
func InstallPackVersion(dest string, pack *platform.Modpack, version string) error {
	fmt.Printf("Installing %s (%s)...\n", pack.DisplayName, pack.Name)

	destination, err := filepath.Abs(dest)
	if err != nil {
		return err
	}

	var mcVersion *minecraft.Version

	// Download the files
	if pack.Solder != "" {
		client := solder.NewClient(nil)
		solderUrl, err := url.Parse(pack.Solder)
		if err != nil {
			return err
		}
		client.BaseURL = solderUrl

		build, err := client.Modpack.GetBuild(pack.Name, version)
		if err != nil {
			return err
		}

		fmt.Printf("Installing from Solder...\n")

		total := len(build.Mods)
		for i, mod := range build.Mods {
			fmt.Printf("[%d / %d] Installing %s...\n", i + 1, total, mod.Name)

			err := downloadAndExtractZip(mod.URL, dest)
			if err != nil {
				return err
			}
		}

		mcVersion, err = minecraft.ParseVersion(build.Minecraft)
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

		mcVersion, err = minecraft.ParseVersion(pack.Minecraft)
		if err != nil {
			return err
		}
	}

	_, modpackJarExists := os.Stat(filepath.Join(dest,
		"bin", "modpack.jar",
	))
	_, versionJsonExists := os.Stat(filepath.Join(dest,
		"bin", "version.json",
	))
	versionName := mcVersion.String() + "-" + pack.Name + "-" + version
	_, launcherVersionJsonExists := os.Stat(filepath.Join(launcher.GetLauncherDir(),
		"versions", versionName, versionName + ".json",
	))
	_, launcherVersionJarExists := os.Stat(filepath.Join(launcher.GetLauncherDir(),
		"versions", versionName, versionName + ".jar",
	))

	// Create a version for the pack (mcversion-pack-version)
	versionDir := filepath.Join(launcher.GetLauncherDir(), "versions", versionName)
	if err := os.MkdirAll(versionDir, os.ModePerm); err != nil {
		return err
	}

	if launcherVersionJsonExists != nil || launcherVersionJarExists != nil {
		fmt.Printf("Installing version '%s'...\n", versionName)
	}

	if launcherVersionJsonExists != nil {
		if versionJsonExists == nil {
			// Open bin/version.json and launcher version
			modpackJsonFile, err := os.Open(filepath.Join(dest, "bin", "version.json"))
			if err != nil {
				return err
			}
			defer modpackJsonFile.Close()
			versionJsonFile, err := os.Create(filepath.Join(versionDir, versionName + ".json"))
			if err != nil {
				return err
			}
			defer versionJsonFile.Close()

			// Rewrite version.json, and save to launcher
			versionInfo, err := rewriteVersionJson(modpackJsonFile, versionName)
			encoder := json.NewEncoder(versionJsonFile)
			encoder.SetIndent("", "\t")
			err = encoder.Encode(versionInfo)
			if err != nil {
				return err
			}
		} else {
			// Create simple version for modpack
			launcherVersion := &launcher.Version{
				ID:           versionName,
				Type:         "release",
				InheritsFrom: mcVersion.String(),
			}

			// Use LegacyLaunch for pre-1.6 packs
			if mcVersion.Major < 1 || (mcVersion.Major == 1 && mcVersion.Minor < 6) {
				fmt.Println("Installing LegacyLaunch")

				legacyLaunch, mainClass, err := launcher.InstallLegacyLaunch(launcher.GetLauncherDir())
				if err != nil {
					return err
				}

				launcherVersion.MainClass = mainClass
				launcherVersion.Libraries = append(launcherVersion.Libraries, legacyLaunch)
			}

			versionJsonFile, err := os.Create(filepath.Join(versionDir, versionName + ".json"))
			if err != nil {
				return err
			}
			defer versionJsonFile.Close()

			// Save version.json to launcher
			encoder := json.NewEncoder(versionJsonFile)
			encoder.SetIndent("", "\t")
			err = encoder.Encode(launcherVersion)
			if err != nil {
				return err
			}
		}
	}

	if launcherVersionJarExists != nil && modpackJarExists == nil {
		// Ensure that the client.jar exists
		if err := launcher.InstallClientVersion(launcher.GetLauncherDir(), mcVersion.String()); err != nil {
			return err
		}

		// Open client.jar
		clientJarFile, err := os.Open(filepath.Join(launcher.GetLauncherDir(),
			"versions", mcVersion.String(), mcVersion.String() + ".jar",
		))
		if err != nil {
			return err
		}
		defer clientJarFile.Close()
		clientJarStat, err := clientJarFile.Stat()
		if err != nil {
			return err
		}

		// Open modpack.jar
		modpackJarFile, err := os.Open(filepath.Join(dest, "bin", "modpack.jar"))
		if err != nil {
			return err
		}
		defer modpackJarFile.Close()
		modpackJarStat, err := modpackJarFile.Stat()
		if err != nil {
			return err
		}

		// Open jars
		modpackJar, err := zip.NewReader(modpackJarFile, modpackJarStat.Size())
		if err != nil {
			return err
		}
		clientJar, err := zip.NewReader(clientJarFile, clientJarStat.Size())
		if err != nil {
			return err
		}

		// Create new jar
		versionJarFile, err := os.Create(filepath.Join(launcher.GetLauncherDir(),
			"versions", versionName, versionName + ".jar",
		))
		if err != nil {
			return err
		}
		zw := zip.NewWriter(versionJarFile)
		defer func() {
			zw.Close()
			versionJarFile.Close()
		}()

		var files []string
		files, err = util.MergeZips(zw, modpackJar, files, nil)
		if err != nil {
			return err
		}
		files, err = util.MergeZips(zw, clientJar, files, func(name string) bool {
			return strings.HasPrefix(name, "META-INF/")
		})
		if err != nil {
			return err
		}
	}

	// Create a profile for the Minecraft launcher
	profile := &launcher.Profile{
		Name:    pack.DisplayName + " " + version,
		Type:    "custom",
		GameDir: destination,
		Version: versionName,
	}

	// Attempt to add pack icon to pack
	if pack.Icon != nil {
		icon, err := launcher.CreateIconFromURL(pack.Icon.URL)
		if err != nil {
			fmt.Printf("Failed to get pack icon: %e", err)
		} else {
			profile.Icon = icon
		}
	}

	// Install the profile to the launcher
	return launcher.InstallProfile(pack.Name, profile)
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

func rewriteVersionJson(r io.Reader, id string) (map[string]interface{}, error) {
	var profile map[string]interface{}
	if err := json.NewDecoder(r).Decode(&profile); err != nil {
		return nil, err
	}

	// Make required changes
	profile["id"] = id

	return profile, nil
}
