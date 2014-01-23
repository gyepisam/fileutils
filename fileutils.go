// Copyright 2014 Gyepi Sam. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fileutils

import (
	"os"
)

// fileTest stats a file and runs a boolean function on the resulting finfo structure.
// If the stat fails due to a non-existent file, fileTest returns false and a nil error.
func fileTest(name string, f func(os.FileInfo) bool) (bool, error) {
	info, err := os.Stat(name)

	if err == nil {
		return f(info), nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// IsFile returns a boolean denoting whether name represents a file.
// A non-existent file returns false.
func IsFile(name string) (bool, error) {
	return fileTest(name, func(info os.FileInfo) bool {
		return info.Mode()&os.ModeType == 0
	})
}

// IsDir returns a boolean denoting whether name represents a directory.
// A non-existent directory returns false.
func IsDir(name string) (bool, error) {
	return fileTest(name, func(info os.FileInfo) bool {
		return info.IsDir()
	})
}

// Exists returns a boolean denoting whether the name exists in the file system.
// A non-existent entity returns false.
func Exists(name string) (bool, error) {
	return fileTest(name, func(info os.FileInfo) bool {
		return true
	})
}

// FileExists is an alias for IsFile.
func FileExists(name string) (bool, error) {
	return IsFile(name)
}

// DirExists is an alias for IsDir.
func DirExists(name string) (bool, error) {
	return IsDir(name)
}
