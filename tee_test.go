package sprocess_test

import (
	"bytes"
	. "github.com/hyperboloide/sprocess"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"image"
	"io/ioutil"
	"os"
)

var _ = Describe("Tee", func() {

	testBin := genBlob(1 << 22)

	It("should Encode", func() {
		out1 := new(bytes.Buffer)
		data := NewData()
		id := GenId()
		data.Set("identifier", id)

		output := &File{
			Dir:  "/tmp/" + GenId(),
			Name: "file",
		}
		Ω(output.Start()).To(BeNil())

		tee := &Tee{
			Output: output,
			Name:   "tee",
		}
		Ω(tee.Start()).To(BeNil())

		Ω(tee.Start()).To(BeNil())
		Ω(tee.Encode(
			bytes.NewReader(testBin),
			out1,
			data)).To(BeNil())
		Ω(bytes.Equal(out1.Bytes(), testBin)).To(BeTrue())
		file, err := ioutil.ReadFile(output.Dir + "/" + id)
		Ω(err).To(BeNil())
		Ω(bytes.Equal(file, testBin)).To(BeTrue())

	})

	It("should do service with tee", func() {
		data := NewData()
		id := "pic"
		dir := "/tmp/" + GenId()

		outputLarge := &File{
			Suffix: ".jpg",
			Dir:    dir,
			Name:   "file",
		}
		Ω(outputLarge.Start()).To(BeNil())

		outputSmall := &File{
			Prefix: "small_",
			Suffix: ".jpg",
			Dir:    dir,
			Name:   "file",
		}
		Ω(outputSmall.Start()).To(BeNil())

		size := &Size{
			Name: "size",
		}
		Ω(size.Start()).To(BeNil())

		imgLarge := &Image{
			Operation: ImageResize,
			Height:    300,
			Output:    "jpg",
			Name:      "resizeLarge",
		}
		Ω(imgLarge.Start()).To(BeNil())

		imgSmall := &Image{
			Operation: ImageResize,
			Height:    100,
			Output:    "jpg",
			Name:      "resizeSmall",
		}
		Ω(imgSmall.Start()).To(BeNil())

		tee := &Tee{
			Encoders: []Encoder{imgSmall, size},
			Output:   outputSmall,
			Name:     "tee",
		}
		Ω(tee.Start()).To(BeNil())

		service := &Service{
			EncodingPipe: &EncodingPipeline{
				Encoders: []Encoder{tee, imgLarge, size},
				Output:   outputLarge,
			},
		}

		Ω(service.Encode(id, testFileReader(), data)).To(BeNil())
		large, err := os.Open(dir + "/" + id + ".jpg")
		Ω(err).To(BeNil())
		img, format, err := image.Decode(large)
		Ω(err).To(BeNil())
		Ω(format).To(Equal("jpeg"))
		Ω(img).ToNot(BeNil())
		Ω(img.Bounds().Size().Y).To(Equal(300))

		small, err := os.Open(dir + "/small_" + id + ".jpg")
		Ω(err).To(BeNil())
		img, format, err = image.Decode(small)
		Ω(err).To(BeNil())
		Ω(format).To(Equal("jpeg"))
		Ω(img).ToNot(BeNil())
		Ω(img.Bounds().Size().Y).To(Equal(100))

		d := data.Export()
		dtee := d["tee"].(map[string]interface{})

		Ω(d["size"].(int64) > dtee["size"].(int64)).To(BeTrue())

	})

	It("should do service with tee error", func() {
		data := NewData()
		id := "pic.jpg"
		dir := "/tmp/" + GenId()

		outputLarge := &File{
			Dir:  dir + "/large",
			Name: "file",
		}
		Ω(outputLarge.Start()).To(BeNil())

		outputSmall := &File{
			Dir:  dir + "/small",
			Name: "file",
		}
		Ω(outputSmall.Start()).To(BeNil())

		imgLarge := &Image{
			Operation: ImageResize,
			Height:    300,
			Output:    "jpg",
			Name:      "resizeLarge",
		}
		Ω(imgLarge.Start()).To(BeNil())

		crash := &Bash{
			Cmd:  "exit 1",
			Name: "crash",
		}
		crash.Start()

		tee := &Tee{
			Encoders: []Encoder{crash},
			Output:   outputSmall,
			Name:     "tee",
		}
		Ω(tee.Start()).To(BeNil())

		service := &Service{
			EncodingPipe: &EncodingPipeline{
				Encoders: []Encoder{tee, imgLarge},
				Output:   outputLarge,
			},
		}

		Ω(service.Encode(id, testFileReader(), data)).ToNot(BeNil())
	})

})
