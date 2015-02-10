package sprocess_test

import (
	"bytes"
	"crypto/rand"
	. "github.com/hyperboloide/sprocess"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
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

	It("should Decode", func() {
		out2 := new(bytes.Buffer)
		Ω(aes.Decode(
			bytes.NewReader(out1.Bytes()),
			out2,
			data)).To(BeNil())
		Ω(bytes.Equal(out2.Bytes(), testBin)).To(BeTrue())
	})

	It("should work with []byte Key", func() {
		out1 := new(bytes.Buffer)
		data := NewData()
		key := make([]byte, 32)
		rand.Read(key)

		aes := &AES{
			Key:  key,
			Name: "aes",
		}

		Ω(aes.Start()).To(BeNil())
		Ω(aes.Encode(
			bytes.NewReader(testBin),
			out1,
			data)).To(BeNil())
		Ω(bytes.Equal(out1.Bytes(), testBin)).To(BeFalse())
		d, err := data.Get("iv")
		Ω(err).To(BeNil())
		Ω(len(d.(string)) > 0).To(BeTrue())

		out2 := new(bytes.Buffer)
		Ω(aes.Start()).To(BeNil())
		Ω(aes.Decode(
			bytes.NewReader(out1.Bytes()),
			out2,
			data)).To(BeNil())
		Ω(bytes.Equal(out2.Bytes(), testBin)).To(BeTrue())

	})

	It("should do service with aes", func(){
		data := NewData()
		id := "encrypted"
		
		fs := &File{
			Dir: "/tmp/" + GenId(),
			Name: "fs",
		}
		Ω(fs.Start()).To(BeNil())
		
		key := make([]byte, 32)
		rand.Read(key)
		aes := &AES{
			Key:  key,
			Name: "aes",
		}
		Ω(aes.Start()).To(BeNil())

		chck := &CheckSum{
			Name: "chck",
		}
		Ω(chck.Start()).To(BeNil())

		service := &Service{
			EncodingPipe: &EncodingPipeline{
				Encoders: []Encoder{chck, aes},				
				Output:   fs,
			},
			DecodingPipe: &DecodingPipeline{
				Decoders: []Decoder{aes, chck},
				Input:   fs,
			},
		}

		Ω(service.Encode(id, testFileReader(), data)).To(BeNil())
		d := data.Export()
		_, exists := d["chck"]
		Ω(exists).To(BeTrue())

		out := new(bytes.Buffer)
		r, w := io.Pipe()
		go func() {
			io.Copy(out, r)
		}()
		Ω(service.Decode(id, w, data)).To(BeNil())

		
	})
	
})
