//
// null.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Oct 24 2015.
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of this source code package.
//

package sprocess

import (
	"io"
)

type NullWriterCloser struct{}

func (nwc *NullWriterCloser) Write(p []byte) (int, error) {
	return len(p), nil
}

func (nwc *NullWriterCloser) Close() error {
	return nil
}

type Null struct {
	Name string
}

func (n *Null) GetName() string {
	return n.Name
}

func (n *Null) Start() error {
	return nil
}

func (n *Null) NewWriter(id string, d *Data) (io.WriteCloser, error) {
	return &NullWriterCloser{}, nil
}
