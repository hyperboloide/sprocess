package sprocess_test

import (
	"bytes"
	. "github.com/hyperboloide/sprocess"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	"os"
)

var _ = Describe("File", func() {

	testBin := genBlob(1 << 22)
	data := NewData()
	id := GenId()
	f := &File{
		Dir:  tmpDir(),
		Name: "file",
	}

	It("should Write", func() {
		Ω(f.Start()).To(BeNil())
		w, err := f.NewWriter(id, data)
		Ω(err).To(BeNil())
		Ω(w).ToNot(BeNil())
		l, err := io.Copy(w, bytes.NewReader(testBin))
		w.Close()
		Ω(err).To(BeNil())
		Ω(len(testBin) == int(l)).To(BeTrue())
	})

	It("should read", func() {
		r, err := f.NewReader(id, data)
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
		Ω(f.Delete(id, data)).To(BeNil())
		_, err := os.Stat(f.Dir + "/" + id)
		Ω(os.IsNotExist(err)).To(BeTrue())
	})
})
