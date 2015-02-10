package sprocess_test

import (
	"bytes"
	"crypto/rand"
	. "github.com/hyperboloide/sprocess"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
)

var _ = Describe("S3bucket", func() {

	testBin := make([]byte, 1<<8)
	rand.Read(testBin)
	data := NewData()
	s3 := &S3Bucket{
		Bucket: "sprocess",
		Name:   "s3",
		Domain: "s3-eu-west-1.amazonaws.com",
	}

	id := GenId()

	It("should Write", func() {
		Ω(s3.Start()).To(BeNil())
		w, err := s3.NewWriter(id, data)
		Ω(err).To(BeNil())
		Ω(w).ToNot(BeNil())
		l, err := io.Copy(w, bytes.NewReader(testBin))
		w.Close()
		Ω(err).To(BeNil())
		Ω(len(testBin) == int(l)).To(BeTrue())
	})

	It("should read", func() {
		r, err := s3.NewReader(id, data)
		Ω(err).To(BeNil())
		Ω(r).ToNot(BeNil())
		out1 := new(bytes.Buffer)
		l, err := io.Copy(out1, r)
		Ω(err).To(BeNil())
		r.Close()
		Ω(len(testBin) == int(l)).To(BeTrue())
		Ω(bytes.Equal(testBin, out1.Bytes())).To(BeTrue())
	})

	It("should delete", func() {
		Ω(s3.Delete(id, data)).To(BeNil())

		_, err := s3.NewReader(id, data)
		Ω(err).ToNot(BeNil())
	})

	It("should do service with s3", func() {
		data := NewData()
		id := "pic.jpg"

		size := &Size{
			Name: "size",
		}
		Ω(size.Start()).To(BeNil())

		img := &Image{
			Operation: ImageResize,
			Height:    100,
			Output:    "jpg",
			Name:      "resize",
		}
		Ω(img.Start()).To(BeNil())

		service := &Service{
			EncodingPipe: &EncodingPipeline{
				Encoders: []Encoder{img, size},
				Output:   s3,
			},
			DecodingPipe: &DecodingPipeline{
				Decoders: []Decoder{size},
				Input:    s3,
			},
		}

		Ω(service.Encode(id, testFileReader(), data)).To(BeNil())

		out := new(bytes.Buffer)
		r, w := io.Pipe()
		go func() {
			io.Copy(out, r)
		}()
		Ω(service.Decode(id, w, data)).To(BeNil())

	})

})
