// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package forgeinstall

//go:generate go run github.com/wlbr/mule -o modernclient.mule.go -p forgeinstall tool/build/forgetool.jar

import (
	"bytes"
	"fmt"
	"github.com/jamiemansfield/ftbinstall/mcinstall"
	"github.com/jamiemansfield/ftbinstall/util"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// See InstallForge
// Installs Minecraft Forge for Minecraft >= 1.13
func installModernForge(target mcinstall.InstallTarget, dest string, mcVersion *mcinstall.McVersion, forgeVersion string) error {
	version := mcVersion.String() + "-" + forgeVersion

	// Check whether we need to install Minecraft Forge
	_, serverCheck := os.Stat(filepath.Join(dest,
		"forge-" + version + ".jar",
	))
	_, clientCheck := os.Stat(filepath.Join(dest,
		"libraries", "net", "minecraftforge", "forge", version, "forge-" + version + ".jar",
	))
	if (serverCheck == nil && target == mcinstall.Server) ||
		(clientCheck == nil && target == mcinstall.Client) {
		fmt.Println("Minecraft Forge install found, skipping...")
		return nil
	}
	fmt.Printf("Installing Minecraft Forge %s-%s using modern installer...\n", mcVersion, forgeVersion)

	// Download installer
	installerJar, err := downloadForgeInstaller(version)
	if err != nil {
		return err
	}
	defer func() {
		installerJar.Close()
		os.Remove(installerJar.Name())
	}()

	// Create the appropriate arguments for the install target
	var args []string
	if target == mcinstall.Client {
		toolJar, err := copyClientInstallTool()
		if err != nil {
			return err
		}
		defer os.Remove(toolJar)

		args = append(args, "-cp", util.BuildClasspath(toolJar, installerJar.Name()), "ModernForgeClientTool", dest)
	} else {
		args = append(args, "-jar", installerJar.Name(), "--installServer", dest)
	}

	// Run installer
	return util.RunCommand("java", args...)
}

// Downloads the Forge Client Installer tool.
// The temporary file should be removed after usage.
func copyClientInstallTool() (string, error) {
	client, err := modernclientResource()
	if err != nil {
		return "", err
	}

	tmp, err := ioutil.TempFile("", "forgetool*.jar")
	if err != nil {
		return "", err
	}
	defer tmp.Close()

	if _, err := io.Copy(tmp, bytes.NewReader(client)); err != nil {
		return "", err
	}

	return tmp.Name(), nil
}
