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
	"github.com/jamiemansfield/ftbinstall/util"
	"github.com/jamiemansfield/go-technic/platform"
	"github.com/jamiemansfield/go-technic/solder"
	"io"
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

	var mcVersion *minecraft.Version

	// Download the files
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

		mcVersion, err = minecraft.ParseVersion(solderVersion.Minecraft)
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
	_, libraryExists := os.Stat(filepath.Join(minecraft.GetLauncherDir(),
		"libraries", "net", "technicpack", "pack", pack.Name, version, pack.Name + "-" + version + ".jar",
	))
	versionName := mcVersion.String() + "-" + pack.Name + "-" + version
	_, versionExists := os.Stat(filepath.Join(minecraft.GetLauncherDir(),
		"versions", versionName, versionName + ".json",
	))

	// Create a library for the pack (pack-version), if bin/modpack.jar exists
	if modpackJarExists == nil {
		if libraryExists == nil {
			fmt.Printf("Library already installed\n")
		} else {
			fmt.Printf("Installing library...\n")

			libraryDir := filepath.Join(minecraft.GetLauncherDir(),
				"libraries", "net", "technicpack", "pack", pack.Name, version,
			)
			if err := os.MkdirAll(libraryDir, os.ModePerm); err != nil {
				return err
			}

			// Open bin/modpack.jar and launcher library
			modpackJarFile, err := os.Open(filepath.Join(dest, "bin", "modpack.jar"))
			if err != nil {
				return err
			}
			defer modpackJarFile.Close()
			versionJarFile, err := os.Create(filepath.Join(libraryDir, pack.Name + "-" + version + ".jar"))
			if err != nil {
				return err
			}
			defer versionJarFile.Close()

			_, err = io.Copy(versionJarFile, modpackJarFile)
			if err != nil {
				return err
			}
		}
	}

	// Create a version for the pack (mcversion-pack-version)
	if versionExists != nil {
		fmt.Printf("Installing version '%s'...\n", versionName)

		versionDir := filepath.Join(minecraft.GetLauncherDir(), "versions", versionName)
		if err := os.MkdirAll(versionDir, os.ModePerm); err != nil {
			return err
		}

		if versionJsonExists == nil {
			// Open bin/modpack.json and launcher version
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
			launcherVersion := &minecraft.LauncherVersion{
				ID:           versionName,
				Type:         "release",
				InheritsFrom: mcVersion.String(),
				Libraries: []*minecraft.LauncherVersionLibrary{
					{
						Name: "net.technicpack.pack:" + pack.Name + ":" + version,
					},
				},
			}
			versionJsonFile, err := os.Create(filepath.Join(versionDir, versionName + ".json"))
			if err != nil {
				return err
			}
			defer versionJsonFile.Close()

			// Rewrite version.json, and save to launcher
			encoder := json.NewEncoder(versionJsonFile)
			encoder.SetIndent("", "\t")
			err = encoder.Encode(launcherVersion)
			if err != nil {
				return err
			}
		}
	}

	// Create a profile for the Minecraft launcher
	return minecraft.InstallProfile(pack.Name, &minecraft.Profile{
		Name:    pack.Name + " " + version,
		Type:    "custom",
		GameDir: destination,
		Version: versionName,
	})
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

func rewriteVersionJson(r io.Reader, id string) (map[string]interface{}, error) {
	var profile map[string]interface{}
	if err := json.NewDecoder(r).Decode(&profile); err != nil {
		return nil, err
	}

	// Make required changes
	profile["id"] = id

	return profile, nil
}
