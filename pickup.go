package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
)

type pickup string

func (p pickup) submit(msg *message) (err error) {
	file, err := newTempFile()
	if err != nil {
		return
	}

	encode(file, msg)
	err = file.Close()
	if err != nil {
		return
	}

	// move to pickup directory
	// Instead of creating the file elsewhere and then moving to the pickup directory,
	// it is also possible to create the file directly in it. However, since PowerMTA
	// must repeatedly attempt to lock the file for exclusive access, it is less efficient
	// to do so.
	name := file.Name()
	dest := filepath.Join(string(p), filepath.Base(name))
	err = os.Rename(name, dest)
	if err != nil {
		return
	}
	return
}

var tempDir = os.TempDir()

func newTempFile() (file *os.File, err error) {
	for n := 0; n < 10; n++ {
		name := filepath.Join(tempDir, randName())
		file, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
		if err == nil || !os.IsExist(err) {
			return
		}
	}
	return
}

func randName() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x.eml", b)
}
