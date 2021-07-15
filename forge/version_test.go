// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package forge

import "testing"

func TestParseVersion(t *testing.T) {
	// 14.23.5.2855
	{
		version, err := ParseVersion("14.23.5.2855")
		if err != nil {
			t.Fatal(err)
		}

		if version.Build != 2855 {
			t.Errorf("Build is %d, should be %d", version.Build, 2855)
		}
	}
}
