//
// checksum.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Feb 10 2015.
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of this source code package.
//

package sprocess

import (
	"crypto/md5"
	"crypto/sha1"
	"errors"
	"fmt"
	"hash"
	"io"
)

type CheckSum struct {
	Hash string
	Name
}

func (c *CheckSum) Start() error {
	switch c.Hash {
	case "", "md5", "sha1":
	default:
		return errors.New("Hash '" + c.Hash + "' is not supported")
	}
	if c.Name == "" {
		return errors.New("Checksum must have a name")
	}
	return nil
}

func (c *CheckSum) Encode(r io.Reader, w io.Writer, d *Data) error {
	var h hash.Hash
	switch c.Hash {
	case "", "sha1":
		h = sha1.New()
	case "md5":
		h = md5.New()
	default:
		return errors.New("Hash '" + c.Hash + "' is not supported")
	}

	reader := io.TeeReader(r, h)
	_, err := io.Copy(w, reader)
	if err != nil {
		return err
	}
	d.Set(c.GetName(), fmt.Sprintf("%x", h.Sum(nil)))
	return nil
}

func (c *CheckSum) Decode(r io.Reader, w io.Writer, d *Data) error {
	var h hash.Hash
	switch c.Hash {
	case "", "sha1":
		h = sha1.New()
	case "md5":
		h = md5.New()
	default:
		return errors.New("Hash '" + c.Hash + "' is not supported")
	}

	value, err := d.Get(c.GetName())
	if err != nil {
		return err
	}

	reader := io.TeeReader(r, h)
	_, err = io.Copy(w, reader)
	if err != nil {
		return err
	}
	if fmt.Sprintf("%x", h.Sum(nil)) != value {
		return errors.New("CheckSums do not match")
	}
	return nil
}
