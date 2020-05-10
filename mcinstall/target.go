// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mcinstall

import (
	"bytes"
	"encoding/json"
)

type InstallTarget int

const (
	Client InstallTarget = iota
	Server
)
var _ json.Marshaler = (*InstallTarget)(nil)
var _ json.Unmarshaler = (*InstallTarget)(nil)

func (i InstallTarget) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	switch i {
	case Server:
		buffer.WriteString("server")
	default:
		buffer.WriteString("client")
	}
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (i *InstallTarget) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	switch j {
	case "server":
		*i = Server
	default:
		*i = Client
	}
	return nil
}
