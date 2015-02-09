//
// pipeline.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Feb  8 2015.
// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
//

package sprocess

import (
	"io"
	"sync"
)

type Pipeline struct {
	lastReader  io.ReadCloser
	errorChain  []chan error
	cancelChain []chan struct{}
	Errors      []error
	size        int
}

type EncodingPipeline struct {
	Encoders []Encoder
	Output   Outputer
}

type DecodingPipeline struct {
	Decoders []Decoder
	Input    Inputer
}

func (p *Pipeline) mergeErrors(chans []chan error) <-chan error {
	var wg sync.WaitGroup
	out := make(chan error, 1)

	output := func(c <-chan error) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(p.size)
	for _, c := range chans {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func (p *Pipeline) stopAll() {
	for _, c := range p.cancelChain {
		close(c)
	}
}

func (p *Pipeline) Exec(w io.Writer) error {
	defer p.stopAll()

	errCh := make(chan error)
	go func() {
		_, err := io.Copy(w, p.lastReader)
		errCh <- err
		close(errCh)
	}()
	errCopy := <-errCh
	if errCopy != nil {
		return errCopy
	}

	p.Errors = make([]error, p.size)
	i := 0
	var firstError error
	for n := range p.mergeErrors(p.errorChain) {
		p.Errors[i] = n
		if firstError == nil && p.Errors[i] != nil {
			firstError = p.Errors[i]
			break
		}
		i += 1
	}
	return firstError
}

func newPipeline(r io.ReadCloser, size int) *Pipeline {
	p := &Pipeline{
		lastReader:  r,
		errorChain:  make([]chan error, size),
		cancelChain: make([]chan struct{}, size),
		size:        size,
	}
	return p
}

func NewEncoding(encoders []Encoder, r io.ReadCloser, d *Data) *Pipeline {
	p := newPipeline(r, len(encoders))
	for i, e := range encoders {
		errCh := make(chan error, 1)
		cancelCh := make(chan struct{}, 1)
		nextReader, writer := io.Pipe()

		go func(e Encoder, r io.ReadCloser, w *io.PipeWriter, d *Data) {
			defer w.Close()
			defer close(errCh)

			go func() {
				select {
				case <-cancelCh:
					r.Close()
				}
			}()

			errCh <- e.Encode(r, w, d)
		}(e, p.lastReader, writer, d)

		p.lastReader = nextReader
		p.errorChain[i] = errCh
		p.cancelChain[i] = cancelCh
	}
	return p
}

func NewDecoding(decoders []Decoder, r io.ReadCloser, d *Data) *Pipeline {
	p := newPipeline(r, len(decoders))
	for i, e := range decoders {
		errCh := make(chan error, 1)
		cancelCh := make(chan struct{}, 1)
		nextReader, writer := io.Pipe()

		go func(e Decoder, r io.ReadCloser, w *io.PipeWriter, d *Data) {
			defer w.Close()
			defer close(errCh)

			go func() {
				select {
				case <-cancelCh:
					r.Close()
				}
			}()

			errCh <- e.Decode(r, w, d)
		}(e, p.lastReader, writer, d)

		p.lastReader = nextReader
		p.errorChain[i] = errCh
		p.cancelChain[i] = cancelCh
	}
	return p
}
