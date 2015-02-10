//
// compress.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Feb  8 2015.
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of this source code package.
//

package sprocess

import (
	"compress/gzip"
	"io"
)

type Gzip struct {
	Algo  string
	Name  string
	level int
}

func (c *Gzip) GetName() string {
	return c.Name
}

func (c *Gzip) Start() error {
	c.level = gzip.DefaultCompression
	switch c.Algo {
	case "best":
		c.level = gzip.BestCompression
	case "speed":
		c.level = gzip.BestSpeed
	}
	return nil
}

func (c *Gzip) Encode(r io.Reader, w io.Writer, d *Data) error {
	gzw, err := gzip.NewWriterLevel(w, c.level)
	if err != nil {
		return err
	}
	defer gzw.Close()
	_, err = io.Copy(gzw, r)
	return err
}

func (c *Gzip) Decode(r io.Reader, w io.Writer, d *Data) error {
	gzr, err := gzip.NewReader(r)
	defer gzr.Close()
	_, err = io.Copy(w, gzr)
	return err
}
