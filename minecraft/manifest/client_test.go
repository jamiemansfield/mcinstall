// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package manifest

import "testing"

func TestGetVersionManifest(t *testing.T) {
	manifest, err := GetVersionManifest(nil)
	if err != nil {
		t.Errorf("failed to get version manifest: %e", err)
		return
	}

	t.Logf("latest release: %s latest snapshot: %s", manifest.Latest.Release, manifest.Latest.Snapshot)
}

func TestGetVersionManifest_1_2_5(t *testing.T) {
	manifest, err := GetVersionManifest(nil)
	if err != nil {
		t.Errorf("failed to get version manifest: %e", err)
		return
	}

	versionInfo := manifest.FindVersion("1.2.5")
	if versionInfo == nil {
		t.Errorf("failed to find Minecraft 1.2.5")
		return
	}

	version, err := versionInfo.GetFull(nil)
	if err != nil {
		t.Errorf("failed to get Minecraft 1.2.5: %e", err)
		return
	}

	t.Logf("Minecraft %s", version.ID)
	t.Logf("Main class: %s", version.MainClass)
}
