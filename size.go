//
// size.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Feb 10 2015.
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of this source code package.
//

package sprocess

import (
	"errors"
	"io"
)

type Size struct {
	Name
}

func (s *Size) Start() error {
	if s.Name == "" {
		return errors.New("Size must have a name")
	}
	return nil
}

func (s *Size) Encode(r io.Reader, w io.Writer, d *Data) error {
	size, err := io.Copy(w, r)
	d.Set(s.GetName(), size)
	return err
}

func (s *Size) Decode(r io.Reader, w io.Writer, d *Data) error {
	size, err := io.Copy(w, r)
	d.Set(s.GetName(), size)
	return err
}
