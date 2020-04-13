// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package util

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// Downloads the file, copying it to the given writer.
func Download(dst io.Writer, req *http.Request) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(dst, resp.Body)
	return err
}

// Downloads the file, copying it to the given writer.
// The temporary file should be removed after usage.
func DownloadTemp(req *http.Request, pattern string) (*os.File, error) {
	file, err := ioutil.TempFile("", pattern)
	if err != nil {
		return nil, err
	}

	if err := Download(file, req); err != nil {
		return nil, err
	}

	return file, nil
}
