//
// file.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Feb  8 2015.
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of this source code package.
//

package sprocess

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

type File struct {
	Prefix string
	Suffix string
	Dir    string
	Name
}

func (s *File) join(path string) string {
	return filepath.Join(s.Dir, path)
}

func (s *File) Start() error {
	if s.Dir == "" {
		return errors.New("attribute 'dir' is empty")
	}
	if err := os.MkdirAll(s.Dir, 0700); err != nil {
		return err
	}
	return nil
}

func (s *File) NewWriter(id string, d *Data) (io.WriteCloser, error) {
	return os.OpenFile(s.join(s.Prefix+id+s.Suffix), os.O_RDWR|os.O_CREATE, 0600)
}

func (s *File) NewReader(id string, d *Data) (io.ReadCloser, error) {
	return os.OpenFile(s.join(s.Prefix+id+s.Suffix), os.O_RDONLY, 0400)
}

func (s *File) Delete(id string, d *Data) error {
	return os.Remove(s.join(s.Prefix + id + s.Suffix))
}
