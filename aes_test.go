package sprocess_test

import (
	"bytes"
	. "github.com/hyperboloide/sprocess"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Aes", func() {

	testBin := genBlob(1 << 22)
	out1 := new(bytes.Buffer)
	data := NewData()

	aes := &AES{
		Base64String: "ETl5QyPnHfi+vF4HrZfFvO2Julv4LVL7HNB1N7vkLGU=",
		Name:         "aes",
	}

	It("should Encode", func() {
		Ω(aes.Start()).To(BeNil())
		Ω(aes.Encode(
			bytes.NewReader(testBin),
			out1,
			data)).To(BeNil())
		Ω(bytes.Equal(out1.Bytes(), testBin)).To(BeFalse())
		d, err := data.Get("iv")
		Ω(err).To(BeNil())
		Ω(len(d.(string)) > 0).To(BeTrue())
	})

	out2 := new(bytes.Buffer)
	It("should Decode", func() {
		Ω(aes.Start()).To(BeNil())
		Ω(aes.Decode(
			bytes.NewReader(out1.Bytes()),
			out2,
			data)).To(BeNil())
		Ω(bytes.Equal(out2.Bytes(), testBin)).To(BeTrue())
	})
})
