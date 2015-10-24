//
// aes.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Feb  8 2015.
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of this source code package.
//

package sprocess

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

type AES struct {
	block        cipher.Block
	Key          []byte
	Base64String string
	Name
}

func (c *AES) Start() error {
	var err error

	switch {
	case c.Key != nil:
	case c.Base64String != "":
		c.Key, err = base64.StdEncoding.DecodeString(c.Base64String)
		if err != nil {
			return err
		}
	default:
		return errors.New("needs a cycpher key")
	}
	block, err := aes.NewCipher(c.Key)
	if err != nil {
		return err
	}
	c.block = block
	return nil
}

func (c *AES) Encode(r io.Reader, w io.Writer, d *Data) error {
	iv := generateIV(c.block.BlockSize())
	d.Set("iv", base64.StdEncoding.EncodeToString(iv))

	stream := cipher.NewCFBEncrypter(c.block, iv)
	writer := &cipher.StreamWriter{S: stream, W: w}
	if _, err := io.Copy(writer, r); err != nil {
		return err
	}
	return nil
}

func (c *AES) Decode(r io.Reader, w io.Writer, d *Data) error {
	iv64, err := d.Get("iv")
	if err != nil {
		return err
	}
	if iv64.(string) == "" {
		return errors.New("no initialization vector provided")
	}
	iv, err := base64.StdEncoding.DecodeString(iv64.(string))
	if err != nil {
		return err
	}
	stream := cipher.NewCFBDecrypter(c.block, iv)
	reader := &cipher.StreamReader{S: stream, R: r}
	if _, err := io.Copy(w, reader); err != nil {
		return err
	}
	return nil
}

func generateIV(size int) []byte {
	b := make([]byte, size)
	rand.Read(b)
	return b
}
