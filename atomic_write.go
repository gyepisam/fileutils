// Copyright 2014 Gyepi Sam. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fileutils

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

//AtomicWrite allows a function to write atomically to a named file.
//It first opens a file handle to a temporary file in the same directory as the target
//file and invokes the function, with the file handle as a parameter.
//If the function returns nil, AtomicWrite renames the temporary file over the target file.
//Any necessary directories to the target file are created.
func AtomicWrite(path string, writer func(*os.File) error) (err error) {

	dir := filepath.Dir(path)
	name := filepath.Base(path)

	if err = os.MkdirAll(dir, 0755); err != nil {
		return
	}

	f, err := ioutil.TempFile(dir, name)
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			os.Remove(f.Name())
		}
	}()

	err = writer(f)
	if err != nil {
		return
	}

	err = f.Close()
	if err != nil {
		return
	}

	err = os.Rename(f.Name(), path)
	return
}
