package sprocess_test

import (
	"bytes"
	. "github.com/hyperboloide/sprocess"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"image"
)

var _ = Describe("Image", func() {

	It("should resize image", func() {
		i := &Image{
			Operation: ImageResize,
			Height:    300,
			Output:    "jpg",
			Name:      "resize",
		}
		out := new(bytes.Buffer)
		data := NewData()

		Ω(i.Start()).To(BeNil())
		Ω(i.Encode(
			testFileReader(),
			out,
			data)).To(BeNil())
		img, format, err := image.Decode(out)
		Ω(err).To(BeNil())
		Ω(format).To(Equal("jpeg"))
		Ω(img).ToNot(BeNil())
		Ω(img.Bounds().Size().Y).To(Equal(300))
		Ω(img.Bounds().Size().X).ToNot(Equal(300))
	})

	It("should do thumbnail image", func() {
		i := &Image{
			Operation: ImageThumbnail,
			Height:    300,
			Width:     100,
			Output:    "png",
			Name:      "thumbnail",
		}
		out := new(bytes.Buffer)
		data := NewData()

		Ω(i.Start()).To(BeNil())
		Ω(i.Encode(
			testFileReader(),
			out,
			data)).To(BeNil())
		img, format, err := image.Decode(out)
		Ω(err).To(BeNil())
		Ω(format).To(Equal("png"))
		Ω(img).ToNot(BeNil())
		Ω(img.Bounds().Size().Y).ToNot(Equal(300))
		Ω(img.Bounds().Size().X).To(Equal(100))
	})

})
