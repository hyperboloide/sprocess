package sprocess_test

import (
	"bytes"
	. "github.com/hyperboloide/sprocess"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Compress", func() {

	testBin := genBlob(1 << 22)
	out1 := new(bytes.Buffer)
	data := NewData()

	gz := &Gzip{
		Algo: "speed",
	}

	It("should Encode", func() {
		Ω(gz.Start()).To(BeNil())
		Ω(gz.Encode(
			bytes.NewReader(testBin),
			out1,
			data)).To(BeNil())
		output := out1.Bytes()
		Ω(bytes.Equal(output, testBin)).To(BeFalse())
		Ω(len(output) < len(testBin)).To(BeTrue())
	})

	out2 := new(bytes.Buffer)
	It("should Decode", func() {
		Ω(gz.Decode(
			bytes.NewReader(out1.Bytes()),
			out2,
			data)).To(BeNil())
		Ω(bytes.Equal(out2.Bytes(), testBin)).To(BeTrue())
	})

})
