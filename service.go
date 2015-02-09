//
// service.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Feb  8 2015.
// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
//

package sprocess

import (
	"errors"
	"io"
)

type Access int

var NoEncodindError = errors.New("No encoding pipeline defined")
var NoDecodindError = errors.New("No decoding pipeline defined")

type Service struct {
	EncodingPipe *EncodingPipeline
	DecodingPipe *DecodingPipeline
}

func (s *Service) Encode(id string, r io.ReadCloser, data *Data) error {
	defer data.Set("identifier", id)
	
	if s.EncodingPipe == nil {
		return NoEncodindError
	}
	w, err := s.EncodingPipe.Output.NewWriter(id, data)
	if err != nil {
		return err
	}
	if len(s.EncodingPipe.Encoders) == 0 {
		
		l, err := io.Copy(w, r)
		if err != nil {
			return err
		}
		w.Close()
		data.Set("size", l)
	} else {
		p := NewEncoding(s.EncodingPipe.Encoders, r, data)
		if err := p.Exec(w); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) Decode(id string, w io.WriteCloser, data *Data) error {
	if s.DecodingPipe == nil {
		return NoDecodindError
	}
	r, err := s.DecodingPipe.Input.NewReader(id, data)
	if err != nil {
		return err
	}
	if len(s.DecodingPipe.Decoders) == 0 {
		done := make(chan error)

		go func() {
			defer r.Close()
			_, err := io.Copy(w, r)
			done <- err
		}()
		err := <-done
		if err != nil {
			return err
		}
	} else {
		p := NewDecoding(s.DecodingPipe.Decoders, r, data)
		if err := p.Exec(w); err != nil {
			return err
		}
	}
	return nil
}