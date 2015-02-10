package sprocess_test

import (
	"bytes"
	. "github.com/hyperboloide/sprocess"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Size", func() {

	var res int64
	res = 1 << 22
	testBin := genBlob(int(res))
	out1 := new(bytes.Buffer)
	data := NewData()

	size := &Size{
		Name: "size",
	}

	It("should Encode", func() {
		Ω(size.Start()).To(BeNil())
		Ω(size.Encode(
			bytes.NewReader(testBin),
			out1,
			data)).To(BeNil())
		Ω(bytes.Equal(out1.Bytes(), testBin)).To(BeTrue())
		s, err := data.Get("size")
		Ω(err).To(BeNil())
		Ω(s).To(Equal(res))
	})

	out2 := new(bytes.Buffer)
	It("should Decode", func() {
		Ω(size.Decode(
			bytes.NewReader(out1.Bytes()),
			out2,
			data)).To(BeNil())
		Ω(bytes.Equal(out2.Bytes(), testBin)).To(BeTrue())
		s, err := data.Get("size")
		Ω(err).To(BeNil())
		Ω(s).To(Equal(res))

	})

})
