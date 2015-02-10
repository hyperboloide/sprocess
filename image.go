//
// resize.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Feb  8 2015.
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of this source code package.
//

package sprocess

import (
	"errors"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
)

type ImageOperation int

const (
	ImageThumbnail ImageOperation = iota
	ImageResize
)

type Image struct {
	Operation     ImageOperation
	Height        uint
	Width         uint
	Interpolation string
	Output        string
	Name          string
	interpolation resize.InterpolationFunction
}

func (i *Image) GetName() string {
	return i.Name
}

func (i *Image) Start() error {

	switch i.Operation {
	case ImageThumbnail:
		if i.Height == 0 || i.Width == 0 {
			return errors.New("height and width cannot be equal to 0")
		}
	case ImageResize:
		if i.Height == 0 && i.Width == 0 {
			return errors.New("height and width cannot be both equal to 0")
		}
	default:
		return errors.New("invalid image operation")
	}

	switch i.Interpolation {
	case "", "NearestNeighbor":
		i.interpolation = resize.NearestNeighbor
	case "Bilinear":
		i.interpolation = resize.Bilinear
	case "Bicubic":
		i.interpolation = resize.Bicubic
	case "MitchellNetravali":
		i.interpolation = resize.MitchellNetravali
	case "Lanczos2":
		i.interpolation = resize.Lanczos2
	case "Lanczos3":
		i.interpolation = resize.Lanczos3
	default:
		return errors.New(fmt.Sprintf("unknow interpolation algorithm '%s'", i.Interpolation))
	}

	switch i.Output {
	case "", "jpg", "png", "gif":
	default:
		return errors.New(fmt.Sprintf("unsuported output format '%s'", i.Output))
	}
	return nil
}

func (i *Image) Encode(r io.Reader, w io.Writer, d *Data) error {
	img, _, err := image.Decode(r)
	if err != nil {
		return err
	}
	var newImage image.Image
	if i.Operation == ImageResize {
		newImage = resize.Resize(i.Width, i.Height, img, i.interpolation)
	} else {
		newImage = resize.Thumbnail(i.Width, i.Height, img, i.interpolation)
	}
	switch i.Output {
	case "jpg", "":
		err = jpeg.Encode(w, newImage, nil)
	case "png":
		err = png.Encode(w, newImage)
	case "gif":
		err = gif.Encode(w, newImage, nil)
	}
	return err
}
