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
	Prefix      string
	Suffix      string
	Dir         string
	AllowSub    bool
	RemoveEmpty bool
	Name        string
}

func (s *File) GetName() string {
	return s.Name
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
	name := s.Prefix + id + s.Suffix

	if filepath.Dir(name) == "." {
		return os.OpenFile(s.join(name), os.O_RDWR|os.O_CREATE, 0600)
	} else if s.AllowSub == false {
		return nil, errors.New("sub directories not allowed")
	}
	if err := os.MkdirAll(s.join(filepath.Dir(name)), 0700); err != nil {
		return nil, err
	}
	return os.OpenFile(s.join(name), os.O_RDWR|os.O_CREATE, 0600)
}

func (s *File) NewReader(id string, d *Data) (io.ReadCloser, error) {
	return os.OpenFile(s.join(s.Prefix+id+s.Suffix), os.O_RDONLY, 0400)
}

func (s *File) RemoveIfEmpty(dir string) error {
	if dir == "." {
		return nil
	} else {
		f, err := os.Open(s.join(dir))
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = f.Readdir(1)
		if err != io.EOF {
			return err
		}
	}
	os.Remove(s.join(dir))
	return s.RemoveIfEmpty(filepath.Dir(dir))
}

func (s *File) Delete(id string, d *Data) error {
	name := s.Prefix + id + s.Suffix
	if err := os.Remove(s.join(name)); err != nil {
		return err
	}
	if s.RemoveEmpty == true && filepath.Dir(name) != "." {
		return s.RemoveIfEmpty(filepath.Dir(name))
	}
	return nil
}
