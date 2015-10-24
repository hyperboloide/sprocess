package sprocess_test

import (
	"bytes"
	"crypto/rand"
	. "github.com/hyperboloide/sprocess"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
)

var _ = Describe("GoogleCloud", func() {

	testBin := make([]byte, 1<<8)
	rand.Read(testBin)
	data := NewData()
	gc := &GoogleCloud{
		ProjectId:   "sprocess-1108",
		Bucket:      "sprocess",
		JsonKeyPath: "google_cloud_key.json",
		Name:        "google_cloud",
	}
	id := GenId()

	It("should Write", func() {
		Ω(gc.Start()).To(BeNil())
		w, err := gc.NewWriter(id, data)
		Ω(err).To(BeNil())
		Ω(w).ToNot(BeNil())
		l, err := io.Copy(w, bytes.NewReader(testBin))
		w.Close()
		Ω(err).To(BeNil())
		Ω(len(testBin) == int(l)).To(BeTrue())
	})

	It("should read", func() {
		r, err := gc.NewReader(id, data)
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
		Ω(gc.Delete(id, data)).To(BeNil())

		_, err := gc.NewReader(id, data)
		Ω(err).ToNot(BeNil())
	})

	It("should do service with google cloud", func() {
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
				Output:   gc,
			},
			DecodingPipe: &DecodingPipeline{
				Decoders: []Decoder{size},
				Input:    gc,
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
