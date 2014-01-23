// Copyright 2014 Gyepi Sam. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fileutils

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

type Case struct {
	Name       string // filename
	Cleanup    func()
	ExpectFile bool
	ExpectDir  bool
}

// None Existent file
func NoneCase() *Case {
	f, err := ioutil.TempFile("", "fileutils-test-none-")
	if err != nil {
		panic(err)
	}
	defer os.Remove(f.Name())
	return &Case{Name: f.Name()}
}

// File
func FileCase() *Case {
	f, err := ioutil.TempFile("", "fileutils-test-file-")
	if err != nil {
		panic(err)
	}
	return &Case{Name: f.Name(), Cleanup: func() { os.Remove(f.Name()) }, ExpectFile: true}
}

// Directory
func DirCase() *Case {
	name, err := ioutil.TempDir("", "fileutils-test-dir-")
	if err != nil {
		panic(err)
	}
	return &Case{Name: name, Cleanup: func() { os.RemoveAll(name) }, ExpectDir: true}
}

// Existing entity, neither file nor directory
func FifoCase() *Case {
	f, err := ioutil.TempFile("", "fileutils-test-fifo-")
	if err != nil {
		panic(err)
	}
	name := f.Name()
	os.Remove(f.Name())
	cmd := exec.Command("mkfifo", name)
	b, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	if len(b) > 0 {
		panic(string(b))
	}
	return &Case{Name: name, Cleanup: func() { os.RemoveAll(name) }}
}

type testFunc struct {
	N    string
	F    func(string) (bool, error)
	File bool
	Dir  bool
}

func TestAll(t *testing.T) {

	testCases := []func() *Case{
		NoneCase,
		FileCase,
		DirCase,
		FifoCase,
	}

	testFuncs := []testFunc{
		{N: "IsFile", F: IsFile, File: true},
		{N: "FileExists", F: FileExists, File: true},
		{N: "IsDir", F: IsDir, Dir: true},
		{N: "DirExists", F: DirExists, Dir: true},
		{N: "Exists", F: Exists, Dir: true, File:true},
	}

	for _, cf := range testCases {
		c := cf()
		name := c.Name

		for _, tf := range testFuncs {
			exists, err := tf.F(name)
			if err != nil {
				t.Fatalf("%s %s", name, err)
			}

			if exists {
				if c.ExpectDir && tf.File && !tf.Dir {
					t.Errorf("unexpected directory exists: %s, %s", name, tf.N)
				}
				if c.ExpectFile && tf.Dir && !tf.File {
					t.Errorf("unexpected file exists: %s, %s", name, tf.N)
				}
			} else {
				if c.ExpectDir && tf.Dir {
					t.Errorf("expected directory does not exist: %s, %s", name, tf.N)
				}
				if c.ExpectFile && tf.File {
					t.Errorf("expected file does not exist: %s, %s", name, tf.N)
				}
			}
		}

		if c.Cleanup != nil {
			c.Cleanup()
		}

	}
}
