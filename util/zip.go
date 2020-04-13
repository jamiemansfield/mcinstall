package util

import (
	"archive/zip"
	"errors"
	"io"
	"os"
)

// Gets the *zip.File of the given name in the given zip file.
func GetFileInZip(zipFile *zip.Reader, name string) (*zip.File, error) {
	for _, file := range zipFile.File {
		if file.Name == name {
			return file, nil
		}
	}
	return nil, errors.New("zip file did not contain: " + name)
}

// Reads the given zip file, and writes it to the given destination.
func CopyZipFileToDisk(file *zip.File, dest string) error {
	r, err := file.Open()
	if err != nil {
		return err
	}
	defer r.Close()
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	return err
}
