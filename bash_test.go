package sprocess_test

import (
	"bytes"
	. "github.com/hyperboloide/sprocess"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bash", func() {

	testBin := genBlob(1 << 22)

	It("should Encode", func() {

		out1 := new(bytes.Buffer)
		data := NewData()
		zip := &Bash{
			Cmd:  "gzip",
			Name: "zip",
		}

		Ω(zip.Start()).To(BeNil())
		Ω(zip.Encode(
			bytes.NewReader(testBin),
			out1,
			data)).To(BeNil())
		Ω(bytes.Equal(out1.Bytes(), testBin)).To(BeFalse())

		out2 := new(bytes.Buffer)
		unzip := &Bash{
			Cmd:  "gzip -d",
			Name: "unzip",
		}

		Ω(unzip.Start()).To(BeNil())
		Ω(unzip.Encode(
			bytes.NewReader(out1.Bytes()),
			out2,
			data)).To(BeNil())
		Ω(bytes.Equal(out2.Bytes(), testBin)).To(BeTrue())
	})

	It("should crash", func() {
		out1 := new(bytes.Buffer)
		data := NewData()
		crash := &Bash{
			Cmd:  "exit 1",
			Name: "crash",
		}

		Ω(crash.Start()).To(BeNil())
		Ω(crash.Encode(
			bytes.NewReader(testBin),
			out1,
			data)).ToNot(BeNil())
	})

})
