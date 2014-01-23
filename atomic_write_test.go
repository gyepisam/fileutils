// Copyright 2014 Gyepi Sam. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fileutils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func dirFileName() (string, string) {
	dir, err := ioutil.TempDir("", "atomic-write-test-")
	if err != nil {
		panic(err)
	}
	f, err := ioutil.TempFile(dir, "rand-file-")
	if err != nil {
		panic(err)
	}
	f.Close()
	os.Remove(f.Name())
	return dir, f.Name()
}

func dirCount(dir string) int {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	return len(entries)
}

type testCase struct {
	Name     string
	In       []byte
	Out      []byte
	Err      error
	Dir      string
	Filename string
}

func newCase(name string, in, out int, err error) *testCase {
	t := &testCase{Name: name, Err: err}
	t.Dir, t.Filename = dirFileName()
	if in > 0 {
		t.In = []byte(strings.Repeat("X", in))
	}
	if out > 0 {
		t.Out = []byte(strings.Repeat("123", out))
	}
	return t
}

func run(t *testing.T, obj *testCase) {

	// If there is input, we create one input file.
	// Otherwise, if the operation does not return an error, we expect to have one file.
	fileCount := 0
	if len(obj.In) > 0 {
		err := ioutil.WriteFile(obj.Filename, obj.In, 0655)
		if err != nil {
			panic(err)
		}
		fileCount++
	} else if obj.Err == nil {
		fileCount++
	}

	err := AtomicWrite(obj.Filename, func(f *os.File) error {
		n, err := f.Write(obj.Out)
		if err != nil {
			return err
		}
		if z := len(obj.Out); z != n {
			return fmt.Errorf("Short write: %s -- %d < %d", obj.Name, n, z)
		}
		return obj.Err
	})

	if err != obj.Err {
		t.Errorf("%s: Expected error: %s, but got: %s\n", obj.Name, obj.Err, err)
		return
	}

	// no littering
	if got := dirCount(obj.Dir); fileCount != got {
		t.Errorf("%s: Expected %d object in dir, but got: %d", obj.Name, fileCount, got)
		return
	}

	gotb, err := ioutil.ReadFile(obj.Filename)
	if os.IsNotExist(err) && len(obj.In) == 0 {
		return
	}

	if err != nil {
		panic(err)
	}

	expb := obj.In
	if obj.Err == nil {
		expb = obj.Out
	}

	if e, g := string(expb), string(gotb); e != g {
		t.Errorf("%s: Expected %s in file, got %s", obj.Name, e, g)
	}
}

func TestAtomicWrite(t *testing.T) {
	testCases := []*testCase{
		newCase("No File -> No File", 0, 0, errors.New("No File Case")),
		newCase("No File -> Has File", 0, 12, nil),
		newCase("Has File -> Changed File", 2, 13, nil),
		newCase("Has File -> Unchanged File", 12, 23, errors.New("Unchanged File Case")),
	}

	for _, c := range testCases {
		run(t, c)
		os.RemoveAll(c.Dir)
	}
}
