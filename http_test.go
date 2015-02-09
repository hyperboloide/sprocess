package sprocess_test

import (
	. "github.com/hyperboloide/sprocess"
	"crypto/rand"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"bytes"
	"mime/multipart"
	"os"
	"io/ioutil"
)

var _ = Describe("Http", func() {

	var data map[string]interface{}
	var id string
	var process *HTTP
	var fileContents []byte
	
	It("should start encoders, decoders, Outputer", func(){
		aesKey := make([]byte, 32)
		rand.Read(aesKey)

		aes := &AES{
			Key: aesKey,
			Name: "encrypt",
		}
		Ω(aes.Start()).To(BeNil())

		zip := &Gzip{
			Name: "compress",
		}
		Ω(zip.Start()).To(BeNil())

		file := &File{
			Dir:  "/tmp/" + GenId(),
			Name: "file",
		}
		Ω(file.Start()).To(BeNil())
		
		process = &HTTP{
			Encoders: []Encoder{zip, aes},
			Decoders: []Decoder{aes, zip},
			Input: file,
			Output: file,
			Delete: file,
		}
	})
	
	It("should POST file", func() {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer GinkgoRecover()
			id = GenId()
			
			d, err := process.Encode(w, r, id)
			Ω(err).To(BeNil())
			Ω(d).ToNot(BeNil())
			data = d
		}))
		defer ts.Close()

		file, err := os.Open("tests/test.jpg")
		Ω(err).To(BeNil())
		fileContents, err = ioutil.ReadAll(file)
		Ω(err).To(BeNil())
		fi, err := file.Stat()
		Ω(err).To(BeNil())
		file.Close()
		
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("somefile", fi.Name())
		Ω(err).To(BeNil())
		part.Write(fileContents)
		Ω(writer.Close()).To(BeNil())
				
		req, err := http.NewRequest("POST", ts.URL, body)
		Ω(err).To(BeNil())
		req.Header.Add("Content-Type", writer.FormDataContentType())

		client := &http.Client{}
		resp, err := client.Do(req)
		Ω(err).To(BeNil())
		Ω(resp.StatusCode).To(Equal(201))
		Ω(data["filename"]).To(Equal(fi.Name()))
		Ω(data["identifier"]).To(Equal(id))
	})
	
	It("sould GET file", func(){
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer GinkgoRecover()			

			Ω(process.Decode(w, r, data)).To(BeNil())
		}))
		defer ts.Close()

		resp, err := http.Get(ts.URL)
		Ω(err).To(BeNil())
		Ω(resp.StatusCode).To(Equal(200))
		buff := new(bytes.Buffer)
		buff.ReadFrom(resp.Body)
		original, err := ioutil.ReadFile("tests/test.jpg")
		Ω(err).To(BeNil())
		Ω(bytes.Equal(buff.Bytes(), original)).To(BeTrue())		
	})

	It("sould DELETE file", func(){
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer GinkgoRecover()			

			Ω(process.Remove(w, r, data)).To(BeNil())
		}))
		defer ts.Close()

		req, err := http.NewRequest("DELETE", ts.URL, nil)
		Ω(err).To(BeNil())
		client := &http.Client{}
		resp, err := client.Do(req)
		Ω(err).To(BeNil())
		Ω(resp.StatusCode).To(Equal(204))
	})

	
	
})
