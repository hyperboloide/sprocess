//
// sprocess.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Nov  8 2014.
//

package sprocess

import (
	//	"errors"
	//	"fmt"
	"io"
)

type DataMap map[string]interface{}

type Base interface {
	GetName() string
	Start() error
}

type Encoder interface {
	Base
	Encode(io.Reader, io.Writer, *Data) error
}

type Decoder interface {
	Base
	Decode(io.Reader, io.Writer, *Data) error
}

type EncodeDecoder interface {
	Base
	Encode(io.Reader, io.Writer, *Data) error
	Decode(io.Reader, io.Writer, *Data) error
}

type Outputer interface {
	Base
	NewWriter(string, *Data) (io.WriteCloser, error)
}

type Inputer interface {
	Base
	NewReader(string, *Data) (io.ReadCloser, error)
}

type Deleter interface {
	Base
	Delete(string, *Data) error
}
