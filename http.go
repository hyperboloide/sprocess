//
// http.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Feb  8 2015.
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of this source code package.
//

package sprocess

import (
	"errors"
	"github.com/dchest/uniuri"
	"io"
	"net/http"
)

type HTTP struct {
	initialized bool

	Encoders []Encoder
	Decoders []Decoder

	Input  Inputer
	Output Outputer
	Delete Deleter
}

func GenId() string {
	return uniuri.New()
}

func badRequest(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func internalError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (h *HTTP) Encode(w http.ResponseWriter, r *http.Request, id string) (map[string]interface{}, error) {
	if h.Output == nil {
		internalError(w)
		return nil, errors.New("No Output provided")
	}

	mr, err := r.MultipartReader()
	if err != nil {
		badRequest(w)
		return nil, err
	}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			badRequest(w)
			return nil, err
		}
		filename := part.FileName()
		if filename == "" {
			continue
		} else {
			data := NewData()
			data.Set("filename", filename)
			service := &Service{
				EncodingPipe: &EncodingPipeline{
					Encoders: h.Encoders,
					Output:   h.Output,
				},
			}
			if err = service.Encode(id, part, data); err != nil {
				internalError(w)
				return nil, err
			}
			w.WriteHeader(http.StatusCreated)
			return data.Export(), nil
		}
	}
	badRequest(w)
	return nil, errors.New("No file in request")
}

func (h *HTTP) Decode(w http.ResponseWriter, r *http.Request, dataMap map[string]interface{}) error {
	if h.Input == nil {
		internalError(w)
		return errors.New("No Input provided")
	}
	service := &Service{
		DecodingPipe: &DecodingPipeline{
			Decoders: h.Decoders,
			Input:    h.Input,
		},
	}

	data := NewDataFrom(dataMap)
	idIf, err := data.Get("identifier")
	if err != nil {
		internalError(w)
		return err
	}
	id, ok := idIf.(string)
	if ok == false {
		internalError(w)
		return errors.New("Key 'identifier' is not a string")
	}

	pr, pw := io.Pipe()

	go io.Copy(w, pr)

	if err := service.Decode(id, pw, data); err != nil {
		internalError(w)
		return err
	}
	return nil
}

func (h *HTTP) Remove(w http.ResponseWriter, r *http.Request, dataMap map[string]interface{}) error {

	data := NewDataFrom(dataMap)
	idIf, err := data.Get("identifier")
	if err != nil {
		internalError(w)
		return err
	}
	id, ok := idIf.(string)
	if ok == false {
		internalError(w)
		return errors.New("Key 'identifier' is not a string")
	}

	if err := h.Delete.Delete(id, data); err != nil {
		internalError(w)
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}
