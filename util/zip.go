package util

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
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

// Extracts the given zip file to the local file system.
func ExtractZipFileToDisk(zipFile *zip.Reader, dest string) error {
	for _, file := range zipFile.File {
		localPath := filepath.Join(dest, file.Name)

		if file.FileInfo().IsDir() {
			err := os.MkdirAll(localPath, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			err := CopyZipFileToDisk(file, localPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func MergeZips(into *zip.Writer, from *zip.Reader, files []string, filter func(string) bool) ([]string, error) {
	for _, file := range from.File {
		if filter != nil && filter(file.Name) {
			continue
		}
		if contains(files, file.Name) {
			continue
		}
		files = append(files, file.Name)

		if err := copyZipFile(into, file); err != nil {
			return nil, err
		}
	}

	return files, nil
}

func copyZipFile(into *zip.Writer, file *zip.File) error {
	r, err := file.Open()
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := into.Create(file.Name)
	if err != nil {
		return err
	}

	// Copy file contents
	if !file.FileInfo().IsDir() {
		_, err = io.Copy(w, r)
		if err != nil {
			return err
		}
	}

	return nil
}

func contains(files []string, file string) bool {
	for _, f := range files {
		if f == file {
			return true
		}
	}
	return false
}
