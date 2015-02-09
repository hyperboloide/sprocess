//
// compress.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Feb  8 2015.
// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
//

package sprocess

import (
	"compress/gzip"
	"io"
)

type Gzip struct {
	Algo string
	Name string
}

func (c *Gzip) GetName() string {
	return c.Name
}

func (c *Gzip) Start() error {
	return nil
}

func (c *Gzip) Encode(r io.Reader, w io.Writer, d *Data) error {
	level := gzip.DefaultCompression
	switch c.Algo {
	case "best":
		level = gzip.BestCompression
	case "speed":
		level = gzip.BestSpeed
	}
	gzw, err := gzip.NewWriterLevel(w, level)
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
