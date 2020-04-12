// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package forgeinstall

import (
	"fmt"
	"github.com/jamiemansfield/ftbinstall/mcinstall"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// Installs Minecraft Forge to the given destination, for the given target.
// If the target is Server, the destination will be the root directory of
// the server; if the target is Client, the destination will be the
// launcher's root directory.
func InstallForge(target mcinstall.InstallTarget, dest string, mcVersion *mcinstall.McVersion, forgeVersion string) error {
	fmt.Printf("Installing Minecraft Forge %s-%s...\n", mcVersion, forgeVersion)

	// Use modern installer - Minecraft 1.13 and above
	if mcVersion.Major >= 1 && mcVersion.Minor >= 13 {
		return installModernForge(target, dest, mcVersion, forgeVersion)
	} else
	// Use universal install method - Minecraft 1.5 -> Minecraft 1.12
	if mcVersion.Major >= 1 && mcVersion.Minor >= 5 && mcVersion.Minor <= 12 {
		return installUniversalForge(target, dest, mcVersion, forgeVersion)
	}
	// todo: older

	return nil
}

// Downloads the Minecraft Forge installer for the given version (MC-Forge),
// to a temporary file.
// The temporary file should be removed after usage.
func downloadForgeInstaller(version string) (*os.File, error) {
	url := MavenRoot + "net/minecraftforge/forge/" + version + "/forge-" + version + "-installer.jar"

	// Download installer
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/java,application/java-archive,application/x-java-archive")
	req.Header.Set("User-Agent", "ftbinstall/0.1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	file, err := ioutil.TempFile("", "forge*.jar")
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(file, resp.Body); err != nil {
		return nil, err
	}

	return file, nil
}

// Runs the given Minecraft Forge installer, installing for the given target,
// to the given target directory.
func runInstaller(target mcinstall.InstallTarget, dest string, installerJar *os.File) error {
	// Use the absolute directory
	destination, err := filepath.Abs(dest)
	if err != nil {
		return err
	}

	// Create the appropriate arguments for the install target
	var args []string
	args = append(args, "-jar", installerJar.Name())
	if target == mcinstall.Client {
		args = append(args, "--extract", destination)
	} else {
		args = append(args, "--installServer", destination)
	}

	// Run the command
	cmd := exec.Command("java", args...)

	// todo: improve
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Print(string(out))

	return nil
}
