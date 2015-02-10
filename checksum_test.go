package sprocess_test

import (
	"bytes"
	. "github.com/hyperboloide/sprocess"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Checksum", func() {

	testBin := genBlob(1 << 22)
	out1 := new(bytes.Buffer)
	data := NewData()

	chck := &CheckSum{
		Name: "chck",
	}

	It("should Encode", func() {
		Ω(chck.Start()).To(BeNil())
		Ω(chck.Encode(
			bytes.NewReader(testBin),
			out1,
			data)).To(BeNil())
		output := out1.Bytes()
		Ω(bytes.Equal(output, testBin)).To(BeTrue())
		d := data.Export()
		_, exists := d["chck"]
		Ω(exists).To(BeTrue())
	})

	out2 := new(bytes.Buffer)
	It("should Decode", func() {
		Ω(chck.Decode(
			bytes.NewReader(out1.Bytes()),
			out2,
			data)).To(BeNil())
		Ω(bytes.Equal(out2.Bytes(), testBin)).To(BeTrue())
		d := data.Export()
		_, exists := d["chck"]
		Ω(exists).To(BeTrue())

	})

	It("should not Decode if mismatch", func() {
		out := new(bytes.Buffer)
		Ω(chck.Decode(testFileReader(), out, data)).ToNot(BeNil())
	})

})
