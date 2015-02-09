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

	testBin := make([]byte, 1<<16)
	rand.Read(testBin)
	data := NewData()
	s3 := &S3Bucket{
		Bucket: "sprocess",
		Name:   "s3",
		Domain: "s3-eu-west-1.amazonaws.com",
	}

	id := GenId()

	It("should Write", func() {
		Ω(s3.Init()).To(BeNil())
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
})
