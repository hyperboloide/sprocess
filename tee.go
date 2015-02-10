//
// tee.go
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

type Tee struct {
	Encoders []Encoder
	Output   Outputer
	Name     string
}

func (t *Tee) GetName() string {
	return t.Name
}

func (t *Tee) Start() error {
	if t.Output == nil {
		return errors.New("Output cannot be nil")
	}
	return nil
}

func (t *Tee) Encode(r io.Reader, w io.Writer, d *Data) error {
	tr, tw := io.Pipe()
	reader := io.TeeReader(r, tw)

	idIf, err := d.Get("identifier")
	if err != nil {
		return err
	}
	id, ok := idIf.(string)
	if ok == false {
		return errors.New("Identifier is not of type string")
	}

	service := &Service{
		EncodingPipe: &EncodingPipeline{
			Encoders: t.Encoders,
			Output:   t.Output,
		},
	}

	var teeData *Data = NewData()
	teeData.Set("identifier", id)
	errorTee := make(chan error, 1)

	go func(r io.ReadCloser, data *Data) {
		err := service.Encode(id, r, data)
		errorTee <- err
		close(errorTee)
	}(tr, teeData)

	_, err = io.Copy(w, reader)
	tw.Close()
	if err != nil {
		return err
	}
	err = <-errorTee
	d.Set(t.Name, teeData)
	return err
}
